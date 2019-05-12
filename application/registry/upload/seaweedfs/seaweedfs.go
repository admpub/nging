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
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/admpub/goseaweedfs"
)

const Name = `seaweedfs`

func NewSeaweedfs(typ string) *Seaweedfs {
	a := DefaultConfig.New()
	return &Seaweedfs{
		config:   DefaultConfig,
		instance: a,
		Type:     typ,
	}
}

type Seaweedfs struct {
	config   *Config
	instance *goseaweedfs.Seaweed
	Type     string
}

func (s *Seaweedfs) Engine() string {
	return Name
}

func (f *Seaweedfs) filepath(fname string) string {
	return path.Join(f.Type, fname)
}

func (f *Seaweedfs) Put(dstFile string, src *os.File) (string, error) {
	var size int64
	fi, fiErr := src.Stat()
	if fiErr != nil {
		return "", fiErr
	}
	size = fi.Size()
	if len(dstFile) == 0 {
		dstFile = src.Name()
	}
	_, fID, err := f.instance.Upload(src, dstFile, size, f.Type, f.config.TTL)
	if err != nil {
		return "", err
	}
	view, err := f.instance.LookupFileID(fID, url.Values{}, true)
	if err != nil {
		return view, err
	}
	return view, nil
}

func (f *Seaweedfs) Get(fileID string) (io.ReadCloser, error) {
	_, fileData, err := f.instance.DownloadFile(fileID, nil)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewBuffer(fileData)), nil
}

func (f *Seaweedfs) Delete(fileID string) error {
	return f.instance.DeleteFile(fileID, nil)
}

func (f *Seaweedfs) DeleteDir(dstDir string) error {
	if len(f.instance.Filers) > 0 {
		return f.instance.Filers[0].Delete(dstDir, true)
	}
	return nil
}
