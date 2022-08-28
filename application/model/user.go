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

package model

import (
	"encoding/gob"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v4/application/dbschema"
)

func init() {
	gob.Register(dbschema.NewNgingUser(nil))
}

func NewUser(ctx echo.Context) *User {
	m := &User{
		NgingUser: dbschema.NewNgingUser(ctx),
	}
	return m
}

type User struct {
	*dbschema.NgingUser
}

func (u *User) Exists(username string) (bool, error) {
	return u.NgingUser.Exists(nil, db.Cond{`username`: username})
}

func (u *User) Exists2(username string, excludeUID uint) (bool, error) {
	return u.NgingUser.Exists(nil, db.And(
		db.Cond{`username`: username},
		db.Cond{`id`: db.NotEq(excludeUID)},
	))
}

func (u *User) check(editMode bool) (err error) {
	ctx := u.Context()
	if len(u.Username) == 0 {
		return ctx.NewError(code.InvalidParameter, `用户名不能为空`).SetZone(`username`)
	}
	if len(u.Email) == 0 {
		return ctx.NewError(code.InvalidParameter, `Email不能为空`).SetZone(`email`)
	}
	if !com.IsUsername(u.Username) {
		return ctx.NewError(code.InvalidParameter, `用户名不能包含特殊字符(只能由字母、数字、下划线和汉字组成)`).SetZone(`username`)
	}
	if !ctx.Validate(`email`, u.Email, `email`).Ok() {
		return ctx.NewError(code.InvalidParameter, `Email地址"%s"格式不正确`, u.Email).SetZone(`email`)
	}
	if len(u.Mobile) > 0 && !ctx.Validate(`mobile`, u.Mobile, `mobile`).Ok() {
		return ctx.NewError(code.InvalidParameter, `手机号"%s"格式不正确`, u.Mobile).SetZone(`mobile`)
	}
	if !editMode || ctx.Form(`modifyPwd`) == `1` {
		if len(u.Password) < 8 {
			return ctx.NewError(code.InvalidParameter, `密码不能少于8个字符`).SetZone(`password`)
		}
	}
	if len(u.Disabled) == 0 {
		u.Disabled = `N`
	}
	if len(u.Online) == 0 {
		u.Online = `N`
	}
	var exists bool
	if editMode {
		exists, err = u.Exists2(u.Username, u.Id)
	} else {
		exists, err = u.Exists(u.Username)
	}
	if err != nil {
		return
	}
	if exists {
		err = ctx.NewError(code.DataAlreadyExists, `用户名已经存在`).SetZone(`username`)
	}
	return
}

func (u *User) Add() (err error) {
	err = u.check(false)
	if err != nil {
		return
	}
	u.Salt = com.Salt()
	u.Password = com.MakePassword(u.Password, u.Salt)
	_, err = u.NgingUser.Insert()
	return
}

func (u *User) UpdateField(uid uint, set map[string]interface{}) (err error) {
	err = u.check(true)
	if err != nil {
		return
	}
	ctx := u.Context()
	if ctx.Form(`modifyPwd`) == `1` {
		u.Password = com.MakePassword(u.Password, u.Salt)
		set[`password`] = u.Password
	}
	err = u.NgingUser.UpdateFields(nil, set, `id`, uid)
	return
}

func (u *User) ClearPasswordData(users ...*dbschema.NgingUser) dbschema.NgingUser {
	var user dbschema.NgingUser
	if len(users) > 0 {
		user = *(users[0])
	} else {
		user = *(u.NgingUser)
	}
	user.Password = ``
	user.Salt = ``
	user.SafePwd = ``
	user.SessionId = ``
	return user
}

func (u *User) NewLoginLog(username string) *LoginLog {
	loginLogM := NewLoginLog(u.Context())
	loginLogM.OwnerType = `user`
	loginLogM.Username = username
	loginLogM.Success = `N`
	loginLogM.SessionId = u.Context().Session().MustID()
	return loginLogM
}
