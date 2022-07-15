package utils

import (
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/graymeta/stow/google"
	_ "github.com/graymeta/stow/s3"
)

type OutputFormats []struct {
	Format string `json:"format"`
	Size   []struct {
		Width  string `json:"width"`
		Height string `json:"height"`
	} `json:"size"`
}

type ConverterInput struct {
	InputImage struct {
		StorageType string      `json:"storageType"`
		StorageData interface{} `json:"storageData"`
	} `json:"inputImage"`
	OutputFormats OutputFormats `json:"outputFormats"`
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func ImageConverter(httpBody []byte) ConverterInput {
	var converterInput ConverterInput
	err := json.Unmarshal(httpBody, &converterInput)
	if err != nil {
		fmt.Println(err)
	}
	return converterInput
}
