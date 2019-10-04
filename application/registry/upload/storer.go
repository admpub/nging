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
	"context"
	"io"
	"mime/multipart"
	"net/url"
	"os"

	"github.com/admpub/checksum"
	"github.com/admpub/nging/application/registry/upload/table"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
)

var (
	ErrExistsFile = table.ErrExistsFile
)

func BatchUpload(
	ctx echo.Context,
	fieldName string,
	dstNamer func(*uploadClient.Result) (dst string, err error),
	storer Storer,
	callback func(*uploadClient.Result, multipart.File) error,
) (results uploadClient.Results, err error) {
	req := ctx.Request()
	if req == nil {
		err = ctx.E(`Invalid upload content`)
		return
	}
	m := req.MultipartForm()
	if m == nil || m.File == nil {
		err = ctx.E(`Invalid upload content`)
		return
	}
	files, ok := m.File[fieldName]
	if !ok {
		err = echo.ErrNotFoundFileInput
		return
	}
	var dstFile string
	for _, fileHdr := range files {
		//for each fileheader, get a handle to the actual file
		var file multipart.File
		file, err = fileHdr.Open()
		if err != nil {
			if file != nil {
				file.Close()
			}
			return
		}
		result := &uploadClient.Result{
			FileName: fileHdr.Filename,
			FileSize: fileHdr.Size,
		}
		result.Md5, err = checksum.MD5sumReader(file)
		if err != nil {
			file.Close()
			return
		}

		dstFile, err = dstNamer(result)
		if err != nil {
			file.Close()
			if err == ErrExistsFile {
				results.Add(result)
				err = nil
				continue
			}
			return
		}
		if len(dstFile) == 0 {
			file.Close()
			continue
		}
		if len(result.SavePath) > 0 {
			file.Close()
			results.Add(result)
			continue
		}
		file.Seek(0, 0)
		result.SavePath, result.FileURL, err = storer.Put(dstFile, file, fileHdr.Size)
		if err != nil {
			file.Close()
			return
		}
		file.Seek(0, 0)
		if err = callback(result, file); err != nil {
			file.Close()
			return
		}
		file.Close()
		results.Add(result)
	}
	return
}

type Sizer interface {
	Size() int64
}

type Storer interface {
	// 引擎名
	Name() string

	// FileDir 文件夹物理路径
	FileDir(subpath string) string

	// URLDir 文件夹网址路径
	URLDir(subpath string) string

	// Put 保存文件
	Put(dst string, src io.Reader, size int64) (savePath string, viewURL string, err error)

	// Get 获取文件
	Get(file string) (io.ReadCloser, error)

	// Exists 文件是否存在
	Exists(file string) (bool, error)

	// FileInfo 文件信息
	FileInfo(file string) (os.FileInfo, error)

	// SendFile 输出文件到浏览器
	SendFile(ctx echo.Context, file string) error

	// Delete 删除文件
	Delete(file string) error

	// DeleteDir 删除目录
	DeleteDir(dir string) error

	// PublicURL 文件网址
	PublicURL(dst string) string

	// URLToFile 网址转文件物理路径
	URLToFile(viewURL string) string

	// FixURL 修正网址
	FixURL(content string, embedded ...bool) string

	// FixURLWithParams 修正网址并增加网址参数
	FixURLWithParams(content string, values url.Values, embedded ...bool) string

	// Close 关闭连接
	Close() error
}

type Constructor func(ctx context.Context, typ string) Storer

var storers = map[string]Constructor{}

var DefaultConstructor Constructor

func StorerRegister(engine string, constructor Constructor) {
	storers[engine] = constructor
}

func StorerGet(engine string) Constructor {
	constructor, ok := storers[engine]
	if !ok {
		return DefaultConstructor
	}
	return constructor
}

func StorerAll(engine string) map[string]Constructor {
	return storers
}
