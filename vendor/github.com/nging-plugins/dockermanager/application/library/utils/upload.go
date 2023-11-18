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

package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/webx-top/com"
)

func SaveUploadedFile(file *multipart.FileHeader, srcPath string) (string, error) {
	fp, err := file.Open()
	if err != nil {
		return ``, err
	}
	defer fp.Close()
	return SaveMultipartFile(fp, file.Filename, srcPath)
}

func SaveMultipartFile(fp multipart.File, fileName string, srcPath string) (string, error) {
	sf := filepath.Join(srcPath, fileName)
	sp, err := os.Create(sf)
	if err != nil {
		if !os.IsNotExist(err) {
			return sf, err
		}
		err = com.MkdirAll(filepath.Dir(sf), os.ModePerm)
		if err != nil {
			return sf, err
		}
		sp, err = os.Create(sf)
		if err != nil {
			return sf, err
		}
	}
	defer sp.Close()
	_, err = io.Copy(sp, fp)
	if err != nil {
		return sf, err
	}
	err = sp.Sync()
	return sf, err
}
