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
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"
)

// DataGrid 分页信息
func DataGrid(ctx echo.Context, ls interface{}, searchField string, args ...string) (map[string]interface{}, error) {
	pageRowsKey := `recPerPage`
	totalRowsKey := `totalRows`
	pageKey := `page`
	switch len(args) {
	case 3:
		pageRowsKey = args[2]
		fallthrough
	case 2:
		totalRowsKey = args[1]
		fallthrough
	case 1:
		pageKey = args[0]
	}
	size := ctx.Formx(pageRowsKey).Int()
	page := ctx.Formx(pageKey).Int()
	if page < 1 {
		page = 1
	}
	var (
		cnt     func() int64
		err     error
		recv    interface{}
		buildMW = func(fMW func(r db.Result) db.Result) func(r db.Result) db.Result {
			sortBy := ctx.Formx(`sortBy`).String()
			if len(sortBy) > 0 {
				if ctx.Form(`order`) != `asc` {
					sortBy = `-` + sortBy
				}
				typeMap := sqlbuilder.Mapper().StructMap(recv)
				ctx.Request().Form().Set(`sort`, sortBy)
				sorts := Sorts(ctx, factory.DBIGet().TableName(typeMap.Tree.Name))
				if len(sorts) > 0 {
					return func(r db.Result) db.Result {
						if fMW != nil {
							r = fMW(r)
						}
						return r.OrderBy(sorts...)
					}
				}
			}
			return fMW
		}

		buildCond = func(args []interface{}) []interface{} {
			if len(searchField) > 0 {
				search := ctx.Formx(`search`).String()
				if len(search) > 0 {
					newCond := func() db.Compound {
						return db.Cond{searchField: db.Like(`%` + com.AddSlashes(search, '_', '%') + `%`)}
					}
					switch len(args) {
					case 0:
						args = append(args, newCond())
					case 1:
						_, ok := args[0].(db.Compound)
						if ok {
							args = append(args, newCond())
						}
					}
				}
			}
			return args
		}
	)
	switch f := ls.(type) {
	case Lister:
		f.SetConds(buildCond(f.Conds()))
		recv = f.Model()
		cnt, err = f.List(nil, buildMW(f.Middleware()), page, size)
	case OffsetLister:
		f.SetConds(buildCond(f.Conds()))
		recv = f.Model()
		offset := com.Offset(uint(page), uint(size))
		cnt, err = f.ListByOffset(nil, buildMW(f.Middleware()), int(offset), size)
	default:
		panic(ctx.T(`不支持的分页类型: %T`, f))
	}
	totalRows := ctx.Formx(totalRowsKey).Int()
	if totalRows < 1 {
		totalRows = int(cnt())
	}
	data := map[string]interface{}{
		"result":  `success`, // 数据请求结果（必须），当值为 'succes'、'ok' 或 200 时则是为请求成功
		"data":    recv,      //原创数据对象数组
		"message": ``,        //用户界面提示消息，当请求结果失败时，可以使用此属性返回文本显示在用户界面上
		"pager": map[string]int{ //远程数据的分页信息对象（必须）
			"page":       page,      // 当前数据对应的页码
			"recTotal":   totalRows, // 总的数据数目
			"recPerPage": size,      // 每页数据数目
		},
	}
	return data, err
}
