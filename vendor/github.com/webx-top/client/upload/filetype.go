/*

   Copyright 2016 Wenhui Shen <www.webx.top>

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

import "strings"

type FileType string

func (f FileType) Extensions() (r []string) {
	if v, ok := FileTypeExts[f]; ok {
		r = v
	}
	return
}

func (f FileType) String() string {
	return string(f)
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
)

var FileTypeExts = map[FileType][]string{
	TypeImage:     []string{`jpeg`, `jpg`, `gif`, `png`},
	TypeFlash:     []string{`swf`},
	TypeAudio:     []string{`mp3`, `mid`},
	TypeVideo:     []string{`mp4`, `mp5`, `flv`, `mpg`, `mkv`, `rmvb`, `avi`, `rm`, `asf`, `divx`, `mpeg`, `mpe`, `wmv`, `mkv`, `vob`, `3gp`, `mov`},
	TypeArchive:   []string{`zip`, `7z`, `rar`, `tar`, `gz`},
	TypeOffice:    []string{`xls`, `doc`, `docx`, `ppt`, `pptx`, `et`, `wps`, `rtf`, `dps`},
	TypeBT:        []string{`torrent`},
	TypePhotoshop: []string{`psd`},
}

func AddFileType(fileType FileType, extensions ...string) {
	if _, y := FileTypeExts[fileType]; !y {
		FileTypeExts[fileType] = []string{}
	}
	FileTypeExts[fileType] = append(FileTypeExts[fileType], extensions...)
	for _, extension := range extensions {
		fileTypes[extension] = fileType
	}
}

var fileTypes = map[string]FileType{}

func InitFileTypes() {
	for fileType, extensions := range FileTypeExts {
		for _, extension := range extensions {
			fileTypes[extension] = fileType
		}
	}
}

func DetectType(extension string) string {
	if len(extension) > 0 && extension[0] == '.' {
		extension = extension[1:]
	}
	extension = strings.ToLower(extension)
	if v, ok := fileTypes[extension]; ok {
		return v.String()
	}
	return `file`
}

func init() {
	InitFileTypes()
}
