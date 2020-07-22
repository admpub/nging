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
package model

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/perm"
	"github.com/admpub/nging/application/model/base"
	"github.com/admpub/nging/application/registry/navigate"
)

func NewUserRole(ctx echo.Context) *UserRole {
	return &UserRole{
		NgingUserRole: &dbschema.NgingUserRole{},
		base:          base.New(ctx),
	}
}

type UserRole struct {
	*dbschema.NgingUserRole
	base        *base.Base
	permActions *perm.Map
	permCmds    *perm.Map
}

func (u *UserRole) check() error {
	if len(u.Name) == 0 {
		return u.base.NewError(code.InvalidParameter, `角色名不能为空`)
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
		err = u.base.NewError(code.DataAlreadyExists, `角色名已经存在`)
	}
	return err
}

func (u *UserRole) Add() (interface{}, error) {
	if err := u.check(); err != nil {
		return nil, err
	}
	return u.NgingUserRole.Add()
}

func (u *UserRole) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := u.check(); err != nil {
		return err
	}
	return u.NgingUserRole.Edit(mw, args...)
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
	}
	roleList = u.Objects()
	return
}

func (u *UserRole) CheckPerm2(roleList []*dbschema.NgingUserRole, perm string) (hasPerm bool) {
	perm = strings.TrimPrefix(perm, `/`)
	for _, role := range roleList {
		r := NewUserRole(nil)
		r.NgingUserRole = role
		if r.CheckPerm(perm) {
			hasPerm = true
			break
		}
	}
	return
}

func (u *UserRole) Exists2(name string, excludeID uint) (bool, error) {
	return u.NgingUserRole.Exists(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`id`: db.NotEq(excludeID)},
	))
}

func (u *UserRole) CleanPermAction(values []string) *UserRole {
	u.PermAction = perm.BuildPermActions(values)
	return u
}

func (u *UserRole) CheckPerm(permPath string) bool {
	if len(u.PermAction) == 0 {
		return false
	}
	if u.PermAction == `*` {
		return true
	}
	navTree := perm.NavTreeCached()
	if u.permActions == nil {
		u.permActions = perm.NewMap()
		u.permActions.Parse(u.PermAction, navTree)
	}
	return u.permActions.Check(permPath, navTree)
}

func (u *UserRole) CheckCmdPerm2(roleList []*dbschema.NgingUserRole, perm string) (hasPerm bool) {
	for _, role := range roleList {
		r := NewUserRole(nil)
		r.NgingUserRole = role
		if r.CheckCmdPerm(perm) {
			hasPerm = true
			break
		}
	}
	return
}

func (u *UserRole) CheckCmdPerm(permPath string) bool {
	if len(u.PermCmd) == 0 {
		return false
	}
	if u.PermCmd == `*` {
		return true
	}
	if u.permCmds == nil {
		u.permCmds = perm.NewMap()
		perms := strings.Split(u.PermCmd, `,`)
		for _, _perm := range perms {
			arr := strings.Split(_perm, `/`)
			result := u.permCmds.V
			for _, a := range arr {
				if _, y := result[a]; !y {
					result[a] = perm.NewMap()
				}
				result = result[a].V
			}
		}
	}

	arr := strings.Split(permPath, `/`)
	result := u.permCmds.V
	for _, a := range arr {
		v, y := result[a]
		if !y {
			return false
		}
		result = v.V
	}
	return true
}

//FilterNavigate 过滤导航菜单，只显示有权限的菜单
func (u *UserRole) FilterNavigate(roleList []*dbschema.NgingUserRole, navList *navigate.List) navigate.List {
	var result navigate.List
	if navList == nil {
		return result
	}
	for _, nav := range *navList {
		if !nav.Unlimited && !u.CheckPerm2(roleList, nav.Action) {
			continue
		}
		navCopy := *nav
		navCopy.Children = &navigate.List{}
		for _, child := range *nav.Children {
			var perm string
			if len(child.Action) > 0 {
				perm = nav.Action + `/` + child.Action
			} else {
				perm = nav.Action
			}
			if !nav.Unlimited && !u.CheckPerm2(roleList, perm) {
				continue
			}
			*navCopy.Children = append(*navCopy.Children, child)
		}
		result = append(result, &navCopy)
	}
	return result
}
