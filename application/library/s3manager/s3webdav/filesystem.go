package s3webdav

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/admpub/nging/v4/application/library/s3manager"
	"github.com/webx-top/com"
	"golang.org/x/net/webdav"
)

var _ webdav.FileSystem = &FileSystem{}

func New(mgr *s3manager.S3Manager, scope string, memoryUploadMode bool, uploadTmpPath string) *FileSystem {
	f := &FileSystem{
		mgr:              mgr,
		scope:            scope,
		memoryUploadMode: memoryUploadMode,
		uploadTmpPath:    uploadTmpPath,
	}
	if len(f.uploadTmpPath) == 0 {
		f.uploadTmpPath = defaultUploadTempDir()
	}
	return f
}

type FileSystem struct {
	mgr              *s3manager.S3Manager
	scope            string
	memoryUploadMode bool
	uploadTmpPath    string
}

func defaultUploadTempDir() string {
	uploadTmpPath := filepath.Join(os.TempDir(), `s3webdav`)
	com.MkdirAll(uploadTmpPath, os.ModePerm)
	return uploadTmpPath
}

// slashClean is equivalent to but slightly more efficient than
// path.Clean("/" + name).
func slashClean(name string) string {
	if name == "" || name[0] != '/' {
		name = "/" + name
	}
	return path.Clean(name)
}

func (f *FileSystem) resolve(name string) string {
	// This implementation is based on Dir.Open's code in the standard net/http package.
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) ||
		strings.Contains(name, "\x00") {
		return ""
	}
	dir := f.scope
	if len(dir) == 0 {
		dir = "."
	}
	return filepath.Join(dir, filepath.FromSlash(slashClean(name)))
}

func (f *FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	if name = f.resolve(name); len(name) == 0 {
		return os.ErrNotExist
	}
	return f.mgr.Mkdir(ctx, name, ``)
}

func (f *FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if name = f.resolve(name); len(name) == 0 {
		return nil, os.ErrNotExist
	}
	object, err := f.mgr.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	return s3manager.NewFile(f.mgr, object, name, f.memoryUploadMode, f.uploadTmpPath), nil
}

func (f *FileSystem) RemoveAll(ctx context.Context, name string) error {
	if name = f.resolve(name); len(name) == 0 {
		return os.ErrNotExist
	}
	return f.mgr.Remove(ctx, name)
}

func (f *FileSystem) Rename(ctx context.Context, oldName, newName string) error {
	if oldName = f.resolve(oldName); len(oldName) == 0 {
		return os.ErrNotExist
	}
	if newName = f.resolve(newName); len(newName) == 0 {
		return os.ErrNotExist
	}
	if root := filepath.Clean(f.scope); root == oldName || root == newName {
		// Prohibit renaming from or to the virtual root directory.
		return os.ErrInvalid
	}
	return f.mgr.Rename(ctx, oldName, newName)
}

func (f *FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	if name = f.resolve(name); len(name) == 0 {
		return nil, os.ErrNotExist
	}
	objectInfo, err := f.mgr.Stat(ctx, name)
	if err != nil {
		return nil, err
	}
	return s3manager.NewFileInfo(objectInfo), nil
}
