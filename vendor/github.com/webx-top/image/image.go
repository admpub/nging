package image

import (
	"image"
	"os"
)

func ParseImage(filePath string) (img image.Image, err error) {
	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	img, _, err = image.Decode(fp)
	fp.Close()
	return img, err
}
