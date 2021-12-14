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

package cloud

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/nging/v4/application/registry/upload"
	"github.com/admpub/nging/v4/application/registry/upload/driver/local"
	"github.com/admpub/nging/v4/application/registry/upload/helper"
	"github.com/webx-top/echo"
	"gocloud.dev/blob"
)

//! incomplete

const Name = `cloud`

var _ upload.Storer = &Cloud{}

func init() {
	upload.StorerRegister(Name, func(ctx context.Context, subdir string) (upload.Storer, error) {
		return NewCloud(ctx, subdir)
	})
}

func NewCloud(ctx context.Context, subdir string) (*Cloud, error) {
	bucket, err := DefaultConfig.New(ctx)
	if err != nil {
		return nil, err
	}
	return &Cloud{
		config:     DefaultConfig,
		bucket:     bucket,
		Filesystem: local.NewFilesystem(ctx, subdir),
	}, nil
}

type Cloud struct {
	config *Config
	bucket *blob.Bucket
	*local.Filesystem
}

func (f *Cloud) Name() string {
	return Name
}

func (f *Cloud) filepath(fname string) string {
	return f.URLDir(fname)
}

func (f *Cloud) Exists(file string) (bool, error) {
	return f.bucket.Exists(f.Context, file)
}

func (f *Cloud) FileInfo(file string) (os.FileInfo, error) {
	r, err := f.bucket.NewReader(f.Context, file, nil)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// Access storage.Reader via sr here.
	sr := blob.Reader{}
	if r.As(&sr) {
		_ = sr.ModTime()
	}
	return os.Stat(file)
}

func (f *Cloud) SendFile(ctx echo.Context, file string) error {
	return ctx.Redirect(file)
}

func (f *Cloud) Put(dstFile string, src io.Reader, size int64) (savePath string, viewURL string, err error) {
	savePath = f.filepath(dstFile)
	viewURL = f.PublicURL(dstFile)

	var dst io.WriteCloser
	dst, err = f.bucket.NewWriter(f.Context, savePath, nil)
	if err != nil {
		return
	}
	defer dst.Close()

	//copy the uploaded file to the destination file
	_, err = io.Copy(dst, src)
	return
}

func (f *Cloud) PublicURL(dstFile string) string {
	return f.config.PublicBaseURL + f.URLDir(dstFile)
}

func (f *Cloud) URLToFile(publicURL string) string {
	dstFile := strings.TrimPrefix(publicURL, strings.TrimRight(f.PublicURL(``), `/`)+`/`)
	return dstFile
}

func (f *Cloud) FixURL(content string, embedded ...bool) string {
	return content
}

func (f *Cloud) FixURLWithParams(content string, values url.Values, embedded ...bool) string {
	if len(embedded) > 0 && embedded[0] {
		return helper.ReplaceAnyFileName(content, func(r string) string {
			return f.URLWithParams(f.PublicURL(r), values)
		})
	}
	return f.URLWithParams(f.PublicURL(content), values)
}

func (f *Cloud) Get(dstFile string) (io.ReadCloser, error) {
	return f.bucket.NewReader(f.Context, dstFile, nil)
}

func (f *Cloud) Delete(dstFile string) error {
	return f.bucket.Delete(f.Context, dstFile)
}

func (f *Cloud) DeleteDir(dstDir string) error {
	dir := filepath.Join(echo.Wd(), dstDir)
	return os.RemoveAll(dir)
}

func (f *Cloud) Close() error {
	if f.bucket != nil {
		return f.bucket.Close()
	}
	return nil
}
