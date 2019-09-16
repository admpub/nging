/*

   Copyright 2016-present Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package upload

import (
	"mime"
	"strings"

	"github.com/webx-top/com"
)

// FileType 文件类型名
type FileType string

// Extensions 文件类型相应扩展名
func (f FileType) Extensions() (r []string) {
	if v, ok := FileTypeExts[f]; ok {
		r = v
	}
	return
}

func (f FileType) String() string {
	return string(f)
}

// Icon 文件类型 Font Awesome 图标(不含“fa fa-”)
func (f FileType) Icon() string {
	return FileTypeIcon(f.String())
}

const (
	TypeImage     FileType = `image`
	TypeFlash     FileType = `flash`
	TypeAudio     FileType = `audio`
	TypeVideo     FileType = `video`
	TypeArchive   FileType = `archive`
	TypeOffice    FileType = `office`
	TypeBT        FileType = `bt`
	TypePhotoshop FileType = `photoshop`
	TypePDF       FileType = `pdf`
)

var (
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

	// FileTypeMimeKeywords 文件类型对应mime关键词
	FileTypeMimeKeywords = map[string][]string{
		`image`:   []string{`image`},
		`video`:   []string{`video`},
		`audio`:   []string{`audio`},
		`archive`: []string{`compressed`},
		`pdf`:     []string{`pdf`},
		`xls`:     []string{`csv`, `excel`},
		`ppt`:     []string{`powerpoint`},
		`doc`:     []string{`msword`, `text`},
	}

	// FileTypeExts 文件类型对应扩展名(不含".")
	FileTypeExts = map[FileType][]string{
		TypeImage:   []string{`jpeg`, `jpg`, `gif`, `png`, `svg`, `webp`},
		TypeFlash:   []string{`swf`},
		TypeAudio:   []string{`mp3`, `mid`},
		TypeVideo:   []string{`mp4`, `mp5`, `flv`, `mpg`, `mkv`, `rmvb`, `avi`, `rm`, `asf`, `divx`, `mpeg`, `mpe`, `wmv`, `mkv`, `vob`, `3gp`, `mov`},
		TypeArchive: []string{`zip`, `7z`, `rar`, `tar`, `gz`, `bz`, `gzip`},
		TypeOffice: []string{
			`xls`, `xlsx`, `csv`, //xls
			`doc`, `docx`, //doc
			`ppt`, `pptx`, //ppt
			`et`, `wps`, `rtf`, `dps`,
		},
		TypeBT:        []string{`torrent`},
		TypePhotoshop: []string{`psd`},
		TypePDF:       []string{`pdf`},
	}

	// 扩展名对应类型
	fileTypes = map[string]FileType{}
)

// TypeRegister 文件类型注册
func TypeRegister(fileType FileType, extensions ...string) {
	if _, ok := FileTypeExts[fileType]; !ok {
		FileTypeExts[fileType] = []string{}
	}
	for _, extension := range extensions {
		if len(extension) > 0 && extension[0] == '.' {
			extension = extension[1:]
		}
		extension = strings.ToLower(extension)
		if _, ok := fileTypes[extension]; ok {
			continue
		}
		if com.InSlice(extension, FileTypeExts[fileType]) {
			continue
		}
		FileTypeExts[fileType] = append(FileTypeExts[fileType], extension)
		fileTypes[extension] = fileType
	}
}

// InitFileTypes 初始化文件扩展名与类型对应关系
func InitFileTypes() {
	for fileType, extensions := range FileTypeExts {
		for _, extension := range extensions {
			fileTypes[extension] = fileType
		}
	}
}

// DetectType 根据扩展名判断类型
func DetectType(extension string) string {
	if len(extension) > 0 && extension[0] == '.' {
		extension = extension[1:]
	}
	extension = strings.ToLower(extension)
	if v, ok := fileTypes[extension]; ok {
		return v.String()
	}
	mimeType := mime.TypeByExtension(`.` + extension)
	mimeType = strings.SplitN(mimeType, ";", 2)[0]
	for typeK, keywords := range FileTypeMimeKeywords {
		for _, words := range keywords {
			if strings.Contains(mimeType, words) {
				return typeK
			}
		}
	}
	return `file`
}

// FileTypeIcon 文件类型icon
func FileTypeIcon(typ string) string {
	icon, ok := FileTypeIcons[typ]
	if ok {
		return icon
	}
	return `file-o`
}

func init() {
	InitFileTypes()
}
