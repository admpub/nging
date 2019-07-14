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

type PageListFunc func(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error)

func (f PageListFunc) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
	return f(recv, mw, page, size, args...)
}

type OffsetListFunc func(recv interface{}, mw func(db.Result) db.Result, pageOrOffset, size int, args ...interface{}) (func() int64, error)

func (f OffsetListFunc) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
	return f(recv, mw, offset, size, args...)
}

type ListFunc = PageListFunc
