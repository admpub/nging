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

package file

import (
	"strings"

	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/model/file"
)

func List(ctx echo.Context, ownerType string, ownerID uint64) error {
	fileM := file.NewFile(ctx)
	cond := db.NewCompounds()
	if len(ownerType) > 0 {
		cond.AddKV(`owner_id`, ownerID)
		cond.AddKV(`owner_type`, ownerType)
	}
	table := ctx.Formx("table").String()
	if len(table) > 0 {
		tableAndField := strings.SplitN(table, ".", 2)
		cond.AddKV(`table_name`, tableAndField[0])
		if len(tableAndField) > 1 {
			cond.AddKV(`field_name`, tableAndField[1])
		}
	}
	used := ctx.Formx("used").String()
	if len(used) > 0 {
		switch used {
		case `0`:
			cond.AddKV(`used_times`, 0)
		case `1`:
			cond.AddKV(`used_times`, db.Gt(0))
		}
	}
	typ := ctx.Formx(`type`).String()
	if len(typ) == 0 {
		typ = ctx.Formx(`filetype`).String()
	}
	if len(typ) > 0 {
		cond.AddKV(`type`, typ)
	}
	timerange := ctx.Formx(`timerange`).String()
	if len(timerange) > 0 {
		cond.Add(mysql.GenDateRange(`created`, timerange).V()...)
	}
	if len(ctx.Formx(`searchValue`).String()) > 0 {
		cond.AddKV(`id`, ctx.Formx(`searchValue`).Uint64())
	}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.Add(
			db.Or(
				db.Cond{`save_name`: db.Like(q + `%`)},
				db.Cond{`name`: db.Like(`%` + q + `%`)},
			),
		)
	}
	sorts := common.Sorts(ctx, `file`, `-id`)
	_, err := common.NewLister(fileM.NgingFile, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(ctx)
	if err != nil {
		return err
	}
	list := fileM.Objects()
	ctx.Set(`listData`, list)
	ctx.Set(`fileTypes`, uploadClient.FileTypeExts)
	ctx.Set(`fileType`, typ)
	return err
}
