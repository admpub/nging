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
package user

import (
	"github.com/admpub/nging/v4/application/handler"
	"github.com/webx-top/echo"
	ws "github.com/webx-top/echo/handler/websocket"
)

func init() {
	handler.RegisterToGroup(`/user`, func(g echo.RouteRegister) {
		metaHandler := handler.IRegister().MetaHandler
		g.Route("GET,POST", `/edit`, metaHandler(echo.H{`name`: `修改个人资料`}, Edit))
		g.Route("GET,POST", `/gauth_bind`, metaHandler(echo.H{`name`: `绑定两步验证`}, GAuthBind))
		g.Route("GET,POST", `/autocomplete_path`, AutoCompletePath)
		wsOpts := ws.Options{
			Handle: Notice,
			Prefix: "/notice",
		}
		wsOpts.Wrapper(g)
	})
}
