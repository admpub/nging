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

import (
	"errors"
	"io"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/webx-top/echo"
)

var (
	ErrUndefinedFileName           = errors.New(`Undefined file name`)
	ErrUndefinedContentDisposition = errors.New(`Not found Content-Disposition in header`)
)

type Sizer interface {
	Size() int64
}

type ReadCloserWithSize interface {
	Sizer
	io.ReadCloser
}

func WrapSizer(sizer Sizer, readCloser io.ReadCloser) ReadCloserWithSize {
	return &wrapSizer{sizer: sizer, ReadCloser: readCloser}
}

func WrapSize(size int64, readCloser io.ReadCloser) ReadCloserWithSize {
	return &wrapSize{size: size, ReadCloser: readCloser}
}

type wrapSizer struct {
	sizer Sizer
	io.ReadCloser
}

func (w *wrapSizer) Size() int64 {
	return w.sizer.Size()
}

type wrapSize struct {
	size int64
	io.ReadCloser
}

func (w *wrapSize) Size() int64 {
	return w.size
}

func Receive(name string, ctx echo.Context) (f ReadCloserWithSize, fileName string, err error) {
	switch ctx.ResolveContentType() {
	case "application/octet-stream":
		val := ctx.Request().Header().Get("Content-Disposition")
		if len(val) == 0 {
			return nil, ``, ErrUndefinedContentDisposition
		}
		fileNameMark := `; filename="`
		pos := strings.LastIndex(val, fileNameMark)
		if pos < 0 {
			return nil, ``, ErrUndefinedFileName
		}
		fileName = val[pos+len(fileNameMark):]
		fileName = strings.TrimRight(fileName, `"`)
		fileName, err = url.QueryUnescape(fileName)
		if err != nil {
			return
		}
		f = WrapSizer(ctx.Request(), ctx.Request().Body())
		return

	default:
		var header *multipart.FileHeader
		var file multipart.File
		file, header, err = ctx.Request().FormFile(name)
		if err != nil {
			return
		}
		fileName = header.Filename
		f = WrapSize(header.Size, file)
		return
	}
}
