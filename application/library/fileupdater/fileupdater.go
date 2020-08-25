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

// Package fileupdater 编辑内容时更新相关引用文件的关联关系
// 用法：fileupdater.New(fileModel.NewEmbedded(ctx)).Set(`表名称`,`字段名称`,主键ID).Add(`/test/image.jpg`,false)
package fileupdater

import (
	"github.com/admpub/log"
	"github.com/webx-top/echo"

	uploadHelper "github.com/admpub/nging/application/registry/upload/helper"
)

var Debug = false

func New(reler Reler) *FileUpdater {
	return &FileUpdater{
		rel: reler,
	}
}

type FileUpdater struct {
	rel       Reler
	project   string
	table     string
	field     string
	tableID   string
	seperator string
}

func (f *FileUpdater) Handle(event string, content *string, embedded bool) error {
	switch event {
	case `creating`, `created`:
		return f.Add(content, embedded)
	case `updating`, `updated`:
		return f.Edit(content, embedded)
	case `deleting`, `deleted`:
		return f.Delete()
	default:
		panic(`Unsupported event: ` + event)
	}
}

func (f *FileUpdater) Add(content *string, embedded bool) (err error) {
	if len(*content) == 0 {
		return
	}
	err = f.Edit(content, embedded)
	return
}

func (f *FileUpdater) Edit(content *string, embedded bool) (err error) {
	if embedded {
		err = f.rel.RelationEmbeddedFiles(f.project, f.table, f.field, f.tableID, *content)
	} else {
		err = f.rel.RelationFiles(f.project, f.table, f.field, f.tableID, *content, f.seperator)
	}
	if err != nil {
		return
	}
	return f.replace(content, embedded)
}

func (f *FileUpdater) replace(content *string, embedded bool) (err error) {
	if len(f.tableID) == 0 || f.tableID == `0` {
		log.Error(`FileUpdater: tableID is empty`)
		return
	}
	var replaces map[string]string
	fileIDs := f.rel.FileIDs()
	if len(fileIDs) > 0 {
		replaces, err = f.rel.MoveFileToOwner(f.table, fileIDs, f.tableID)
		if err != nil {
			return
		}
	}
	if mp, ok := f.rel.Context().Internal().Get(`FileReplaces`).(map[string]string); ok {
		if replaces == nil {
			replaces = mp
		} else {
			for k, v := range mp {
				replaces[k] = v
			}
		}
	}
	if replaces != nil {
		f.rel.Context().Internal().Set(`FileReplaces`, replaces)
	}
	if f.rel.ReplacedViewURLs() != nil {
		if replaces == nil {
			replaces = f.rel.ReplacedViewURLs()
		} else {
			for k, v := range f.rel.ReplacedViewURLs() {
				replaces[k] = v
			}
		}
	}
	if Debug {
		echo.Dump(echo.H{
			`project`: f.project,
			`table`:   f.table,
			`field`:   f.field,
			`tableID`: f.tableID,
			`replace`: replaces,
		})
	}
	if embedded {
		*content = uploadHelper.ReplaceEmbeddedRes(*content, replaces)
	} else {
		*content = uploadHelper.ReplaceRelatedRes(*content, replaces, f.seperator)
	}
	return
}

func (f *FileUpdater) Delete() (err error) {
	err = f.rel.DeleteByTableID(f.project, f.table, f.tableID)
	return
}
