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

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
)

func NewForeverProcess(ctx echo.Context) *ForeverProcess {
	return &ForeverProcess{
		ForeverProcess: &dbschema.ForeverProcess{},
		Base:           base.New(ctx),
	}
}

type ForeverProcess struct {
	*dbschema.ForeverProcess
	*base.Base
}

func (u *ForeverProcess) Exists(name string) (bool, error) {
	n, e := u.Param(nil, db.Cond{`name`: name}).Count()
	return n > 0, e
}

func (u *ForeverProcess) Add() (pk interface{}, err error) {
	u.Status = `idle`
	return u.ForeverProcess.Add()
}

func (u *ForeverProcess) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return u.ForeverProcess.Edit(mw, args...)
}

func (u *ForeverProcess) Exists2(name string, excludeID uint) (bool, error) {
	n, e := u.Param(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`id`: db.NotEq(excludeID)},
	)).Count()
	return n > 0, e
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
