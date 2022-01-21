/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
	"context"
	"io"
	"net/url"

	"github.com/admpub/errors"
	"github.com/admpub/goseaweedfs"
	"github.com/admpub/nging/v4/application/model"
	"github.com/admpub/nging/v4/application/registry/upload"
	"github.com/admpub/nging/v4/application/registry/upload/driver/local"
	"github.com/admpub/nging/v4/application/registry/upload/helper"
)

const Name = `seaweedfs`

var _ upload.Storer = &Seaweedfs{}

func init() {
	upload.StorerRegister(Name, func(ctx context.Context, subdir string) (upload.Storer, error) {
		return NewSeaweedfs(ctx, subdir)
	})
}

func NewSeaweedfs(ctx context.Context, subdir string) (*Seaweedfs, error) {
	m, err := model.GetCloudStorage(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	a, err := DefaultConfig.New()
	if err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	return &Seaweedfs{
		config:     DefaultConfig,
		model:      m,
		instance:   a,
		Filesystem: local.NewFilesystem(ctx, subdir),
	}, nil
}

type Seaweedfs struct {
	config   *Config
	model    *model.CloudStorage
	instance *goseaweedfs.Seaweed
	*local.Filesystem
}

func (s *Seaweedfs) Name() string {
	return Name
}

func (s *Seaweedfs) filepath(fname string) string {
	return s.URLDir(fname)
}

func (s *Seaweedfs) xPut(dstFile string, src io.Reader, size int64) (savePath string, viewURL string, err error) {
	savePath = s.filepath(dstFile)
	var rs *goseaweedfs.FilerUploadResult
	rs, err = s.instance.Filers()[0].Upload(src, size, savePath, s.Subdir, s.config.TTL)
	if err != nil {
		err = errors.WithMessage(err, Name)
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

/*
func (s *Seaweedfs) xGet(dstFile string) (io.ReadCloser, error) {
	filer := s.instance.Filers()[0]
	_, readCloser, err := filer.Get(dstFile, nil, nil)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return readCloser, err
}
*/

func (s *Seaweedfs) PublicURL(dstFile string) string {
	return s.config.Filers[0].Public + dstFile
}

func (f *Seaweedfs) FixURL(content string, embedded ...bool) string {
	rowsByID := f.model.CachedList()
	return helper.ReplacePlaceholder(content, func(id string) string {
		r, y := rowsByID[id]
		if !y {
			return ``
		}
		if len(r.Baseurl) > 0 {
			return r.Baseurl
		}
		return f.PublicURL(``)
	})
}

func (s *Seaweedfs) xDelete(dstFile string) error {
	filer := s.instance.Filers()[0]
	err := filer.Delete(dstFile, nil)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return err
}

func (s *Seaweedfs) xDeleteDir(dstDir string) error {
	err := s.instance.Filers()[0].Delete(dstDir, nil)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return err
}

func (s *Seaweedfs) apiPut(dstFile string, src io.Reader, size int64) (fID string, viewURL string, err error) {
	var part *goseaweedfs.FilePart
	part, err = s.instance.Upload(src, dstFile, size, s.Subdir, s.config.TTL)
	if err != nil {
		return
	}

	viewURL, err = s.instance.LookupFileID(part.FileID, url.Values{}, true)
	return
}

/*
func (s *Seaweedfs) apiGet(fileID string) (io.ReadCloser, error) {
	_, readCloser, err := s.instance.Download(fileID, nil)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return readCloser, err
}
*/

func (s *Seaweedfs) apiDelete(fileID string) error {
	err := s.instance.DeleteFile(fileID, nil)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return err
}
