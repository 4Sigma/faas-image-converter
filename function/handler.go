package function

import (
	"encoding/json"
	"fmt"
	"handler/function/pkg/imageprocessing"
	"handler/function/pkg/storage"
	"handler/function/pkg/utils"
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

	fileName, err := storage.DownloadFile(dataImage)
	if err != nil {
		return handler.Response{}, err
	}

	fmt.Println("File name: ", fileName)

	vips.Startup(nil)
	defer vips.Shutdown()

	fmt.Println("File name: ", fileName)

	/* 	for _, format := range outputFormat {
		for _, size := range format.Size {
			fmt.Println("Format: ", format.Format, " Size: ", size.Width, "x", size.Height)
			ep := vips.NewDefaultJPEGExportParams()
		}
	} */
	//ep := vips.NewThumbnailWithSizeFromFile(image, 200, 200)
	imageprocessing.ImageConverter(dataImage.OutputFormats, fileName)

	return handler.Response{
		Body:       []byte(out),
		StatusCode: http.StatusOK,
	}, err
}
