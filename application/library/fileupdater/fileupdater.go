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

func (f *FileUpdater) Add(content string, embedded bool) (err error) {
	if len(content) == 0 {
		return
	}
	err = f.Edit(content, embedded)
	return
}

func (f *FileUpdater) Edit(content string, embedded bool) (err error) {
	if !embedded {
		err = f.rel.RelationFiles(f.project, f.table, f.field, f.tableID, content, f.seperator)
		return
	}
	err = f.rel.RelationEmbeddedFiles(f.project, f.table, f.field, f.tableID, content)
	return
}

func (f *FileUpdater) Delete() (err error) {
	err = f.rel.DeleteByTableID(f.project, f.table, f.tableID)
	return
}
