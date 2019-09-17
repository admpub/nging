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
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/admpub/checksum"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

var (
	ErrUndefinedFileName           = errors.New(`Undefined file name`)
	ErrUndefinedContentDisposition = errors.New(`Not found Content-Disposition in header`)
)

type Sizer interface {
	Size() int64
}

type SetReadCloser interface {
	SetReadCloser(rc io.ReadCloser)
}

type ReadCloserWithSize interface {
	Sizer
	io.ReadCloser
	Md5() (string, error)
}

func WrapBodyWithSize(req engine.Request) ReadCloserWithSize {
	return &wrapBodyWithSize{Request: req}
}

func WrapFileWithSize(size int64, file multipart.File) ReadCloserWithSize {
	return &wrapFileWithSize{size: size, File: file}
}

type wrapBodyWithSize struct {
	engine.Request
}

func (w *wrapBodyWithSize) Read(p []byte) (n int, err error) {
	return w.Body().Read(p)
}

func (w *wrapBodyWithSize) Close() error {
	return w.Body().Close()
}

func (w *wrapBodyWithSize) Md5() (md5 string, err error) {
	var b []byte
	b, err = ioutil.ReadAll(w.Body())
	if err != nil {
		return
	}
	w.Body().Close()
	defer func() {
		w.SetBody(bytes.NewReader(b))
	}()
	md5, err = checksum.MD5sumReader(bytes.NewReader(b))
	return
}

type wrapFileWithSize struct {
	size int64
	multipart.File
}

func (w *wrapFileWithSize) Size() int64 {
	return w.size
}

func (w *wrapFileWithSize) Md5() (md5 string, err error) {
	md5, err = checksum.MD5sumReader(w)
	if err != nil {
		return
	}
	_, err = w.Seek(0, 0)
	return
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
		f = WrapBodyWithSize(ctx.Request())
		return

	default:
		var header *multipart.FileHeader
		var file multipart.File
		file, header, err = ctx.Request().FormFile(name)
		if err != nil {
			return
		}
		fileName = header.Filename
		f = WrapFileWithSize(header.Size, file)
		return
	}
}
