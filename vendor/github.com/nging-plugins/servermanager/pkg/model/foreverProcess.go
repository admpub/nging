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

	"github.com/nging-plugins/servermanager/pkg/dbschema"
)

func NewForeverProcess(ctx echo.Context) *ForeverProcess {
	return &ForeverProcess{
		NgingForeverProcess: dbschema.NewNgingForeverProcess(ctx),
	}
}

type ForeverProcess struct {
	*dbschema.NgingForeverProcess
}

func (u *ForeverProcess) Exists(name string) (bool, error) {
	return u.NgingForeverProcess.Exists(nil, db.Cond{`name`: name})
}

func (u *ForeverProcess) check() error {
	if len(u.Name) == 0 {
		return u.Context().NewError(code.InvalidParameter, `名称不能为空`).SetZone(`name`)
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
		return u.Context().NewError(code.DataAlreadyExists, `名称已经存在`).SetZone(`name`)
	}
	return err
}

func (u *ForeverProcess) Add() (pk interface{}, err error) {
	if err := u.check(); err != nil {
		return nil, err
	}
	u.Status = `idle`
	return u.NgingForeverProcess.Insert()
}

func (u *ForeverProcess) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := u.check(); err != nil {
		return err
	}
	return u.NgingForeverProcess.Update(mw, args...)
}

func (u *ForeverProcess) Exists2(name string, excludeID uint) (bool, error) {
	return u.NgingForeverProcess.Exists(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`id`: db.NotEq(excludeID)},
	))
}

func (u *ForeverProcess) EnvSlice() []string {
	var env []string
	u.Env = strings.TrimSpace(u.Env)
	if len(u.Env) > 0 {
		for _, row := range strings.Split(u.Env, "\n") {
			row = strings.TrimSpace(row)
			if len(row) > 0 {
				env = append(env, row)
			}
		}
	}
	return env
}

func (u *ForeverProcess) ArgsSlice() []string {
	var args []string
	u.Args = strings.TrimSpace(u.Args)
	if len(u.Args) > 0 {
		for _, row := range strings.Split(u.Args, "\n") {
			row = strings.TrimSpace(row)
			if len(row) > 0 {
				args = append(args, row)
			}
		}
	}
	return args
}
