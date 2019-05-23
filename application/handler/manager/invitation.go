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
	"fmt"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func dateString2UnixString(k string, v []string) (string, []string) {
	switch strings.Title(k) {
	case `Start`:
		fallthrough
	case `End`:
		if len(v) > 0 && len(v[0]) > 0 {
			t, e := time.Parse(`2006-01-02`, v[0])
			if e != nil {
				log.Error(e)
				return k, []string{`0`}
			}
			return k, []string{fmt.Sprint(t.Unix())}
		}
		return k, []string{`0`}
	case `RoleIds`:
		return k, []string{strings.Join(v, `,`)}
	}
	return k, v
}

func Invitation(ctx echo.Context) error {
	m := model.NewInvitation(ctx)
	cond := db.Compounds{}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`code`, q)
	}
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ret := handler.Err(ctx, err)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`/manager/invitation`, ret)
}

func InvitationAdd(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := model.NewInvitation(ctx)
		err = ctx.MustBind(m.CodeInvitation, dateString2UnixString)
		if err == nil {
			if len(m.Code) == 0 {
				err = ctx.E(`邀请码不能为空`)
			} else if exists, erro := m.Exists(m.Code); erro != nil {
				err = erro
			} else if exists {
				err = ctx.E(`邀请码已经存在`)
			} else {
				_, err = m.Add()
			}
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/manager/invitation`))
		}
	} else {
		ctx.Request().Form().Set(`code`, com.RandomAlphanumeric(16))
	}
	ctx.Set(`activeURL`, `/manager/invitation`)
	userRoleMdl := model.NewUserRole(ctx)
	userRoleMdl.List(nil, nil, 1, -1, `disabled`, `N`)
	ctx.Set(`roleList`, userRoleMdl.Objects())
	ctx.SetFunc(`isChecked`, func(roleId uint) bool {
		return false
	})
	return ctx.Render(`/manager/invitation_edit`, handler.Err(ctx, err))
}

func InvitationEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewInvitation(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/manager/invitation`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.CodeInvitation, dateString2UnixString)
		if err == nil {
			m.Id = id
			if len(m.Code) == 0 {
				err = ctx.E(`邀请码不能为空`)
			} else if exists, erro := m.Exists2(m.Code, id); erro != nil {
				err = erro
			} else if exists {
				err = ctx.E(`邀请码已经存在`)
			} else {
				err = m.Edit(nil, `id`, id)
			}
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/manager/invitation`))
		}
	} else {
		echo.StructToForm(ctx, m.CodeInvitation, ``, echo.LowerCaseFirstLetter)
		var startDate, endDate string
		if m.CodeInvitation.Start > 0 {
			startDate = time.Unix(int64(m.CodeInvitation.Start), 0).Format(`2006-01-02`)
		}
		ctx.Request().Form().Set(`start`, startDate)
		if m.CodeInvitation.End > 0 {
			endDate = time.Unix(int64(m.CodeInvitation.End), 0).Format(`2006-01-02`)
		}
		ctx.Request().Form().Set(`end`, endDate)
	}

	ctx.Set(`activeURL`, `/manager/invitation`)
	userRoleMdl := model.NewUserRole(ctx)
	userRoleMdl.List(nil, nil, 1, -1, `disabled`, `N`)
	ctx.Set(`roleList`, userRoleMdl.Objects())
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
	return ctx.Render(`/manager/invitation_edit`, handler.Err(ctx, err))
}

func InvitationDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewInvitation(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/manager/invitation`))
}
