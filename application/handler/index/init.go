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

package index

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/middleware"
	"github.com/admpub/nging/v4/application/registry/dashboard"
)

func init() {
	handler.Register(func(e echo.RouteRegister) {
		e.Route("GET", ``, Index)
		e.Route("GET", `/`, Index)
		e.Route("GET", `/project/:ident`, Project)
		e.Route("GET", `/index`, Index)
		e.Route("GET,POST", `/login`, Login)
		e.Route("GET,POST", `/register`, Register)
		e.Route("GET", `/logout`, Logout)
		if dashboard.TopButtonFindTmpl(`manager/topbutton/donation`) > -1 {
			e.Route("GET", `/donation/:type`, Donation)
		}
		//e.Route(`GET,POST`, `/ping`, Ping)
		e.Get(`/icon`, Icon, middleware.AuthCheck)
		e.Get(`/routeList`, RouteList, middleware.AuthCheck)
		e.Get(`/routeNotin`, RouteNotin, middleware.AuthCheck)
		e.Get(`/navTree`, NavTree, middleware.AuthCheck)
	})
}
