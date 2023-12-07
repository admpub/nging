/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package s3manager

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/nging/v5/application/library/s3manager/fileinfo"
	minio "github.com/minio/minio-go/v7"
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
	fi, err := f.mgr.Stat(context.Background(), f.name)
	if err != nil {
		return nil, err
	}
	return fileinfo.New(fi), nil
}

func (f *file) ReadFrom(r io.Reader) (n int64, err error) {
	ctx := context.Background()
	// memory mode
	if f.memoryUploadMode {
		return f.mgr.PutObject(ctx, r, f.name, -1)
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
	return f.mgr.FPutObject(ctx, tmpFilePath, f.name)
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
	return f.mgr.listByMinio(context.Background(), objectPrefix)
}
