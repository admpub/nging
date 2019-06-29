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
