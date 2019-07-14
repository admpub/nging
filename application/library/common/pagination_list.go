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
	"github.com/webx-top/pagination"
)

// Lister 页码分页列表查询接口
type Lister interface {
	List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error)
}

// OffsetLister 偏移值分页列表查询接口
type OffsetLister interface {
	ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error)
}

// NewLister 创建页码分页列表查询
func NewLister(list Lister, recv interface{}, mw func(db.Result) db.Result, args ...interface{}) *List {
	return &List{
		mw:   mw,
		recv: recv,
		ls:   list,
		args: args,
	}
}

// List 页码分页列表封装
type List struct {
	recv interface{}
	mw   func(db.Result) db.Result
	ls   Lister
	args []interface{}
}

// List 分页查询
func (f *List) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = f.recv
	}
	if mw == nil {
		mw = f.mw
	}
	return f.ls.List(recv, mw, page, size, f.args...)
}

// Paging 分页信息
func (f *List) Paging(ctx echo.Context, varSuffix ...string) (*pagination.Pagination, error) {
	return PagingWithLister(ctx, f, varSuffix...)
}

// DataTable 分页信息
func (f *List) DataTable(ctx echo.Context, args ...string) (echo.H, error) {
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
	page := (offset + size) / size
	totalRows := ctx.Formx(totalRowsKey).Int()
	cnt, err := f.List(nil, nil, page, size)
	if totalRows < 1 {
		totalRows = int(cnt())
	}
	data := echo.H{
		"draw":            ctx.Form(`draw`),
		"recordsTotal":    totalRows,
		"recordsFiltered": totalRows,
		"list":            f.recv,
	}
	return data, err
}

// NewOffsetLister 创建偏移值分页列表查询
func NewOffsetLister(list OffsetLister, recv interface{}, mw func(db.Result) db.Result, args ...interface{}) *OffsetList {
	return &OffsetList{
		mw:   mw,
		recv: recv,
		ls:   list,
		args: args,
	}
}

// OffsetList 偏移值分页列表查询封装
type OffsetList struct {
	recv interface{}
	mw   func(db.Result) db.Result
	ls   OffsetLister
	args []interface{}
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
