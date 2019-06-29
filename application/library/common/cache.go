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

package common

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/webx-top/echo"
)

// WriteCache 写缓存文件
func WriteCache(dir string, name string, content []byte) (err error) {
	savePath := filepath.Join(echo.Wd(), `data`, `cache`, dir)
	if _, err = os.Stat(savePath); os.IsNotExist(err) {
		err = os.MkdirAll(savePath, os.ModePerm)
		if err != nil {
			return
		}
	}
	err = ioutil.WriteFile(savePath+echo.FilePathSeparator+name, content, os.ModePerm)
	return
}

// ReadCache 读缓存文件
func ReadCache(dir string, name string) (content []byte, err error) {
	savePath := filepath.Join(echo.Wd(), `data`, `cache`, dir, name)
	return ioutil.ReadFile(savePath)
}

// RemoveCache 删除缓存文件
func RemoveCache(dir string, names ...string) (err error) {
	savePath := filepath.Join(echo.Wd(), `data`, `cache`, dir)
	if len(names) < 1 {
		if _, err = os.Stat(savePath); os.IsExist(err) {
			return os.RemoveAll(savePath)
		}
		return
	}
	for _, name := range names {
		filePath := savePath + echo.FilePathSeparator + name
		if _, err = os.Stat(filePath); os.IsExist(err) {
			err = os.Remove(filePath)
			if err != nil {
				return
			}
		}
	}
	return err
}
