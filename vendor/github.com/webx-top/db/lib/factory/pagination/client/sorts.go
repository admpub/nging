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

package client

import (
	"strings"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

// Sorts 获取数据查询时的排序方式
func Sorts(ctx echo.Context, table interface{}, defaultSorts ...string) []interface{} {
	sorts := make([]interface{}, 0, len(defaultSorts)+1)
	sort := ctx.Form(`sort`)
	field := strings.TrimPrefix(sort, `-`)
	noPrefixTableName := factory.NoPrefixTableName(table)
	if len(field) > 0 && factory.ExistField(noPrefixTableName, field) {
		sorts = append(sorts, sort)
		for _, defaultSort := range defaultSorts {
			if field != strings.TrimPrefix(defaultSort, `-`) {
				sorts = append(sorts, defaultSort)
			}
		}
	} else {
		for _, defaultSort := range defaultSorts {
			sorts = append(sorts, defaultSort)
		}
	}
	return sorts
}
