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

package file

import (
	"fmt"
	"strings"

	"github.com/admpub/db"
	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/echo"
)

//UpdateUnrelation 更新未关联的附件
func (f *File) UpdateUnrelation(project string, table string, field string, tableid uint64, fileIds ...interface{}) (err error) {
	if len(fileIds) < 1 {
		return nil
	}
	err = f.SetFields(nil, echo.H{
		`project`:    project,
		`table_name`: table,
		`field_name`: field,
		`table_id`:   tableID,
	}, db.And(
		db.Cond{`table_id`: 0},
		db.Cond{`id`: db.In(fileIds)},
	))
	return
}

//Incr 增加使用次数
func (f *File) Incr(fileIds ...interface{}) (err error) {
	err = f.SetField(nil, `used_times`, db.Raw(`used_times+1`), db.Cond{`id`: db.In(fileIds)})
	return
}

//Decr 减少使用次数
func (f *File) Decr(fileIds ...interface{}) (err error) {
	err = f.SetField(nil, `used_times`, db.Raw(`used_times-1`), db.Cond{`id`: db.In(fileIds)})
	return
}

// DeleteEmbeddedByTableID 删除嵌入文件
func (f *File) DeleteEmbeddedByTableID(project string, table string, tableid uint64) error {
	_, err := f.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`project`: project},
		db.Cond{`table_id`: tableid},
		db.Cond{`table_name`: table},
	))
	if err != nil {
		return err
	}
	var ids []uint64
	for _, row := range f.Objects() {
		err = f.SetField(nil, `used_times`, db.Raw(`used_times-1`), db.And(
			db.Cond{`used_times`: db.Gt(0)},
			db.Cond{`id`: db.In(strings.Split(row.FileIds, `,`))},
		))
		if err != nil {
			return err
		}
		ids = append(ids, row.Id)
	}
	if len(ids) > 0 {
		embM := NewEmbedded(f.base.Context)
		err = embM.Delete(nil, db.Cond{`id`: db.In(ids)})
	}
	return err
}

func (f *File) UpdateEmbeddedByFileID(project string, table string, field string, tableID uint64, fileID uint64) error {
	_, err := f.UpdateUnrelation(project, table, field, tableID, fileID)
	if err != nil {
		return err
	}
	err = f.Incr(fileID)
	if err != nil {
		return err
	}
	m := &dbschema.FileEmbedded{}
	err = m.Get(nil, db.And(
		db.Cond{`table_id`: tableID},
		db.Cond{`field_name`: field},
		db.Cond{`table_name`: table},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		m.Reset()
		m.FieldName = field
		m.TableName = table
		m.Project = project
		m.TableId = tableID
		m.FileIds = fmt.Sprint(fileID)
		err = m.Add()
	}
	return err
}

func (f *File) UpdateEmbedded(project string, table string, field string, tableid uint64, fileIds ...interface{}) error {

	err := f.UpdateUnrelation(project, table, field, tableid, fileIds...)
	if err != nil {
		return err
	}

	m := &dbschema.FileEmbedded{}
	err = m.Get(nil, db.And(
		db.Cond{`table_id`: tableID},
		db.Cond{`field_name`: field},
		db.Cond{`table_name`: table},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		if len(fileIds) < 1 {
			return nil
		}
		m.Reset()
		m.FieldName = field
		m.TableName = table
		m.Project = project
		m.TableId = tableid
		m.FileIds = ""
		err = f.Incr(fileIds...)
		if err != nil {
			return err
		}
		for _, v := range fileIds {
			m.FileIds += fmt.Sprintf("%v,", v)
		}
		m.FileIds = strings.TrimSuffix(m.FileIds, ",")
		err = m.Add()
		return err
	}
	if len(fileIds) < 1 {
		err = f.Delete(nil, `id`, m.Id)
		if err != nil {
			return err
		}
	}
	var fidsString string
	for _, v := range fileIds {
		fidsString += fmt.Sprintf("%v,", v)
	}
	fidsString = strings.TrimSuffix(fidsString, ",")
	if m.FileIds == fidsString {
		return nil
	}
	ids := strings.Split(m.FileIds, ",")
	var (
		delIds []interface{}
		newIds []interface{}
	)
	//已删除引用
	for _, v := range ids {
		var has bool
		for _, v2 := range fileIds {
			if fmt.Sprint(v2) == v {
				has = true
			}
		}
		if has == false {
			delIds = append(delIds, v)
		}
	}
	//新增引用
	for _, v2 := range fileIds {
		var has bool
		for _, v := range ids {
			if fmt.Sprint(v2) == v {
				has = true
			}
		}
		if has == false {
			newIds = append(newIds, v2)
		}
	}
	if len(delIds) > 0 {
		_, err := f.Decr(delIds...)
		if err != nil {
			return err
		}
		err = f.SetFields(nil, echo.H{
			`table_id`:   0,
			`table_name`: ``,
			`field_name`: ``,
		}, db.Cond{`used_times`: 0})
		if err != nil {
			return err
		}
	}
	if len(newIds) > 0 {
		_, err := f.Incr(newIds...)
		if err != nil {
			return err
		}
	}
	m.FileIds = fidsString
	err = f.SetField(nil, `file_ids`, m.FileIds, db.Cond{`id`: m.Id})
	return err
}

// RelationEmbeddedFiles 关联嵌入的文件
// @param tableID 表中行主键ID
// @param field 被嵌入的字段名
// @param table 表名称
// @param v 内容
// @return
// @author AdamShen <swh@admpub.com>
func (f *File) RelationEmbeddedFiles(project string, table string, field string, tableID uint64, v string) error {
	var (
		files []interface{}
		fids  []interface{}
	)
	EmbeddedRes(v, func(file string, fid int64) {
		var exists bool
		if fid > 0 {
			for _, id := range fids {
				if fid == id {
					exists = true
					break
				}
			}
		} else {
			for _, rfile := range files {
				if rfile == file {
					exists = true
					break
				}
			}
		}
		if exists == false {
			files = append(files, file)
			if fid > 0 {
				fids = append(fids, fid)
			}
		}
	})
	if len(fids) < 1 && len(files) > 0 {
		fids = f.GetIDByViewURLs(files)
	}
	err := f.UpdateEmbedded(project, table, field, tableID, fids...)
	return err
}

func (f *File) GetViewURLByIds(ids ...interface{}) (r map[string]interface{}) {
	_, err := f.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `view_url`)
	}, 0, -1, db.Cond{
		`id`: db.In(ids),
	})
	if err != nil {
		return
	}
	r = f.AsKV(`id`, `view_url`)
	return
}

func (f *File) GetIDByViewURLs(viewURLs []interface{}) (r []interface{}) {
	err := f.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `view_url`)
	}, 0, -1, db.Cond{
		`view_url`: db.In(viewURLs),
	})
	if err != nil {
		return []interface{}{}
	}
	for _, v := range f.Objects() {
		r = append(r, v.Id)
	}
	return
}
