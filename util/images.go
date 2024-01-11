package util

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"mime/multipart"
)

func ConvertToJPEG(imageToConvert *multipart.File) (*bytes.Buffer, error) {
	originalImage, _, err := image.Decode(*imageToConvert)
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
