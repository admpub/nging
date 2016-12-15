package server

import "os"

type FileInfo interface {
	os.FileInfo

	Owner() string
	Group() string
}
