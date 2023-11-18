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

package utils

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/webx-top/echo"
)

func AjaxListSelectpage[T any](ctx echo.Context, list []T, callback func(v T) echo.H) error {
	data := ctx.Data()
	rows := make([]echo.H, 0, len(list))
	var sk, n int
	_, size, _, pg := handler.PagingWithPagination(ctx)
	pg.SetRows(len(list))
	offset := pg.Offset()
	for _, v := range list {
		if n >= size {
			break
		}
		sk++
		if sk-1 < offset {
			continue
		}
		row := callback(v)
		if row == nil {
			continue
		}
		rows = append(rows, row)
		n++
	}
	data.SetData(echo.H{`listData`: rows, `pagination`: pg})
	return ctx.JSON(data)
}

func AjaxListTypeahead[T any](ctx echo.Context, list []T, callback func(v T) string) error {
	data := ctx.Data()
	names := make([]string, 0, len(list))
	for _, v := range list {
		name := callback(v)
		if len(name) == 0 {
			continue
		}
		names = append(names, name)
	}
	data.SetData(echo.H{`listData`: names})
	return ctx.JSON(data)
}
