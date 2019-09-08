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

var FileTypeExts = map[string][]string{
	`image`:     []string{`jpeg`, `jpg`, `gif`, `png`},
	`flash`:     []string{`swf`},
	`audio`:     []string{`mp3`, `mid`},
	`video`:     []string{`mp4`, `mp5`, `flv`, `mpg`, `mkv`, `rmvb`, `avi`, `rm`, `asf`, `divx`, `mpeg`, `mpe`, `wmv`, `mkv`, `vob`, `3gp`, `mov`},
	`archive`:   []string{`zip`, `7z`, `rar`, `tar`, `gz`},
	`office`:    []string{`xls`, `doc`, `docx`, `ppt`, `pptx`, `et`, `wps`, `rtf`, `dps`},
	`bt`:        []string{`torrent`},
	`photoshop`: []string{`psd`},
}

var fileTypes = map[string]string{}

func InitFileTypes() {
	for typ, exts := range FileTypeExts {
		for _, ext := range exts {
			fileTypes[ext] = typ
		}
	}
}

func GetType(extName string) string {
	if len(extName) > 0 && extName[0] == '.' {
		extName = extName[1:]
	}
	if v, ok := fileTypes[extName]; ok {
		return v
	}
	return `file`
}

type FileType string

func (f FileType) ExtNames() (r []string) {
	if v, ok := FileTypeExts[string(f)]; ok {
		r = v
	}
	return
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

func init() {
	InitFileTypes()
}
