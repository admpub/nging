package common

import (
	"errors"
	"io/fs"
)

var FileSystem FileSystems

type (
	FileSystems []fs.FS
)

func (f FileSystems) Open(name string) (file fs.File, err error) {
	for _, i := range f {
		file, err = i.Open(name)
		if err == nil || !errors.Is(err, fs.ErrNotExist) {
			return
		}
	}
	return
}

func (f FileSystems) Size() int {
	return len(f)
}

func (f FileSystems) IsEmpty() bool {
	return f.Size() == 0
}

func (f *FileSystems) Register(fileSystem fs.FS) {
	*f = append(*f, fileSystem)
}
