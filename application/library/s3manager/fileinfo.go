/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

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
	"os"
	"strings"
	"time"

	minio "github.com/minio/minio-go"
)

func NewFileInfo(objectInfo minio.ObjectInfo) os.FileInfo {
	return &fileInfo{objectInfo: objectInfo}
}

type fileInfo struct {
	objectInfo minio.ObjectInfo
}

func (f *fileInfo) Name() string {
	return f.objectInfo.Key
}

func (f *fileInfo) Size() int64 {
	return f.objectInfo.Size
}

func (f *fileInfo) Mode() os.FileMode {
	return 0
}

func (f *fileInfo) ModTime() time.Time {
	return f.objectInfo.LastModified
}

func (f *fileInfo) IsDir() bool {
	return strings.HasSuffix(f.Name(), "/")
}

func (f *fileInfo) Sys() interface{} {
	return f.objectInfo
}
