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
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

// OffsetLister 偏移值分页列表查询接口
type OffsetLister interface {
	ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error)
}

// NewOffsetLister 创建偏移值分页列表查询
func NewOffsetLister(list OffsetLister, recv interface{}, mw func(db.Result) db.Result, args ...interface{}) *OffsetList {
	return &OffsetList{
		ListParam: NewListParam(recv, mw, args...),
		ls:        list,
	}
}

// OffsetList 偏移值分页列表查询封装
type OffsetList struct {
	*ListParam
	ls OffsetLister
}

// ListByOffset 分页查询
func (f *OffsetList) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = f.recv
	}
	if mw == nil {
		mw = f.mw
	}
	if len(args) < 1 {
		args = f.args
	}
	return f.ls.ListByOffset(recv, mw, offset, size, args...)
}

// ChunkList 分批查询列表
func (f *OffsetList) ChunkList(eachPageCallback func() error, size int, offset int) error {
	cnt, err := f.ListByOffset(nil, nil, offset, size)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return err
	}
	for total := cnt(); int64(offset) < total; offset += size {
		if offset > 0 {
			_, err = f.ListByOffset(nil, nil, offset, size)
			if err != nil {
				if err == db.ErrNoMoreRows {
					return nil
				}
				return err
			}
		}
		err = eachPageCallback()
		if err != nil {
			return err
		}
	}
	return err
}

// DataTable 分页信息
func (f *OffsetList) DataTable(ctx echo.Context, args ...string) (map[string]interface{}, error) {
	pageRowsKey := `length`
	totalRowsKey := `totalrows`
	offsetKey := `start`
	switch len(args) {
	case 3:
		pageRowsKey = args[2]
		fallthrough
	case 2:
		totalRowsKey = args[1]
		fallthrough
	case 1:
		offsetKey = args[0]
	}
	size := ctx.Formx(pageRowsKey).Int()
	offset := ctx.Formx(offsetKey).Int()
	if size < 1 || size > 1000 {
		size = 10
	}
	if offset < 0 {
		offset = 0
	}
	totalRows := ctx.Formx(totalRowsKey).Int()
	cnt, err := f.ListByOffset(nil, nil, offset, size)
	if totalRows < 1 {
		totalRows = int(cnt())
	}
	data := map[string]interface{}{
		"draw":            ctx.Form(`draw`),
		"recordsTotal":    totalRows,
		"recordsFiltered": totalRows,
		"list":            f.recv,
	}
	return data, err
}
