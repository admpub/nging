// +build !zbar

package qrcode

import (
	"io"
	"os"

	"github.com/tuotoo/qrcode"
)

func Decode(reader io.Reader, imageType string) (string, error) {
	qrmatrix, err := qrcode.Decode(reader)
	if err != nil {
		return ``, err
	}
	return qrmatrix.Content, err
}

func DecodeFile(imgPath string) (string, error) {
	fi, err := os.Open(imgPath)
	if err != nil {
		return ``, err
	}
	defer fi.Close()
	return Decode(fi, ``)
}
