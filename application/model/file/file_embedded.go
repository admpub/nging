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
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

//UpdateUnrelation 更新未关联的附件
func (f *File) UpdateUnrelation(project string, table string, field string, tableID uint64, fileIds ...interface{}) (err error) {
	if len(fileIds) < 1 {
		return nil
	}
	err = f.SetFields(nil, echo.H{
		`table_id`:   tableID,
		`table_name`: table,
		`field_name`: field,
		`project`:    project,
	}, db.And(
		db.Cond{`table_id`: 0},
		db.Cond{`id`: db.In(fileIds)},
	))
	return
}

// Incr 增加使用次数
func (f *File) Incr(fileIds ...interface{}) (err error) {
	err = f.SetField(nil, `used_times`, db.Raw(`used_times+1`), db.Cond{
		`id`: db.In(fileIds),
	})
	return
}

// Decr 减少使用次数
func (f *File) Decr(fileIds ...interface{}) (err error) {
	err = f.SetField(nil, `used_times`, db.Raw(`used_times-1`), db.And(
		db.Cond{`id`: db.In(fileIds)},
		db.Cond{`used_times`: db.NotEq(0)},
	))
	return
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
	_, err := f.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `view_url`)
	}, 0, -1, db.Cond{
		`view_url`: db.In(viewURLs),
	})
	if err != nil {
		return r
	}
	for _, v := range f.Objects() {
		r = append(r, v.Id)
	}
	return
}
