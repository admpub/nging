// +build !zbar

package qrcode

import (
	"os"

	"github.com/tuotoo/qrcode"
)

func Decode(file *os.File) (string, error) {
	qrmatrix, err := qrcode.Decode(file)
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
	return Decode(fi)
}
