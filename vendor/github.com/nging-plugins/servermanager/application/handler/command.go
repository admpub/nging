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

	"github.com/nging-plugins/servermanager/application/model"
	sshmodel "github.com/nging-plugins/sshmanager/application/model"
)

func Command(ctx echo.Context) error {
	m := model.NewCommand(ctx)
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}))
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`server/command`, handler.Err(ctx, err))
}

func CommandAdd(ctx echo.Context) error {
	var err error
	m := model.NewCommand(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingCommand)
		if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/server/command`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.Command, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	ctx.Set(`activeURL`, `/server/command`)
	sshUser := sshmodel.NewSshUser(ctx)
	_, err = sshUser.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, 0, -1)
	ctx.Set(`sshAccountList`, sshUser.Objects())
	return ctx.Render(`server/command_edit`, handler.Err(ctx, err))
}

func CommandEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewCommand(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/server/command`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingCommand)
		if err == nil {
			m.Id = id
			err = m.Edit(nil, `id`, id)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/server/command`))
		}
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			m.Disabled = disabled
			data := ctx.Data()
			err = m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
	}

	echo.StructToForm(ctx, m.NgingCommand, ``, echo.LowerCaseFirstLetter)
	ctx.Set(`activeURL`, `/server/command`)
	sshUser := sshmodel.NewSshUser(ctx)
	_, err = sshUser.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, 0, -1)
	ctx.Set(`sshAccountList`, sshUser.Objects())
	return ctx.Render(`server/command_edit`, handler.Err(ctx, err))
}

func CommandDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewCommand(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/server/command`))
}
