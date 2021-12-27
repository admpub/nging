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

package filemanager

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/webx-top/echo"
)

// Search 自动完成查询文件
func Search(prefix string, nums ...int) []string {
	var paths []string
	root := filepath.Dir(prefix)
	if len(root) == 0 {
		root = echo.Wd()
	}
	num := 10
	if len(nums) > 0 {
		num = nums[0]
	}
	dir, _ := os.ReadDir(root)
	for _, d := range dir {
		if len(paths) >= num {
			break
		}
		path := filepath.Join(root, d.Name())
		if d.IsDir() {
			path += echo.FilePathSeparator
		}
		if len(prefix) == 0 {
			paths = append(paths, path)
			continue
		}
		if strings.HasPrefix(path, prefix) {
			paths = append(paths, path)
			continue
		}
	}
	return paths
}

// SearchDir 自动完成查询文件
func SearchDir(prefix string, nums ...int) []string {
	var paths []string
	root := filepath.Dir(prefix)
	if len(root) == 0 {
		root = echo.Wd()
	}
	num := 10
	if len(nums) > 0 {
		num = nums[0]
	}
	dir, _ := os.ReadDir(root)
	for _, d := range dir {
		if len(paths) >= num {
			break
		}
		if !d.IsDir() {
			continue
		}
		path := filepath.Join(root, d.Name())
		path += echo.FilePathSeparator
		if len(prefix) == 0 {
			paths = append(paths, path)
			continue
		}
		if strings.HasPrefix(path, prefix) {
			paths = append(paths, path)
			continue
		}
	}
	return paths
}

// SearchFile 自动完成查询文件
func SearchFile(prefix string, nums ...int) []string {
	var paths []string
	root := filepath.Dir(prefix)
	if len(root) == 0 {
		root = echo.Wd()
	}
	num := 10
	if len(nums) > 0 {
		num = nums[0]
	}
	dir, _ := os.ReadDir(root)
	for _, d := range dir {
		if len(paths) >= num {
			break
		}
		if d.IsDir() {
			continue
		}
		path := filepath.Join(root, d.Name())
		if len(prefix) == 0 {
			paths = append(paths, path)
			continue
		}
		if strings.HasPrefix(path, prefix) {
			paths = append(paths, path)
			continue
		}
	}
	return paths
}
