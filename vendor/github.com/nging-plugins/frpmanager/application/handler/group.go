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
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"

	"github.com/nging-plugins/frpmanager/application/model"
)

func GroupIndex(ctx echo.Context) error {
	m := model.NewFrpGroup(ctx)
	_, err := handler.PagingWithLister(ctx, m)
	ret := handler.Err(ctx, err)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`frp/group_index`, ret)
}

func GroupAdd(ctx echo.Context) error {
	var err error
	m := model.NewFrpGroup(ctx)
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`用户组名称不能为空`)
		} else if y, e := m.Exists(name); e != nil {
			err = e
		} else if y {
			err = ctx.E(`用户组名称已经存在`)
		} else {
			err = ctx.MustBind(m.NgingFrpGroup)
		}
		if err == nil {
			_, err = m.Insert()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/frp/group_index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.NgingFrpGroup, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/frp/group_index`)
	return ctx.Render(`frp/group_edit`, err)
}

func GroupEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewFrpGroup(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`用户组名称不能为空`)
		} else if y, e := m.ExistsOther(name, id); e != nil {
			err = e
		} else if y {
			err = ctx.E(`用户组名称已经存在`)
		} else {
			err = ctx.MustBind(m.NgingFrpGroup, echo.ExcludeFieldName(`created`))
		}

		if err == nil {
			m.Id = id
			err = m.Update(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/frp/group_index`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.NgingFrpGroup, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/frp/group_index`)
	return ctx.Render(`frp/group_edit`, err)
}

func GroupDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewFrpGroup(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/frp/group_index`))
}
