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
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/model"
	"github.com/admpub/nging/v3/application/registry/navigate"
	"github.com/admpub/nging/v3/application/registry/perm"
)

func Role(ctx echo.Context) error {
	m := model.NewUserRole(ctx)
	_, err := handler.PagingWithLister(ctx, m)
	ret := handler.Err(ctx, err)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`/manager/role`, ret)
}

func RoleAdd(ctx echo.Context) error {
	var err error
	m := model.NewUserRole(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingUserRole)
		if err == nil {
			m.BuildPermAction(ctx.FormValues(`permAction[]`))
			m.PermCmd = strings.Join(ctx.FormValues(`permCmd[]`), `,`)
			err = m.BuildPermBehavior(ctx.FormValues(`permBehavior[]`))
		}
		if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/manager/role`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.NgingUserRole, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	ctx.Set(`activeURL`, `/manager/role`)
	ctx.Set(`topNavigate`, navigate.TopNavigate)
	cmdM := model.NewCommand(ctx)
	cmdM.ListByOffset(nil, nil, 0, -1, `disabled`, `N`)
	ctx.Set(`cmdList`, cmdM.Objects())
	ctx.Set(`behaviorList`, perm.Behaviors.Slice())
	ctx.Set(`data`, m)
	return ctx.Render(`/manager/role_edit`, handler.Err(ctx, err))
}

func RoleEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewUserRole(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/manager/role`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingUserRole)
		if err == nil {
			m.Id = id
			m.BuildPermAction(ctx.FormValues(`permAction[]`))
			m.PermCmd = strings.Join(ctx.FormValues(`permCmd[]`), `,`)
			err = m.BuildPermBehavior(ctx.FormValues(`permBehavior[]`))
		}
		if err == nil {
			err = m.Edit(nil, `id`, id)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/manager/role`))
		}
	}

	echo.StructToForm(ctx, m.NgingUserRole, ``, echo.LowerCaseFirstLetter)
	ctx.Set(`activeURL`, `/manager/role`)
	ctx.Set(`topNavigate`, navigate.TopNavigate)
	cmdM := model.NewCommand(ctx)
	cmdM.ListByOffset(nil, nil, 0, -1, `disabled`, `N`)
	ctx.Set(`cmdList`, cmdM.Objects())
	ctx.Set(`behaviorList`, perm.Behaviors.Slice())
	ctx.Set(`data`, m)
	return ctx.Render(`/manager/role_edit`, handler.Err(ctx, err))
}

func RoleDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewUserRole(ctx)
	if id == 1 {
		handler.SendFail(ctx, ctx.T(`超级管理员角色不可删除`))
		return ctx.Redirect(handler.URLFor(`/manager/role`))
	}
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/manager/role`))
}
