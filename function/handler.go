package function

import (
	"encoding/json"
	"handler/function/pkg/imageprocessing"
	"handler/function/pkg/storage"
	"handler/function/pkg/utils"
	"net/http"

	"github.com/rs/zerolog/log"

	handler "github.com/openfaas/templates-sdk/go-http"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {

	var converterInput utils.ConverterInput = utils.ImageConverter(req.Body)

	out, err := json.Marshal(converterInput)
	if err != nil {
		return handler.Response{}, err
	}

	fileName, err := storage.DownloadFile(converterInput)
	if err != nil {
		return handler.Response{}, err
	}

	log.Print("File name: ", fileName)

	imageprocessing.ImageConverter(fileName, converterInput.OutputFormats)

	return handler.Response{
		Body:       []byte(out),
		StatusCode: http.StatusOK,
	}, err
}
