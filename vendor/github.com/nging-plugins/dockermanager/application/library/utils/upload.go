package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/webx-top/com"
)

func SaveUploadedFile(file *multipart.FileHeader, srcPath string) (string, error) {
	fp, err := file.Open()
	if err != nil {
		return ``, err
	}
	defer fp.Close()
	return SaveMultipartFile(fp, file.Filename, srcPath)
}

func SaveMultipartFile(fp multipart.File, fileName string, srcPath string) (string, error) {
	sf := filepath.Join(srcPath, fileName)
	sp, err := os.Create(sf)
	if err != nil {
		if !os.IsNotExist(err) {
			return sf, err
		}
		err = com.MkdirAll(filepath.Dir(sf), os.ModePerm)
		if err != nil {
			return sf, err
		}
		sp, err = os.Create(sf)
		if err != nil {
			return sf, err
		}
	}
	defer sp.Close()
	_, err = io.Copy(sp, fp)
	if err != nil {
		return sf, err
	}
	err = sp.Sync()
	return sf, err
}
