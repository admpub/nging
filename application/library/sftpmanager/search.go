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

package sftpmanager

import (
	"path"
	"strings"

	"github.com/pkg/sftp"
)

func Search(client *sftp.Client, query string, typ string, nums ...int) []string {
	var (
		paths  []string
		prefix string
		ppath  string
	)
	if strings.HasSuffix(query, `/`) {
		ppath = query
	} else {
		prefix = path.Base(query)
		ppath = path.Dir(query)
	}
	num := 10
	if len(nums) > 0 {
		num = nums[0]
		if num <= 0 {
			num = 10
		}
	}
	if len(ppath) == 0 {
		ppath = `/`
	}
	var onlyDir bool
	switch typ {
	case `dir`:
		onlyDir = true
	case `file`:
		onlyDir = false
	default:
		onlyDir = true
	}
	dirs, _ := client.ReadDir(ppath)
	for _, d := range dirs {
		if onlyDir && d.IsDir() == false {
			continue
		}
		if len(paths) >= num {
			break
		}
		name := d.Name()
		if len(prefix) == 0 || strings.HasPrefix(name, prefix) {
			paths = append(paths, path.Join(ppath, name)+`/`)
			continue
		}
	}
	return paths
}
