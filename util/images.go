package util

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
)

func ConvertToJPEG(imageToConvert *multipart.FileHeader) (io.Reader, error) {
	openedImage, err := imageToConvert.Open()
	if err != nil {
		return nil, err
	}
	defer openedImage.Close()

	if imageToConvert.Header.Get("Content-Type") == "image/jpeg" {
		return openedImage, nil
	}

	originalImage, _, err := image.Decode(openedImage)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = jpeg.Encode(&buffer, originalImage, nil)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}
