package s3manager

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	minio "github.com/minio/minio-go"
	"github.com/webx-top/com"
)

func NewFile(mgr *S3Manager, object *minio.Object, name string, memoryUploadMode bool, uploadTmpPath string) *file {
	return &file{
		mgr:              mgr,
		Object:           object,
		name:             name,
		memoryUploadMode: memoryUploadMode,
		uploadTmpPath:    uploadTmpPath,
	}
}

type file struct {
	mgr *S3Manager
	*minio.Object
	name             string
	memoryUploadMode bool
	uploadTmpPath    string
}

func (f *file) Stat() (os.FileInfo, error) {
	fi, err := f.mgr.Stat(f.name)
	if err != nil {
		return nil, err
	}
	return NewFileInfo(fi), nil
}

func (f *file) ReadFrom(r io.Reader) (n int64, err error) {

	// memory mode
	if f.memoryUploadMode {
		return f.mgr.PutObject(r, f.name, -1)
	}

	// file mode
	tmpFilePath := filepath.Join(f.uploadTmpPath, com.Md5(f.name))
	var fp *os.File
	fp, err = os.Create(tmpFilePath)
	if err != nil {
		return 0, err
	}
	defer fp.Close()
	defer func(p string) {
		err = os.RemoveAll(p)
	}(tmpFilePath)

	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := f.Write(buf[:n]); err != nil {
			return 0, err
		}
	}
	return f.mgr.FPutObject(tmpFilePath, f.name)
}

func (f *file) Write(p []byte) (n int, err error) {
	return len(p), nil // useless
}

func (f *file) Readdir(count int) (fileInfoList []os.FileInfo, err error) {
	objectPrefix := strings.TrimPrefix(f.name, `/`)
	words := len(objectPrefix)
	if words > 0 {
		if !strings.HasSuffix(objectPrefix, `/`) {
			objectPrefix += `/`
		}
	}
	return f.mgr.listByMinio(objectPrefix)
}
