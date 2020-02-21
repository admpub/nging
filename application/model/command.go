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
	"os"
	"os/exec"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
)

func NewCommand(ctx echo.Context) *Command {
	return &Command{
		NgingCommand: &dbschema.NgingCommand{},
		Base:         base.New(ctx),
	}
}

type Command struct {
	*dbschema.NgingCommand
	*base.Base
}

func (u *Command) Exists(name string) (bool, error) {
	n, e := u.Param(nil, db.Cond{`name`: name}).Count()
	return n > 0, e
}

func (u *Command) Exists2(name string, excludeID uint) (bool, error) {
	n, e := u.Param(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`id`: db.NotEq(excludeID)},
	)).Count()
	return n > 0, e
}

func (u *Command) CreateCmd() *exec.Cmd {
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
	cmd := exec.Command(u.NgingCommand.Command)
	cmd.Dir = u.WorkDirectory
	cmd.Env = append(os.Environ(), env...)
	//cmd.Stdout = bufOut
	//cmd.Stderr = bufErr
	//return cmd.Run()
	return cmd
}
