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
	"mime"
	"regexp"
	"strings"

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

	// FileTypeByExtensions 扩展名对应类型
	FileTypeByExtensions = map[string]string{
		// image
		`.jpeg`: `image`,
		`.jpg`:  `image`,
		`.gif`:  `image`,
		`.png`:  `image`,
		`.svg`:  `image`,

		// video
		`.mp4`:  `video`,
		`.mpeg`: `video`,
		`.rmvb`: `video`,
		`.rm`:   `video`,
		`.avi`:  `video`,
		`.mkv`:  `video`,

		// audio
		`.mp3`: `audio`,
		`.mid`: `audio`,

		// archive
		`.rar`:  `archive`,
		`.zip`:  `archive`,
		`.gz`:   `archive`,
		`.bz`:   `archive`,
		`.gzip`: `archive`,
		`.7z`:   `archive`,

		// pdf
		`.pdf`: `pdf`,

		// xls
		`.xls`:  `xls`,
		`.xlsx`: `xls`,
		`.csv`:  `xls`,

		// ppt
		`.ppt`: `ppt`,

		// doc
		`.txt`:  `doc`,
		`.doc`:  `doc`,
		`.docx`: `doc`,
	}

	// FileTypes 文件类型对应mime关键词
	FileTypes = map[string][]string{
		`image`:   []string{`image`},
		`video`:   []string{`video`},
		`audio`:   []string{`audio`},
		`archive`: []string{`compressed`},
		`pdf`:     []string{`pdf`},
		`xls`:     []string{`csv`, `excel`},
		`ppt`:     []string{`powerpoint`},
		`doc`:     []string{`msword`, `text`},
	}

	// FileTypeIcons 文件类型对应icon(不含“fa-”前缀)
	FileTypeIcons = map[string]string{
		`image`:   `picture-o`,
		`video`:   `film`,
		`audio`:   `music`,
		`archive`: `archive`,
		`pdf`:     `file-o`,
		`xls`:     `file-o`,
		`ppt`:     `file-o`,
		`doc`:     `file-text-o`,
	}
)

// FileTypeByName 根据文件名判断文件类型
func FileTypeByName(filename string) string {
	p := strings.LastIndex(filename, `.`)
	if p < 0 {
		return ``
	}
	ext := filename[p:]
	return DetectFileType(ext)
}

// FileTypeIcon 文件类型icon
func FileTypeIcon(typ string) string {
	icon, ok := FileTypeIcons[typ]
	if ok {
		return icon
	}
	return `file-o`
}

// DetectFileType 根据文件扩展名判断文件类型
func DetectFileType(ext string) string {
	ext = strings.ToLower(ext)
	typ, ok := FileTypeByExtensions[ext]
	if ok {
		return typ
	}

	mimeType := mime.TypeByExtension(ext)
	mimeType = strings.SplitN(ext, ";", 2)[0]
	for typeK, keywords := range FileTypes {
		for _, words := range keywords {
			if strings.Contains(mimeType, words) {
				return typeK
			}
		}
	}
	return `file`
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
