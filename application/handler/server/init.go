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
package server

import (
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/webx-top/echo"
	sockjsHandler "github.com/webx-top/echo/handler/sockjs"
	ws "github.com/webx-top/echo/handler/websocket"
)

func init() {
	handler.RegisterToGroup(`/server`, func(g echo.RouteRegister) {
		g.Route("GET", `/sysinfo`, Info)
		g.Route("GET", `/netstat`, Connections)
		g.Route("GET", `/process/:pid`, ProcessInfo)
		g.Route("GET", `/procskill/:pid`, ProcessKill)
		g.Route(`GET,POST`, `/service`, Service)
		g.Route(`GET,POST`, `/daemon_index`, DaemonIndex)
		g.Route(`GET,POST`, `/daemon_add`, DaemonAdd)
		g.Route(`GET,POST`, `/daemon_edit`, DaemonEdit)
		g.Route(`GET,POST`, `/daemon_delete`, DaemonDelete)
		g.Route("GET", `/cmd`, Cmd)
		g.Route(`GET,POST`, `/daemon_log`, DaemonLog)
		g.Route(`GET,POST`, `/log`, func(c echo.Context) error {
			return config.DefaultConfig.Log.Show(c)
		})
		sockjsOpts := sockjsHandler.Options{
			Handle: CmdSendBySockJS,
			Prefix: "/cmdSend",
		}
		//sockjsOpts.Wrapper(g)
		_ = sockjsOpts
		wsOpts := ws.Options{
			Handle: CmdSendByWebsocket,
			Prefix: "/cmdSendWS",
		}
		wsOpts.Wrapper(g)

		wsOptsDynamicInfo := ws.Options{
			Handle: InfoByWebsocket,
			Prefix: "/dynamic",
		}
		wsOptsDynamicInfo.Wrapper(g)
	})
	ListenRealTimeStatus()
}
