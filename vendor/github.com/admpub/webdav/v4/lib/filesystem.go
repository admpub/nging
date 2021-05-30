package lib

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/net/webdav"
)

var _ webdav.FileSystem = WebDavFS{}

type FSGenerator func(scope string, options map[string]string) webdav.FileSystem

type FS struct {
	Scope   string
	FS      FSGenerator
	Options map[string]string
}

func (f FS) Stat(ctx context.Context, name string) (fi os.FileInfo, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("%v", panicErr)
		}
	}()
	return f.FS(f.Scope, f.Options).Stat(ctx, name)
}

func (f FS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (file webdav.File, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("%v", panicErr)
		}
	}()
	return f.FS(f.Scope, f.Options).OpenFile(ctx, name, flag, perm)
}

func (f FS) Mkdir(ctx context.Context, name string, perm os.FileMode) (err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("%v", panicErr)
		}
	}()
	return f.FS(f.Scope, f.Options).Mkdir(ctx, name, perm)
}

func (f FS) RemoveAll(ctx context.Context, name string) (err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("%v", panicErr)
		}
	}()
	return f.FS(f.Scope, f.Options).RemoveAll(ctx, name)
}

func (f FS) Rename(ctx context.Context, oldName, newName string) (err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("%v", panicErr)
		}
	}()
	return f.FS(f.Scope, f.Options).Rename(ctx, oldName, newName)
}

type WebDavFS struct {
	FS
	User    *User
	NoSniff bool
}

func (d WebDavFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	// Skip wrapping if NoSniff is off
	if !d.NoSniff {
		return d.FS.Stat(ctx, name)
	}

	info, err := d.FS.Stat(ctx, name)
	if err != nil {
		return nil, err
	}

	if name != `/` && d.User != nil {
		if !d.User.Allowed(name, true) {
			return nil, filepath.SkipDir
		}
	}

	return NoSniffFileInfo{info}, nil
}

func (d WebDavFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	// Skip wrapping if NoSniff is off
	if !d.NoSniff {
		return d.FS.OpenFile(ctx, name, flag, perm)
	}

	file, err := d.FS.OpenFile(ctx, name, flag, perm)
	if err != nil {
		return nil, err
	}

	return WebDavFile{File: file}, nil
}
