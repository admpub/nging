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

package handler

import (
	"github.com/admpub/nging/v4/application/library/route"
	"github.com/webx-top/echo"
)

func RegisterRoute(r *route.Collection) {
	r.Backend.RegisterToGroup(`/ftp`, registerRoute)
}

func registerRoute(g echo.RouteRegister) {
	g.Route(`GET`, `/account`, AccountIndex)
	g.Route(`GET,POST`, `/account_add`, AccountAdd)
	g.Route(`GET,POST`, `/account_edit`, AccountEdit)
	g.Route(`GET,POST`, `/account_delete`, AccountDelete)

	g.Route(`GET`, `/group`, GroupIndex)
	g.Route(`GET,POST`, `/group_add`, GroupAdd)
	g.Route(`GET,POST`, `/group_edit`, GroupEdit)
	g.Route(`GET,POST`, `/group_delete`, GroupDelete)
	g.Route(`GET,POST`, `/restart`, Restart)
	g.Route(`GET,POST`, `/stop`, Stop)
	g.Route(`GET,POST`, `/log`, Log)
}
