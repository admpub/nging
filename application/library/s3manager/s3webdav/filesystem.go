package s3webdav

import (
	"context"
	"os"
	"path/filepath"

	"github.com/admpub/nging/application/library/s3manager"
	"golang.org/x/net/webdav"
)

var _ webdav.FileSystem = &FileSystem{}

func New(mgr *s3manager.S3Manager, memoryUploadMode bool, uploadTmpPath string) *FileSystem {
	f := &FileSystem{
		mgr:              mgr,
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
	memoryUploadMode bool
	uploadTmpPath    string
}

func defaultUploadTempDir() string {
	uploadTmpPath := filepath.Join(os.TempDir(), `s3webdav`)
	_, err := os.Stat(uploadTmpPath)
	if err != nil && !os.IsNotExist(err) {
		os.MkdirAll(uploadTmpPath, os.ModePerm)
		os.Chmod(uploadTmpPath, os.ModePerm)
	}
	return uploadTmpPath
}

func (f *FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return f.mgr.Mkdir(name, ``)
}

func (f *FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	object, err := f.mgr.GetWithContext(ctx, name)
	if err != nil {
		return nil, err
	}

	return s3manager.NewFile(f.mgr, object, name, f.memoryUploadMode, f.uploadTmpPath), nil
}

func (f *FileSystem) RemoveAll(ctx context.Context, name string) error {
	return f.mgr.RemoveWithContext(ctx, name)
}

func (f *FileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return f.mgr.Rename(oldName, newName)
}

func (f *FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	objectInfo, err := f.mgr.Stat(name)
	if err != nil {
		return nil, err
	}
	return s3manager.NewFileInfo(objectInfo), nil
}
