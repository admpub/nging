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

package seaweedfs

import (
	"io"
	"net/url"
	"path"
	"context"

	"github.com/admpub/nging/application/registry/upload/driver/filesystem"
	"github.com/admpub/nging/application/registry/upload/helper"

	"github.com/admpub/goseaweedfs"
	modelSeaweedfs "github.com/admpub/goseaweedfs/model"
	"github.com/admpub/nging/application/registry/upload"
)

const Name = `seaweedfs`

var _ upload.Storer = &Seaweedfs{}

func init() {
	upload.StorerRegister(Name, func(ctx context.Context, typ string) upload.Storer {
		return NewSeaweedfs(ctx, typ)
	})
}

func NewSeaweedfs(ctx context.Context, typ string) *Seaweedfs {
	a := DefaultConfig.New()
	return &Seaweedfs{
		config:     DefaultConfig,
		instance:   a,
		Filesystem: filesystem.NewFilesystem(ctx, typ),
	}
}

type Seaweedfs struct {
	config   *Config
	instance *goseaweedfs.Seaweed
	*filesystem.Filesystem
}

func (s *Seaweedfs) Engine() string {
	return Name
}

func (s *Seaweedfs) filepath(fname string) string {
	return path.Join(s.UploadDir, fname)
}

func (s *Seaweedfs) xPut(dstFile string, src io.Reader, size int64) (savePath string, viewURL string, err error) {
	savePath = s.filepath(dstFile)
	var rs *modelSeaweedfs.FilerUploadResult
	rs, err = s.instance.Filers[0].Upload(src, size, savePath, s.Type, s.config.TTL)
	if err != nil {
		return
	}
	//com.Dump(rs)
	// {
	//   "name": "config.go",
	//   "url": "http://127.0.0.1:9001/6,070894a14c",
	//   "fid": "6,070894a14c",
	//   "size": 1734
	// }

	viewURL = rs.FileID //TODO: fileID VS filePath
	//viewURL = s.instance.Filers[0]+savePath
	return
}

func (s *Seaweedfs) xGet(dstFile string) (io.ReadCloser, error) {
	filer := s.instance.Filers[0]
	_, readCloser, err := filer.Download(dstFile)
	return readCloser, err
}

func (s *Seaweedfs) PublicURL(dstFile string) string {
	return s.config.Filers[0].Public + dstFile
}

func (f *Seaweedfs) FixURL(content string, embedded ...bool) string {
	if len(embedded) > 0 && embedded[0] {
		return helper.ReplaceAnyFileName(content, func(r string) string {
			return f.PublicURL(r)
		})
	}
	return f.PublicURL(content)
}

func (f *Seaweedfs) FixURLWithParams(content string, values url.Values, embedded ...bool) string {
	if len(embedded) > 0 && embedded[0] {
		return helper.ReplaceAnyFileName(content, func(r string) string {
			return f.URLWithParams(f.PublicURL(r), values)
		})
	}
	return f.URLWithParams(f.PublicURL(content), values)
}

func (s *Seaweedfs) xDelete(dstFile string) error {
	filer := s.instance.Filers[0]
	return filer.Delete(dstFile)
}

func (s *Seaweedfs) xDeleteDir(dstDir string) error {
	return s.instance.Filers[0].Delete(dstDir, true)
}

func (s *Seaweedfs) apiPut(dstFile string, src io.Reader, size int64) (fID string, viewURL string, err error) {
	_, fID, err = s.instance.Upload(src, dstFile, size, s.Type, s.config.TTL)
	if err != nil {
		return
	}
	viewURL, err = s.instance.LookupFileID(fID, url.Values{}, true)
	return
}

func (s *Seaweedfs) apiGet(fileID string) (io.ReadCloser, error) {
	_, readCloser, err := s.instance.Download(fileID, nil)
	return readCloser, err
}

func (s *Seaweedfs) apiDelete(fileID string) error {
	return s.instance.DeleteFile(fileID, nil)
}
