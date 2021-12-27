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
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	minio "github.com/minio/minio-go/v7"
)

func NewFileInfo(objectInfo minio.ObjectInfo) os.FileInfo {
	return &fileInfo{objectInfo: objectInfo}
}

func NewStrFileInfo(prefix string) os.FileInfo {
	return &fileInfo{
		objectInfo: minio.ObjectInfo{
			Key: prefix,
		},
	}
}

func NewS3FileInfo(object *s3.Object) os.FileInfo {
	objectInfo := minio.ObjectInfo{}
	if object.ETag != nil {
		objectInfo.ETag = *object.ETag
	}
	if object.Key != nil {
		objectInfo.Key = *object.Key
	}
	if object.LastModified != nil {
		objectInfo.LastModified = *object.LastModified
	}
	if object.Owner != nil {
		if object.Owner.DisplayName != nil {
			objectInfo.Owner.DisplayName = *object.Owner.DisplayName
		}
		if object.Owner.ID != nil {
			objectInfo.Owner.ID = *object.Owner.ID
		}
	}
	if object.Size != nil {
		objectInfo.Size = *object.Size
	}
	if object.StorageClass != nil {
		objectInfo.StorageClass = *object.StorageClass
	}
	return &fileInfo{
		objectInfo: objectInfo,
	}
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
