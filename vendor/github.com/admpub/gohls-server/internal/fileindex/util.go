package fileindex

import (
	"os"
)

type dirInfo struct {
	path     string
	info     os.FileInfo
	children []os.FileInfo
}

func readDirInfo(path string) (*dirInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return &dirInfo{path, fi, list}, nil
}
