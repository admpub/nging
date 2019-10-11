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

package helper

import (
	"regexp"
	"strings"

	"github.com/webx-top/client/upload"
	"github.com/webx-top/com"
)

const (
	DefaultUploadURLPath = `/public/upload/`
	DefaultUploadDir     = `./public/upload`
)

var (
	// UploadURLPath 上传文件网址访问路径
	UploadURLPath = DefaultUploadURLPath

	// UploadDir 定义上传目录（首尾必须带“/”）
	UploadDir = DefaultUploadDir

	// AllowedUploadFileExtensions 被允许上传的文件的扩展名
	AllowedUploadFileExtensions = []string{
		`.jpeg`, `.jpg`, `.gif`, `.png`,
	}

	// FileTypeIcon 文件类型icon
	FileTypeIcon = upload.FileTypeIcon

	// DetectFileType 根据文件扩展名判断文件类型
	DetectFileType = upload.DetectType

	// TypeRegister 注册文件扩展名
	TypeRegister = upload.TypeRegister
)

// URLToFile 文件网址转为存储路径
func URLToFile(fileURL string) string {
	filePath := strings.TrimPrefix(fileURL, UploadURLPath)
	filePath = strings.TrimSuffix(UploadDir, `/`) + `/` + strings.TrimPrefix(filePath, `/`)
	return filePath
}

// FileTypeByName 根据文件名判断文件类型
func FileTypeByName(filename string) string {
	p := strings.LastIndex(filename, `.`)
	if p < 0 {
		return ``
	}
	ext := filename[p:]
	return DetectFileType(ext)
}

func ExtensionRegister(extensions ...string) {
	AllowedUploadFileExtensions = append(AllowedUploadFileExtensions, extensions...)
}

func ExtensionUnregister(extensions ...string) {
	com.SliceRemoveCallback(len(AllowedUploadFileExtensions), func(i int) func(bool) error {
		if !com.InStringSlice(AllowedUploadFileExtensions[i], extensions) {
			return nil
		}
		return func(inside bool) error {
			if inside {
				AllowedUploadFileExtensions = append(AllowedUploadFileExtensions[0:i], AllowedUploadFileExtensions[i+1:]...)
			} else {
				AllowedUploadFileExtensions = AllowedUploadFileExtensions[0:i]
			}
			return nil
		}
	})
}

func ExtensionRegexpEnd() string {
	extensions := make([]string, len(AllowedUploadFileExtensions))
	for index, extension := range AllowedUploadFileExtensions {
		extensions[index] = regexp.QuoteMeta(extension)
	}
	return `(` + strings.Join(extensions, `|`) + `)`
}
