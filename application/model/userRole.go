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
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/library/perm"
)

func NewUserRole(ctx echo.Context) *UserRole {
	return &UserRole{
		NgingUserRole: dbschema.NewNgingUserRole(ctx),
	}
}

type UserRole struct {
	*dbschema.NgingUserRole
	permActions   *perm.Map
	permCmds      *perm.Map
	permBehaviors perm.BehaviorPerms
}

func (u *UserRole) check() error {
	if len(u.Name) == 0 {
		return u.Context().NewError(code.InvalidParameter, `角色名不能为空`)
	}
	var exists bool
	var err error
	if u.Id > 0 {
		exists, err = u.Exists2(u.Name, u.Id)
	} else {
		exists, err = u.Exists(u.Name)
	}
	if err != nil {
		return err
	}
	if exists {
		err = u.Context().NewError(code.DataAlreadyExists, `角色名已经存在`)
	}
	return err
}

func (u *UserRole) Add() (interface{}, error) {
	if err := u.check(); err != nil {
		return nil, err
	}
	return u.NgingUserRole.Insert()
}

func (u *UserRole) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := u.check(); err != nil {
		return err
	}
	return u.NgingUserRole.Update(mw, args...)
}

func (u *UserRole) Exists(name string) (bool, error) {
	return u.NgingUserRole.Exists(nil, db.Cond{`name`: name})
}

func (u *UserRole) ListByUser(user *dbschema.NgingUser) (roleList []*dbschema.NgingUserRole) {
	if len(user.RoleIds) > 0 {
		u.ListByOffset(nil, nil, 0, -1, db.And(
			db.Cond{`disabled`: `N`},
			db.Cond{`id`: db.In(strings.Split(user.RoleIds, `,`))},
		))
		roleList = u.Objects()
	}
	return
}

func (u *UserRole) Exists2(name string, excludeID uint) (bool, error) {
	return u.NgingUserRole.Exists(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`id`: db.NotEq(excludeID)},
	))
}
