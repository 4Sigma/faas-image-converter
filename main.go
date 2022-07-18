package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"handler/function"

	handler "github.com/openfaas/templates-sdk/go-http"
)

var (
	acceptingConnections int32
)

const defaultTimeout = 10 * time.Second

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	readTimeout := parseIntOrDurationValue(os.Getenv("read_timeout"), defaultTimeout)
	writeTimeout := parseIntOrDurationValue(os.Getenv("write_timeout"), defaultTimeout)
	healthInterval := parseIntOrDurationValue(os.Getenv("healthcheck_interval"), writeTimeout)

	vips.LoggingSettings(vipsLogger, vips.LogLevelInfo)
	vips.Startup(nil)
	defer vips.Shutdown()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8082),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
	}

	http.HandleFunc("/", imageResizeView)
	listenUntilShutdown(s, healthInterval, writeTimeout)
}

func listenUntilShutdown(s *http.Server, shutdownTimeout time.Duration, writeTimeout time.Duration) {
	idleConnsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM)

		<-sig

		log.Printf("[entrypoint] SIGTERM: no connections in: %s", shutdownTimeout.String())
		<-time.Tick(shutdownTimeout)

		ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Printf("[entrypoint] Error in Shutdown: %v", err)
		}

		log.Printf("[entrypoint] Exiting.")

		close(idleConnsClosed)
	}()

	// Run the HTTP server in a separate go-routine.
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("[entrypoint] Error ListenAndServe: %v", err)
			close(idleConnsClosed)
		}
	}()

	atomic.StoreInt32(&acceptingConnections, 1)

	<-idleConnsClosed
}

func imageResizeView(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		bodyBytes, bodyErr := io.ReadAll(r.Body)

		if bodyErr != nil {
			log.Printf("Error reading body from request.")
		}

		input = bodyBytes
	}

	req := handler.Request{
		Body:        input,
		Header:      r.Header,
		Method:      r.Method,
		QueryString: r.URL.RawQuery,
	}
	req.WithContext(r.Context())

	result, resultErr := function.Handle(req)

	if result.Header != nil {
		for k, v := range result.Header {
			w.Header()[k] = v
		}
	}

	if resultErr != nil {
		log.Print(resultErr)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		if result.StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(result.StatusCode)
		}
	}

	w.Write(result.Body)
}

func parseIntOrDurationValue(val string, fallback time.Duration) time.Duration {
	if len(val) > 0 {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal >= 0 {
			return time.Duration(parsedVal) * time.Second
		}
	}

	duration, durationErr := time.ParseDuration(val)
	if durationErr != nil {
		return fallback
	}
	return duration
}

func vipsLogger(messageDomain string, verbosity vips.LogLevel, message string) {
	var messageLevelDescription string
	switch verbosity {
	case vips.LogLevelError:
		messageLevelDescription = "error"
	case vips.LogLevelCritical:
		messageLevelDescription = "critical"
	case vips.LogLevelWarning:
		messageLevelDescription = "warning"
	case vips.LogLevelMessage:
		messageLevelDescription = "message"
	case vips.LogLevelInfo:
		messageLevelDescription = "info"
	case vips.LogLevelDebug:
		messageLevelDescription = "debug"
	}

	log.Print(messageDomain, messageLevelDescription, message)
}
