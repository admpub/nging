/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package cloud

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/model"
)

func StorageIndex(ctx echo.Context) error {
	m := model.NewCloudStorage(ctx)
	cond := db.Compounds{}
	common.SelectPageCond(ctx, &cond)
	_, err := common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()).Paging(ctx)
	list := m.Objects()
	ctx.Set(`listData`, list)
	return ctx.Render(`cloud/storage`, handler.Err(ctx, err))
}

func StorageAdd(ctx echo.Context) error {
	var (
		err error
		id  uint
	)
	m := model.NewCloudStorage(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingCloudStorage)
		if err != nil {
			goto END
		}
		m.Baseurl = ctx.Formx(`baseurl`).String()
		_, err = m.Add()
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(handler.URLFor(`/cloud/storage`))
	}
	id = ctx.Formx(`copyId`).Uint()
	if id > 0 {
		err = m.Get(nil, `id`, id)
		if err == nil {
			echo.StructToForm(ctx, m.NgingCloudStorage, ``, func(topName, fieldName string) string {
				return echo.LowerCaseFirstLetter(topName, fieldName)
			})
			ctx.Request().Form().Set(`id`, `0`)
		}
	}

END:
	ctx.Set(`isAdd`, true)
	ctx.Set(`title`, ctx.T(`添加云存储账号`))
	ctx.Set(`activeURL`, `/cloud/storage`)
	return ctx.Render(`cloud/storage_edit`, err)
}

func StorageEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudStorage(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingCloudStorage, echo.ExcludeFieldName(`created`))
		if err != nil {
			goto END
		}
		m.Id = id
		m.Baseurl = ctx.Formx(`baseurl`).String()
		err = m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(handler.URLFor(`/cloud/storage`))
	}
	if err == nil {
		echo.StructToForm(ctx, m.NgingCloudStorage, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

END:
	ctx.Set(`isAdd`, false)
	ctx.Set(`title`, ctx.T(`修改云存储账号`))
	ctx.Set(`activeURL`, `/cloud/storage`)
	return ctx.Render(`cloud/storage_edit`, err)
}

func StorageDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudStorage(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/cloud/storage`))
}
