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

// NewLister 创建页码分页列表查询
func NewLister(list Lister, recv interface{}, mw func(db.Result) db.Result, args ...interface{}) *List {
	return &List{
		ListParam: NewListParam(recv, mw, args...),
		ls:        list,
	}
}

// List 页码分页列表封装
type List struct {
	*ListParam
	ls Lister
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
func (f *List) DataTable(ctx echo.Context, args ...string) (map[string]interface{}, error) {
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
	data := map[string]interface{}{
		"draw":            ctx.Form(`draw`),
		"recordsTotal":    totalRows,
		"recordsFiltered": totalRows,
		"list":            f.recv,
	}
	return data, err
}
