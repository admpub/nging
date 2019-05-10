// +build zbar

package qrcode

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/admpub/qrcode/decode"
)

func Decode(file *os.File) (string, error) {
	var body string
	var img image.Image
	var err error
	imageTypeArr := strings.Split(file.Name(), ".")
	if len(imageTypeArr) <= 1 {
		err = errors.New("Image file format error")
		return body, err
	}

	imageType := imageTypeArr[len(imageTypeArr)-1]

	switch imageType {
	case "jpeg", "jpg":
		img, err = jpeg.Decode(file)
	case "png":
		img, err = png.Decode(file)
	default:
		err = errors.New("Image file format error")
		return body, err
	}

	if err != nil {
		err = errors.New("decode failed: " + err.Error())
		return body, err
	}

	newImg := decode.NewImage(img)
	scanner := decode.NewScanner().SetEnabledAll(true)

	symbols, _ := scanner.ScanImage(newImg)
	for _, s := range symbols {
		body += s.Data
	}

	return body, err
}

func DecodeFile(imgPath string) (string, error) {
	fi, err := os.Open(imgPath)
	if err != nil {
		return ``, err
	}
	defer fi.Close()
	return Decode(fi)
}
