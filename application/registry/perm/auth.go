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

package perm

import (
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/echo"
)

type AuthChecker func(
	h echo.Handler,
	c echo.Context,
	rpath string,
	user *dbschema.User,
	roleM *model.UserRole,
	roleList []*dbschema.UserRole,
) (err error, ppath string, returning bool)

var SpecialAuths = map[string]AuthChecker{
	`/server/cmdSend/*`: func(
		h echo.Handler,
		c echo.Context,
		rpath string,
		user *dbschema.User,
		roleM *model.UserRole,
		roleList []*dbschema.UserRole,
	) (err error, ppath string, returning bool) {
		returning = true
		c.SetFunc(`CheckPerm`, func(id string) error {
			if user.Id == 1 {
				return nil
			}
			if len(id) > 0 {
				if !roleM.CheckCmdPerm2(roleList, id) {
					return echo.ErrForbidden
				}
			} else {
				if !roleM.CheckPerm2(roleList, `server/cmd`) {
					return echo.ErrForbidden
				}
			}
			return nil
		})
		err = h.Handle(c)
		return
	},
	`server/dynamic`: func(
		h echo.Handler,
		c echo.Context,
		rpath string,
		user *dbschema.User,
		roleM *model.UserRole,
		roleList []*dbschema.UserRole,
	) (err error, ppath string, returning bool) {
		ppath = `server/sysinfo`
		return
	},
	`/server/cmd`: func(
		h echo.Handler,
		c echo.Context,
		rpath string,
		user *dbschema.User,
		roleM *model.UserRole,
		roleList []*dbschema.UserRole,
	) (err error, ppath string, returning bool) {
		id := c.Form(`id`)
		if len(id) > 0 {
			returning = true
			if !roleM.CheckCmdPerm2(roleList, id) {
				err = echo.ErrForbidden
				return
			}
			err = h.Handle(c)
			return
		}
		ppath = `cmd`
		return
	},
}

func init() {
	SpecialAuths[`/server/cmdSendWS`] = SpecialAuths[`/server/cmdSend/*`]
}

func AuthRegister(ppath string, checker AuthChecker) {
	SpecialAuths[ppath] = checker
}

func AuthUnregister(ppath string) {
	if _, ok := SpecialAuths[ppath]; ok {
		delete(SpecialAuths, ppath)
	}
}
