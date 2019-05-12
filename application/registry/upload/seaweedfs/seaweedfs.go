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
	"os"
	"path/filepath"

	"github.com/webx-top/echo"
)

const Name = `seaweedfs`

func NewSeaweedfs(typ string) *Seaweedfs {
	uploadPath := `public/upload/` + typ
	return &Seaweedfs{
		Type:       typ,
		UploadPath: uploadPath,
	}
}

type Seaweedfs struct {
	Type       string
	UploadPath string
}

func (s *Seaweedfs) Engine() string {
	return Name
}

func (f *Seaweedfs) filepath(fname string) string {
	return filepath.Join(echo.Wd(), f.UploadPath, fname)
}

func (f *Seaweedfs) Put(dstFile string, src io.Reader) (string, error) {
	file := f.filepath(dstFile)
	view := `/` + f.UploadPath + `/` + dstFile
	//create destination file making sure the path is writeable.
	dst, err := os.Create(file)
	if err != nil {
		return view, err
	}
	defer dst.Close()
	//copy the uploaded file to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		return view, err
	}
	return view, nil
}

func (f *Seaweedfs) Get(dstFile string) (io.ReadCloser, error) {
	return f.OpenFile(dstFile)
}

func (f *Seaweedfs) OpenFile(dstFile string) (*os.File, error) {
	//file := f.filepath(dstFile)
	file := filepath.Join(echo.Wd(), dstFile)
	return os.Open(file)
}

func (f *Seaweedfs) Delete(dstFile string) error {
	file := filepath.Join(echo.Wd(), dstFile)
	return os.Remove(file)
}

func (f *Seaweedfs) DeleteDir(dstDir string) error {
	dir := filepath.Join(echo.Wd(), dstDir)
	return os.RemoveAll(dir)
}
