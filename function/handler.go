package function

import (
	"encoding/json"
	"handler/function/pkg/imageprocessing"
	"handler/function/pkg/storage"
	"handler/function/pkg/utils"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/davidbyttow/govips/v2/vips"
	handler "github.com/openfaas/templates-sdk/go-http"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Print("hello world")

	var dataImage utils.ImageGeneration

	dataImage = utils.ImageConverter(req.Body)

	out, err := json.Marshal(dataImage)
	if err != nil {
		return handler.Response{}, err
	}

	fileName, err := storage.DownloadFile(dataImage)
	if err != nil {
		return handler.Response{}, err
	}

	log.Print("File name: ", fileName)

	vips.LoggingSettings(vipsLogger, vips.LogLevelInfo)
	vips.Startup(nil)
	defer vips.Shutdown()

	imageprocessing.ImageConverter(dataImage.OutputFormats, fileName)

	return handler.Response{
		Body:       []byte(out),
		StatusCode: http.StatusOK,
	}, err
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
