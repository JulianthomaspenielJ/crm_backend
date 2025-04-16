package models

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

func imageToBase64(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageData), nil
}