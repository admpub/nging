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
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v4/application/dbschema"
)

func NewUserRolePermission(ctx echo.Context) *UserRolePermission {
	return &UserRolePermission{
		NgingUserRolePermission: dbschema.NewNgingUserRolePermission(ctx),
	}
}

type UserRolePermission struct {
	*dbschema.NgingUserRolePermission
}

func (u *UserRolePermission) check() error {
	if len(u.Type) == 0 {
		return u.Context().NewError(code.InvalidParameter, `权限类型不能为空`).SetZone(`type`)
	}
	if u.RoleId < 1 {
		return u.Context().NewError(code.InvalidParameter, `权限的角色ID无效`).SetZone(`roleId`)
	}
	return nil
}

func (u *UserRolePermission) Add() (interface{}, error) {
	if err := u.check(); err != nil {
		return nil, err
	}
	exists, err := u.Exists(u.RoleId, u.Type)
	if err != nil {
		return nil, err
	}
	if exists {
		err = u.UpdateFields(nil, echo.H{`permission`: u.Permission}, db.And(
			db.Cond{`role_id`: u.RoleId},
			db.Cond{`type`: u.Type},
		))
		return nil, err
	}
	return u.NgingUserRolePermission.Insert()
}

func (u *UserRolePermission) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := u.check(); err != nil {
		return err
	}
	return u.NgingUserRolePermission.Update(mw, args...)
}

func (u *UserRolePermission) Exists(roleID uint, typ string) (bool, error) {
	return u.NgingUserRolePermission.Exists(nil, db.And(
		db.Cond{`role_id`: roleID},
		db.Cond{`type`: typ},
	))
}
