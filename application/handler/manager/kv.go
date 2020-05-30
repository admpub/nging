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

package manager

import (
	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/model"
)

func KvIndex(ctx echo.Context) error {
	m := model.NewKv(ctx)
	cond := db.Compounds{}
	t := ctx.Formx(`type`).String()
	if len(t) > 0 {
		cond.AddKV(`type`, t)
	}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`key`, q)
		cond.From(mysql.MatchAnyField(`value`, q))
	}
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ctx.Set(`listData`, m.Objects())
	ctx.Set(`title`, ctx.E(`元数据`))
	typeList := m.KvTypeList()
	ctx.Set(`typeList`, typeList)
	ctx.SetFunc(`typeName`, func(key string) string {
		for _, row := range typeList {
			if row.Key == key {
				return row.Value
			}
		}
		return ``
	})
	return ctx.Render(`/manager/kv`, handler.Err(ctx, err))
}

func KvAdd(ctx echo.Context) error {
	var err error
	m := model.NewKv(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingKv)
		if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/manager/kv`))
		}
	}
	ctx.Set(`activeURL`, `/manager/kv`)
	ctx.Set(`title`, ctx.E(`添加元数据`))
	ctx.Set(`typeList`, m.KvTypeList())
	return ctx.Render(`/manager/kv_edit`, handler.Err(ctx, err))
}

func KvEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewKv(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/manager/tv`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingKv)
		if err == nil {
			m.Id = id
			err = m.Edit(nil, `id`, id)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/manager/kv`))
		}
	} else {
		echo.StructToForm(ctx, m.NgingKv, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/manager/kv`)
	ctx.Set(`title`, ctx.E(`修改元数据`))
	var typeList []*dbschema.NgingKv
	if m.IsRootType(m.Type) {
		typeList = m.KvTypeList(m.Id)
	} else {
		typeList = m.KvTypeList()
	}
	ctx.Set(`typeList`, typeList)
	return ctx.Render(`/manager/kv_edit`, handler.Err(ctx, err))
}

func KvDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewKv(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/manager/kv`))
}
