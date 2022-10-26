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

package role

import (
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/webx-top/echo"
)

type AuthChecker func(
	h echo.Handler,
	c echo.Context,
	user *dbschema.NgingUser,
	permission *RolePermission,
) (ppath string, returning bool, err error)

var SpecialAuths = map[string]AuthChecker{
	`/server/cmdSend/*`: authServerCmdSend,
	`server/dynamic`:    authServerStatus,
	`/server/cmd`:       authCmd,
	`/manager/crop`:     authCrop,
}

func authServerCmdSend(
	h echo.Handler,
	c echo.Context,
	user *dbschema.NgingUser,
	permission *RolePermission,
) (ppath string, returning bool, err error) {
	returning = true
	c.SetFunc(`CheckPerm`, func(id string) error {
		if user.Id == 1 {
			return nil
		}
		if permission == nil {
			return common.ErrUserNoPerm
		}
		if len(id) > 0 {
			if !permission.CheckCmd(c, id) {
				return common.ErrUserNoPerm
			}
		} else {
			if !permission.Check(c, `server/cmd`) {
				return common.ErrUserNoPerm
			}
		}
		return nil
	})
	err = h.Handle(c)
	return
}

func authServerStatus(
	h echo.Handler,
	c echo.Context,
	user *dbschema.NgingUser,
	permission *RolePermission,
) (ppath string, returning bool, err error) {
	ppath = `server/sysinfo`
	return
}

func authCmd(
	h echo.Handler,
	c echo.Context,
	user *dbschema.NgingUser,
	permission *RolePermission,
) (ppath string, returning bool, err error) {
	id := c.Form(`id`)
	if len(id) > 0 {
		returning = true
		if permission == nil {
			err = common.ErrUserNoPerm
			return
		}
		if !permission.CheckCmd(c, id) {
			err = common.ErrUserNoPerm
			return
		}
		err = h.Handle(c)
		return
	}
	ppath = `cmd`
	return
}

func authCrop(
	h echo.Handler,
	c echo.Context,
	user *dbschema.NgingUser,
	permission *RolePermission,
) (ppath string, returning bool, err error) {
	ppath = `/manager/upload/:type`
	return
}

func init() {
	SpecialAuths[`/server/cmdSendWS`] = SpecialAuths[`/server/cmdSend/*`]
}

func AuthRegister(ppath string, checker AuthChecker) {
	SpecialAuths[ppath] = checker
}

func AuthUnregister(ppath string) {
	delete(SpecialAuths, ppath)
}
