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

package handler

import (
	"github.com/admpub/nging/v4/application/handler"
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func Group(ctx echo.Context) error {
	m := dbschema.NewNgingVhostGroup(ctx)
	user := handler.User(ctx)
	_, err := handler.PagingWithListerCond(ctx, m, db.Cond{`uid`: user.Id})
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`caddy/group`, handler.Err(ctx, err))
}

func GroupAdd(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := dbschema.NewNgingVhostGroup(ctx)
		err = ctx.MustBind(m)
		if err == nil {
			if len(m.Name) == 0 {
				err = ctx.E(`分组名称不能为空`)
			} else {
				user := handler.User(ctx)
				m.Uid = user.Id
				_, err = m.Insert()
			}
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/caddy/group`))
		}
	}
	ctx.Set(`activeURL`, `/caddy/group`)
	return ctx.Render(`caddy/group_edit`, handler.Err(ctx, err))
}

func GroupEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := dbschema.NewNgingVhostGroup(ctx)
	user := handler.User(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/caddy/group`))
	}
	if m.Id < 1 || user.Id != m.Uid {
		handler.SendFail(ctx, ctx.T(`分组不存在`))
		return ctx.Redirect(handler.URLFor(`/caddy/group`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m)
		if err == nil {
			m.Id = id
			if len(m.Name) == 0 {
				err = ctx.E(`分组名称不能为空`)
			} else {
				m.Uid = user.Id
				err = m.Update(nil, `id`, id)
			}
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/caddy/group`))
		}
	}
	echo.StructToForm(ctx, m, ``, echo.LowerCaseFirstLetter)
	ctx.Set(`activeURL`, `/caddy/group`)
	return ctx.Render(`caddy/group_edit`, handler.Err(ctx, err))
}

func GroupDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := dbschema.NewNgingVhostGroup(ctx)
	user := handler.User(ctx)
	err := m.Delete(nil, db.And(
		db.Cond{`id`: id},
		db.Cond{`uid`: user.Id},
	))
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}
	return ctx.Redirect(handler.URLFor(`/caddy/group`))
}
