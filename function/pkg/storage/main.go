package storage

import (
	"errors"
	"fmt"
	"handler/function/pkg/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"

	"github.com/graymeta/stow"
	"github.com/graymeta/stow/s3"
	_ "github.com/graymeta/stow/s3"
)

func InitStorage() (stow.Location, error) {
	fmt.Println("Initializing storage")
	kind := "s3"
	config := stow.ConfigMap{
		s3.ConfigEndpoint:    "s3.eu-central-2.wasabisys.com",
		s3.ConfigAccessKeyID: "2HQI5JE5XTG66NGAQUMV",
		s3.ConfigSecretKey:   "2e6BKeDem1j6nK8ieV7qcc554hOKb5DHrSdsNLRU",
		s3.ConfigRegion:      "eu-central-2",
	}

	location, err := stow.Dial(kind, config)
	if err != nil {
		fmt.Println("Error dialing", kind, ":", err)
	}
	return location, err
}

func DownloadFile(image utils.ImageGeneration) error {
	if image.Img.StorageType == "remote-http" {
		err := downloadRemoteHttp(image)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Unknown storage type")
	}
	return nil

}

func downloadRemoteHttp(image utils.ImageGeneration) error {
	imageUrl := image.Img.Src
	fileName := filepath.Base(imageUrl)

	//Get the response bytes from the url
	response, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fExtension := filepath.Ext(fileName)
	if fExtension == "" {
		kind, _ := filetype.Match(buf)
		if kind == filetype.Unknown {
			return errors.New("Unknown file type")
		}
		fExtension = kind.Extension
		fileName = "tmp/" + fileName + "." + fExtension
	}

	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	_, err = file.Write(buf)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

/* err = stow.WalkContainers(location, stow.NoPrefix, 100,
	func(c stow.Container, err error) error {
		if err != nil {
			return err
		}
		log.Println("Container: ", c.Name())
		return nil
	})

if err != nil {
	return handler.Response{}, err
}

err = stow.Walk(container, stow.NoPrefix, 100,
	func(item stow.Item, err error) error {
		if err != nil {
			return err
		}
		log.Println(item.Name())
		return nil
	})
if err != nil {
	return err
} */
