package utils

import (
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/graymeta/stow/google"
	_ "github.com/graymeta/stow/s3"
)

type ImageGeneration struct {
	Img struct {
		StorageType string `json:"StorageType"`
		Src         string `json:"src"`
	} `json:"img"`
	Format []struct {
		Format string `json:"format"`
		Size   []struct {
			Width  string `json:"width"`
			Height string `json:"height"`
		} `json:"size"`
	} `json:"format"`
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func ImageConverter(httpBody []byte) ImageGeneration {
	var imgData ImageGeneration
	err := json.Unmarshal(httpBody, &imgData)
	if err != nil {
		fmt.Println(err)
	}
	return imgData
}
