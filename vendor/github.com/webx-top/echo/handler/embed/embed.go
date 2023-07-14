//go:build go1.16
// +build go1.16

package embed

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/webx-top/echo"
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

func (f FileSystems) ReadFile(name string) (content []byte, err error) {
	for _, fileSystem := range f {
		content, err = fs.ReadFile(fileSystem, name)
		if err == nil || !errors.Is(err, fs.ErrNotExist) {
			return
		}
	}
	return
}

func (f FileSystems) ReadDir(name string) (dirs []fs.DirEntry, err error) {
	unique := map[string]struct{}{}
	for _, fileSystem := range f {
		var _dirs []fs.DirEntry
		_dirs, err = fs.ReadDir(fileSystem, name)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return
		}
		err = nil
		for _, dir := range _dirs {
			if _, ok := unique[dir.Name()]; ok {
				continue
			}
			unique[dir.Name()] = struct{}{}
			dirs = append(dirs, dir)
		}
	}
	return
}

func (f FileSystems) Sub(name string) (sub fs.FS, err error) {
	for _, fileSystem := range f {
		sub, err = fs.Sub(fileSystem, name)
		if err == nil || !errors.Is(err, fs.ErrNotExist) {
			return
		}
	}
	return
}

func (f FileSystems) WalkDir(name string, fn fs.WalkDirFunc) (err error) {
	for _, fileSystem := range f {
		err = fs.WalkDir(fileSystem, name, fn)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return
		}
		err = nil
	}
	return
}

func (f FileSystems) Stat(name string) (fi fs.FileInfo, err error) {
	for _, fileSystem := range f {
		fi, err = fs.Stat(fileSystem, name)
		if err == nil || !errors.Is(err, fs.ErrNotExist) {
			return
		}
	}
	return
}

func (f FileSystems) Glob(pattern string) (matches []string, err error) {
	unique := map[string]struct{}{}
	for _, fileSystem := range f {
		var mch []string
		mch, err = fs.Glob(fileSystem, pattern)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return
		}
		err = nil
		for _, match := range mch {
			if _, ok := unique[match]; ok {
				continue
			}
			unique[match] = struct{}{}
			matches = append(matches, match)
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

// File
// e.Get(`/*`, File(customFS))
func File(fs FileSystems, configs ...Config) func(c echo.Context) error {
	config := DefaultConfig
	if len(configs) > 0 {
		config = configs[0]
		if len(config.Index) == 0 {
			config.Index = DefaultConfig.Index
		}
		if len(config.Prefix) > 0 {
			config.Prefix = strings.TrimPrefix(config.Prefix, `/`)
		}
		if len(config.Prefix) > 0 {
			if !strings.HasSuffix(config.Prefix, `/`) {
				config.Prefix += `/`
			}
		}
		if config.FilePath == nil {
			config.FilePath = DefaultConfig.FilePath
		}
	}
	return func(c echo.Context) error {
		file, err := config.FilePath(c)
		if err != nil {
			return err
		}
		if len(file) == 0 {
			file = config.Index
		}
		if len(config.Prefix) > 0 {
			file = config.Prefix + file
		}
		f, err := fs.Open(file)
		if err != nil {
			return echo.ErrNotFound
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

			file = filepath.Join(file, config.Index)
			if f, err = fs.Open(file); err != nil {
				return echo.ErrNotFound
			}

			if fi, err = f.Stat(); err != nil {
				return err
			}
		}
		return c.ServeContent(f, fi.Name(), fi.ModTime())
	}
}
