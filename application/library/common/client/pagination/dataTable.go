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

package pagination

import (
	"github.com/webx-top/echo"
)

// DataTable 分页信息
func DataTable(ctx echo.Context, ls interface{}, args ...string) (map[string]interface{}, error) {
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
	var (
		cnt  func() int64
		err  error
		recv interface{}
	)
	switch f := ls.(type) {
	case Lister:
		page := (offset + size) / size
		cnt, err = f.List(nil, nil, page, size)
		recv = f.Model()
	case OffsetLister:
		cnt, err = f.ListByOffset(nil, nil, offset, size)
		recv = f.Model()
	default:
		panic(ctx.T(`不支持的分页类型: %T`, f))
	}
	totalRows := ctx.Formx(totalRowsKey).Int()
	if totalRows < 1 {
		totalRows = int(cnt())
	}
	data := map[string]interface{}{
		"draw":            ctx.Form(`draw`),
		"recordsTotal":    totalRows,
		"recordsFiltered": totalRows,
		"list":            recv,
	}
	return data, err
}
