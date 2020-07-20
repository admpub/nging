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
   along with f program.  If not, see <https://www.gnu.org/licenses/>.
*/

package list

import (
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/registry/upload/subdir"
	uploadSubdir "github.com/admpub/nging/application/registry/upload/subdir"
)

type FileInfo struct {
	*dbschema.NgingFile
	SubdirInfo *uploadSubdir.SubdirInfo `db:"-"`
}

var FileList = func(list []*dbschema.NgingFile) interface{} {
	listData := make([]*FileInfo, len(list))
	for k, v := range list {
		listData[k] = &FileInfo{
			NgingFile:  v,
			SubdirInfo: subdir.GetByTable(v.TableName),
		}
	}
	return listData
}
