package function

import (
	"encoding/json"
	"fmt"
	"handler/function/pkg/storage"
	"handler/function/pkg/utils"
	"io/ioutil"
	"net/http"

	"github.com/davidbyttow/govips/v2/vips"
	handler "github.com/openfaas/templates-sdk/go-http"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	var dataImage utils.ImageGeneration

	dataImage = utils.ImageConverter(req.Body)

	out, err := json.Marshal(dataImage)
	if err != nil {
		return handler.Response{}, err
	}
	//location, err := storage.InitStorage()
	// if err != nil {
	// 	return handler.Response{}, err
	// }

	fileName, err := storage.DownloadFile(dataImage)
	if err != nil {
		return handler.Response{}, err
	}

	fmt.Println("File name: ", fileName)

	vips.Startup(nil)
	defer vips.Shutdown()

	image1, err := vips.NewImageFromFile(fileName)
	ep := vips.NewDefaultJPEGExportParams()
	image1bytes, _, err := image1.Export(ep)
	err = ioutil.WriteFile("output.jpg", image1bytes, 0644)

	if err != nil {
		return handler.Response{}, err
	}

	ep = vips.NewDefaultWEBPExportParams()
	image1bytes, _, err = image1.Export(ep)
	err = ioutil.WriteFile("output.webp", image1bytes, 0644)

	if err != nil {
		return handler.Response{}, err
	}

	epaAvif := vips.NewAvifExportParams()
	epaAvif.Quality = 90
	epaAvif.Lossless = false
	image1bytes, _, err = image1.ExportAvif(epaAvif)

	if err != nil {
		return handler.Response{}, err
	}

	err = ioutil.WriteFile("output.avif", image1bytes, 0644)

	return handler.Response{
		Body:       []byte(out),
		StatusCode: http.StatusOK,
	}, err
}
