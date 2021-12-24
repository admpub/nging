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
	"os"
	"path/filepath"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// WriteCache 写缓存文件
func WriteCache(dir string, name string, content []byte) (err error) {
	savePath := filepath.Join(echo.Wd(), `data`, `cache`, dir)
	err = com.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		return
	}
	err = os.WriteFile(savePath+echo.FilePathSeparator+name, content, os.ModePerm)
	return
}

// ReadCache 读缓存文件
func ReadCache(dir string, name string) (content []byte, err error) {
	savePath := filepath.Join(echo.Wd(), `data`, `cache`, dir, name)
	return os.ReadFile(savePath)
}

// ModTimeCache 缓存文件修改时间
func ModTimeCache(dir string, name string) (time.Time, error) {
	savePath := filepath.Join(echo.Wd(), `data`, `cache`, dir, name)
	info, err := os.Stat(savePath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), err
}

// RemoveCache 删除缓存文件
func RemoveCache(dir string, names ...string) (err error) {
	savePath := filepath.Join(echo.Wd(), `data`, `cache`, dir)
	if len(names) < 1 {
		if _, err = os.Stat(savePath); !os.IsNotExist(err) {
			return os.RemoveAll(savePath)
		}
		return
	}
	for _, name := range names {
		filePath := savePath + echo.FilePathSeparator + name
		if _, err = os.Stat(filePath); !os.IsNotExist(err) {
			err = os.Remove(filePath)
			if err != nil {
				return
			}
		}
	}
	return err
}
