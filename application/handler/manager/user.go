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

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func User(ctx echo.Context) error {
	username := ctx.Formx(`username`).String()
	online := ctx.Form(`online`)
	cond := db.Compounds{}
	if len(username) > 0 {
		cond.AddKV(`username`, db.Like(username+`%`))
	} else if len(ctx.Formx(`searchValue`).String()) > 0 {
		cond.AddKV(`id`, ctx.Formx(`searchValue`).Uint64())
	} else if len(online) > 0 {
		cond.AddKV(`online`, online)
	}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`username`, db.Like(`%`+q+`%`))
	}
	m := model.NewUser(ctx)
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.Select(factory.Fields.SortedFieldLists(`user`, `password`, `salt`, `safe_pwd`)...).OrderBy(`-id`)
	}, cond.And()))
	ret := handler.Err(ctx, err)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`/manager/user`, ret)
}

func UserAdd(ctx echo.Context) error {
	var err error
	m := model.NewUser(ctx)
	if ctx.IsPost() {
		m.Username = strings.TrimSpace(ctx.Form(`username`))
		m.Email = strings.TrimSpace(ctx.Form(`email`))
		m.Mobile = strings.TrimSpace(ctx.Form(`mobile`))
		m.Password = strings.TrimSpace(ctx.Form(`password`))
		m.Avatar = strings.TrimSpace(ctx.Form(`avatar`))
		m.Gender = strings.TrimSpace(ctx.Form(`gender`))
		m.RoleIds = strings.Join(ctx.FormValues(`roleIds`), `,`)
		err = m.Add()
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/manager/user`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				setFormData(ctx, m)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	ctx.Set(`activeURL`, `/manager/user`)
	roleM := model.NewUserRole(ctx)
	roleM.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `name`, `description`)
	}, 0, -1, db.And(db.Cond{`parent_id`: 0}))
	ctx.Set(`roleList`, roleM.Objects())
	ctx.SetFunc(`isChecked`, func(roleId uint) bool {
		return false
	})
	return ctx.Render(`/manager/user_edit`, handler.Err(ctx, err))
}

func UserEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewUser(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/manager/user`))
	}
	if ctx.IsPost() {
		m.Username = strings.TrimSpace(ctx.Form(`username`))
		m.Email = strings.TrimSpace(ctx.Form(`email`))
		m.Mobile = strings.TrimSpace(ctx.Form(`mobile`))
		m.Avatar = strings.TrimSpace(ctx.Form(`avatar`))
		m.Gender = strings.TrimSpace(ctx.Form(`gender`))
		m.RoleIds = strings.Join(ctx.FormValues(`roleIds`), `,`)
		if err == nil {
			m.Id = id
			set := map[string]interface{}{
				`email`:    m.Email,
				`mobile`:   m.Mobile,
				`username`: m.Username,
				`role_ids`: m.RoleIds,
				`avatar`:   m.Avatar,
				`gender`:   m.Gender,
			}
			err = m.UpdateField(id, set)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/manager/user`))
		}
	}

	setFormData(ctx, m)
	ctx.Set(`activeURL`, `/manager/user`)
	roleM := model.NewUserRole(ctx)
	roleM.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `name`, `description`)
	}, 0, -1, db.And(db.Cond{`parent_id`: 0}))
	ctx.Set(`roleList`, roleM.Objects())
	return ctx.Render(`/manager/user_edit`, handler.Err(ctx, err))
}

func setFormData(ctx echo.Context, m *model.User) {
	m.Password = ``
	m.Salt = ``
	echo.StructToForm(ctx, m.User, ``, echo.LowerCaseFirstLetter)
	var roleIds []uint
	if len(m.RoleIds) > 0 {
		roleIds = param.StringSlice(strings.Split(m.RoleIds, `,`)).Uint()
	}
	ctx.SetFunc(`isChecked`, func(roleId uint) bool {
		for _, rid := range roleIds {
			if rid == roleId {
				return true
			}
		}
		return false
	})
}

func UserDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := model.NewUser(ctx)
	if id == 1 {
		handler.SendFail(ctx, ctx.T(`创始人不可删除`))
		return ctx.Redirect(handler.URLFor(`/manager/user`))
	}
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		err = common.RemoveUploadedFile(`user-avatar`, id)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/manager/user`))
}
