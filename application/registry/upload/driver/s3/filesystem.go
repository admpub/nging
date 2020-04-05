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

package s3

import (
	"context"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/engine/mock"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/application/registry/upload"
	"github.com/admpub/nging/application/registry/upload/helper"
	"github.com/admpub/nging/application/registry/upload/driver/local"
	"github.com/admpub/nging/application/library/s3manager"
	"github.com/admpub/nging/application/library/s3manager/s3client"
	"github.com/admpub/nging/application/model"
	"github.com/admpub/nging/application/model/file/storer"
)

const (
	Name = `s3`
	AccountIDKey = `storerID`
)

var _ upload.Storer = &Filesystem{}

func init() {
	upload.StorerRegister(Name, func(ctx context.Context, typ string) (upload.Storer, error) {
		return NewFilesystem(ctx, typ)
	})
}

func NewFilesystem(ctx context.Context, typ string) (*Filesystem, error) {
	var cloudAccountID string
	eCtx, ok := ctx.(echo.Context)
	if !ok {
		eCtx = echo.NewContext(mock.NewRequest(), mock.NewResponse(), defaults.Default)
		eCtx.SetStdContext(ctx)
	} else {
		cloudAccountID = eCtx.Internal().String(AccountIDKey)
	}
	m := model.NewCloudStorage(eCtx)
	if len(cloudAccountID) == 0 {
		cloudAccountID = param.AsString(ctx.Value(AccountIDKey))
	}
	if len(cloudAccountID) == 0 || cloudAccountID == `0` {
		storerConfig, ok := storer.GetOk()
		if ok {
			cloudAccountID = storerConfig.ID
		}
	}
	if err := m.Get(nil, `id`, cloudAccountID); err != nil {
		return nil, err
	}
	mgr, err := s3client.New(m.NgingCloudStorage, 0)
	if err != nil {
		return nil, err
	}
	return &Filesystem{
		Filesystem: local.NewFilesystem(ctx, typ),
		cdn: strings.TrimSuffix(m.Baseurl, `/`),
		model: m,
		mgr: mgr,
	}, nil
}

// Filesystem 文件系统存储引擎
type Filesystem struct {
	*local.Filesystem
	cdn string // CDN URL
	model *model.CloudStorage
	mgr *s3manager.S3Manager
}

// Name 引擎名
func (f *Filesystem) Name() string {
	return Name
}

// Exists 判断文件是否存在
func (f *Filesystem) Exists(file string) (bool, error) {
	return f.mgr.Exists(file)
}

// FileInfo 获取文件信息
func (f *Filesystem) FileInfo(file string) (os.FileInfo, error) {
	objectInfo, err := f.mgr.Stat(file)
	if err != nil {
		return nil, err
	}
	return s3manager.NewFileInfo(objectInfo), err
}

// SendFile 下载文件
func (f *Filesystem) SendFile(ctx echo.Context, file string) error {
	fp, err := f.mgr.Get(file)
	if err != nil {
		return err
	}
	defer fp.Close()
	fileName := path.Base(file)
	inline := true
	return ctx.Attachment(fp, fileName, inline)
}

// FileDir 物理路径文件夹
func (f *Filesystem) FileDir(subpath string) string {
	return path.Join(helper.UploadURLPath, f.Type, subpath)
}

// Put 上传文件
func (f *Filesystem) Put(dstFile string, src io.Reader, size int64) (savePath string, viewURL string, err error) {
	savePath = f.FileDir(dstFile)
	//viewURL = `[storage:`+param.AsString(f.model.Id)+`]`+f.URLDir(dstFile)
	viewURL = f.PublicURL(dstFile)
	err = f.mgr.Put(src, savePath, size)
	return
}

// Get 获取文件读取接口
func (f *Filesystem) Get(dstFile string) (io.ReadCloser, error) {
	return f.mgr.Get(dstFile)
}

// Delete 删除文件
func (f *Filesystem) Delete(dstFile string) error {
	return f.mgr.Remove(dstFile)
}

// DeleteDir 删除文件夹及其内部文件
func (f *Filesystem) DeleteDir(dstDir string) error {
	return f.mgr.Remove(dstDir)
}

// Move 移动文件
func (f *Filesystem) Move(src, dst string) error {
	return f.mgr.Rename(src, dst)
}

// Close 关闭连接
func (f *Filesystem) Close() error {
	return nil
}

// PublicURL 文件物理路径转为文件网址
func (f *Filesystem) PublicURL(dstFile string) string {
	return f.cdn + f.URLDir(dstFile)
}

// URLToFile 文件网址转为文件物理路径
func (f *Filesystem) URLToFile(publicURL string) string {
	dstFile := strings.TrimPrefix(publicURL, strings.TrimRight(f.PublicURL(``), `/`)+`/`)
	return dstFile
}

// FixURL 改写文件网址
func (f *Filesystem) FixURL(content string, embedded ...bool) string {
	rowsByID := f.model.CachedList()  
	return helper.ReplacePlaceholder(content, func(id string) string {
		r, y := rowsByID[id]
		if !y {
			return ``
		}
		return r.Baseurl
	})
}

// FixURLWithParams 替换文件网址为带参数的网址
func (f *Filesystem) FixURLWithParams(content string, values url.Values, embedded ...bool) string {
	if len(embedded) > 0 && embedded[0] {
		return helper.ReplaceAnyFileName(content, func(r string) string {
			return f.URLWithParams(f.PublicURL(r), values)
		})
	}
	return f.URLWithParams(f.PublicURL(content), values)
}
