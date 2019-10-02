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

// Package file 上传文件管理
package file

import (
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/model/file"
	"github.com/admpub/nging/application/registry/upload"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"
)

func FileList(ctx echo.Context) error {
	fileM := file.NewFile(ctx)
	cond := db.NewCompounds()
	table := ctx.Formx("table").String()
	if len(table) > 0 {
		cond.AddKV(`table_name`, table)
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
	typ := ctx.Formx("type").String()
	if len(typ) > 0 {
		cond.AddKV(`type`, typ)
	}
	timerange := ctx.Formx("timerange").String()
	if len(timerange) > 0 {
		cond.Add(mysql.GenDateRange(`created`, timerange).V()...)
	}
	sorts := common.Sorts(ctx, `file`, `-id`)
	_, err := common.NewLister(fileM.File, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(ctx)
	if err != nil {
		return err
	}
	list := fileM.Objects()
	listData := file.FileList(list)
	ctx.Set(`listData`, listData)
	ctx.Set(`fileTypes`, uploadClient.FileTypeExts)
	ctx.Set(`tableNames`, upload.SubdirAll())
	return ctx.Render(`manager/file/list`, err)
}

func FileDelete(ctx echo.Context) (err error) {
	user := handler.User(ctx)
	id := ctx.Paramx("id").Uint64()
	fileM := file.NewFile(ctx)
	ownerID := uint64(user.Id)
	if id == 0 {
		ids := ctx.FormxValues(`id`).Uint64()
		for _, id := range ids {
			err = fileM.DeleteByID(id, `user`, ownerID)
			if err != nil {
				return err
			}
		}
		goto END
	}
	err = fileM.DeleteByID(id, `user`, ownerID)
	if err != nil {
		return err
	}

END:
	return ctx.Redirect(handler.URLFor(`/manager/file/list`))
}
