// +build go1.16

package echo

import (
	"errors"
	"io/fs"
	"path/filepath"
)

func NewFileSystems() FileSystems {
	return FileSystems{}
}

type FileSystems []fs.FS

func (f FileSystems) Open(name string) (file fs.File, err error) {
	for _, fileSystem := range f {
		file, err = fileSystem.Open(name)
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

// EmbedFile
// e.Get(`/*`, EmbedFile(customFS))
func EmbedFile(fs FileSystems) func(c Context) error {
	return func(c Context) error {
		file := c.Param(`*`)
		if len(file) == 0 {
			file = `index.html`
		}
		f, err := fs.Open(file)
		if err != nil {
			return ErrNotFound
		}
		defer func() {
			if f != nil {
				f.Close()
			}
		}()
		fi, err := f.Stat()
		if err != nil {
			return err
		}
		if fi.IsDir() {
			f.Close()

			file = filepath.Join(file, "index.html")
			if f, err = fs.Open(file); err != nil {
				return ErrNotFound
			}

			if fi, err = f.Stat(); err != nil {
				return err
			}
		}
		return c.ServeContent(f, fi.Name(), fi.ModTime())
	}
}
