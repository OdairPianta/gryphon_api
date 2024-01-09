package models

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/disintegration/imaging"
)

type Image struct {
	PublicURL string `gorm:"size:2048" json:"public_url"`
}

func ResizeImage(img image.Image, width int, height int) image.Image {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

func ConvertBase64IntoImage(base64Content string) (image.Image, error) {
	base64Decode, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(string(base64Decode))

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ConvertImageIntoByte(image image.Image, extension string) ([]byte, error) {
	buf := new(bytes.Buffer)
	var err error
	if extension == "jpeg" || extension == "jpg" {
		err = jpeg.Encode(buf, image, nil)
	} else if extension == "png" {
		err = png.Encode(buf, image)
	} else if extension == "gif" {
		err = gif.Encode(buf, image, nil)
	} else {
		return nil, fmt.Errorf("formato de imagem inválido")
	}

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
