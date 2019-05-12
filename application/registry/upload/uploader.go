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

package upload

import (
	"io"
	"mime/multipart"

	"github.com/webx-top/echo"
)

func BatchUpload(
	ctx echo.Context,
	fieldName string,
	dstNamer func(*multipart.FileHeader) (dst string, err error),
	uploader Uploader,
) ([]string, error) {
	m := ctx.Request().MultipartForm()
	files, ok := m.File[fieldName]
	if !ok {
		return nil, echo.ErrNotFoundFileInput
	}
	var dstFile string
	var viewURLs []string
	for _, fileHdr := range files {
		//for each fileheader, get a handle to the actual file
		file, err := fileHdr.Open()
		if err != nil {
			file.Close()
			return viewURLs, err
		}
		dstFile, err = dstNamer(fileHdr)
		if err != nil {
			file.Close()
			return viewURLs, err
		}
		if len(dstFile) == 0 {
			file.Close()
			continue
		}
		viewURL, err := uploader.Put(dstFile, file)
		if err != nil {
			file.Close()
			return viewURLs, err
		}
		file.Close()
		viewURLs = append(viewURLs, viewURL)
	}
	return viewURLs, nil
}

type Uploader interface {
	Put(dst string, src io.Reader) (string, error)
}

type Constructor func(typ string) Uploader

var uploaders = map[string]Constructor{}

func UploaderRegister(engine string, constructor Constructor) {
	uploaders[engine] = constructor
}

func UploaderGet(engine string) Constructor {
	constructor, ok := uploaders[engine]
	if !ok {
		return nil
	}
	return constructor
}

func UploaderAll(engine string) map[string]Constructor {
	return uploaders
}
