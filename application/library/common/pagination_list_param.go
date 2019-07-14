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

import "github.com/webx-top/db"

// NewListParam 列表参数
func NewListParam(recv interface{}, mw func(db.Result) db.Result, args ...interface{}) *ListParam {
	return &ListParam{
		mw:   mw,
		recv: recv,
		args: args,
	}
}

// ListParam 列表参数
type ListParam struct {
	recv interface{}
	mw   func(db.Result) db.Result
	args []interface{}
}

// AddMiddleware 添加中间件
func (f *ListParam) AddMiddleware(mw ...func(db.Result) db.Result) {
	if f.mw != nil {
		origin := f.mw
		f.mw = func(r db.Result) db.Result {
			r = origin(r)
			for _, m := range mw {
				r = m(r)
			}
			return r
		}
		return
	}
	f.mw = func(r db.Result) db.Result {
		for _, m := range mw {
			r = m(r)
		}
		return r
	}
}

// AddCond 添加条件
func (f *ListParam) AddCond(args ...interface{}) {
	f.args = append(f.args, args...)
}

// Model 模型实例
func (f *ListParam) Model() interface{} {
	return f.recv
}
