package util

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"mime/multipart"
)

func ConvertToJPEG(imageToConvert *multipart.FileHeader) (*bytes.Buffer, error) {
	openedImage, err := imageToConvert.Open()
	if err != nil {
		return nil, err
	}
	defer openedImage.Close()

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

// func jpegToBuffer(file *multipart.FileHeader) (*bytes.Buffer, error) {
// 	openedFile, err := file.Open()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer openedFile.Close()

// 	fileAsImage, _, err := image.Decode(openedFile)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var buffer bytes.Buffer
// 	err = jpeg.Encode(&buffer, fileAsImage, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &buffer, nil
// }
