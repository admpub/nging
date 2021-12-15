//go:build zbar
// +build zbar

package qrcode

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/admpub/qrcode/decode"
)

func Decode(reader io.Reader, imageType string) (string, error) {
	var body string
	var img image.Image
	var err error
	if len(imageType) == 0 {
		imageType = `png`
	}
	if name, ok := reader.(Name); ok {
		fileName := name.Name()
		p := strings.LastIndex(fileName, `.`)
		if p < 0 || len(fileName) <= p+1 {
			err = errors.New("Image file format error")
			return body, err
		}
		imageType = fileName[p+1:]
	}

	switch strings.ToLower(imageType) {
	case "jpeg", "jpg":
		img, err = jpeg.Decode(reader)
	case "png":
		img, err = png.Decode(reader)
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
	return Decode(fi, ``)
}
