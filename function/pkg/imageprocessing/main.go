package imageprocessing

import (
	"errors"
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	"handler/function/pkg/utils"
	"io/ioutil"
	"strconv"
)

func ImageConverter(size utils.OutputFormats, fileName string) error {

	fmt.Println(fileName)

	for _, format := range size {
		for _, size := range format.Size {
			fmt.Println("Format: ", format.Format, " Size: ", size.Width, "x", size.Height)
			//outputNameNoSuffix := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))

			tmpIntConverter, err := strconv.ParseInt(size.Width, 0, 32)
			sizeWidth := int(tmpIntConverter)
			if err != nil {
				return err
			}

			tmpIntConverter, err = strconv.ParseInt(size.Height, 0, 32)
			if err != nil {
				return err
			}
			sizeHeight := int(tmpIntConverter)

			// Create new thunbnail
			image, err := vips.NewThumbnailFromFile(fileName, sizeWidth, sizeHeight, vips.InterestingNone)

			if err != nil {
				return err
			}

			var imageBytes []byte
			switch format.Format {
			case "jpg":
				ep := vips.NewJpegExportParams()
				ep.Quality = 95

				imageBytes, _, err = image.ExportJpeg(ep)
				if err != nil {
					return err
				}
			case "webp":
				ep := vips.NewWebpExportParams()
				ep.Quality = 95

				imageBytes, _, err = image.ExportWebp(ep)
				if err != nil {
					return err
				}
			case "png":
				ep := vips.NewPngExportParams()
				ep.Quality = 95

				imageBytes, _, err = image.ExportPng(ep)
				if err != nil {
					return err
				}
			case "gif":
				ep := vips.NewGifExportParams()
				ep.Quality = 95

				imageBytes, _, err = image.ExportGIF(ep)
			case "tiff":
				ep := vips.NewTiffExportParams()
				ep.Quality = 95

				imageBytes, _, err = image.ExportTiff(ep)
			case "jpg2k":
				ep := vips.NewJp2kExportParams()
				ep.Quality = 95

				imageBytes, _, err = image.ExportJp2k(ep)
			case "heif":
				ep := vips.NewHeifExportParams()
				ep.Quality = 95

				imageBytes, _, err = image.ExportHeif(ep)
			case "avif":
				fmt.Println("AVIF")
				ep := vips.NewAvifExportParams()
				ep.Quality = 95

				imageBytes, _, err = image.ExportAvif(ep)
			default:
				return errors.New("Format not supported")
			}

			outputFname := fmt.Sprintf("%s_cropped_%sx%s.%s", fileName, string(size.Width), string(size.Height), format.Format)
			err = ioutil.WriteFile(outputFname, imageBytes, 0644)
			fmt.Println("File name: ", outputFname)
		}
	}
	return nil
}
