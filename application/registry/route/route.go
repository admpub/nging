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

package route

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/logger"

	"github.com/admpub/nging/v5/application/library/route"
)

var (
	routeRegister = route.NewRegister(defaults.Default)
)

func init() {
	route.Default.Backend = routeRegister
}

func IRegister() route.IRegister {
	return routeRegister
}

func Routes() []*echo.Route {
	return routeRegister.Routes()
}

func Logger() logger.Logger {
	return routeRegister.Logger()
}

func Pre(middlewares ...interface{}) {
	routeRegister.Pre(middlewares...)
}

func PreToGroup(groupName string, middlewares ...interface{}) {
	routeRegister.PreToGroup(groupName, middlewares...)
}

func Use(middlewares ...interface{}) {
	routeRegister.Use(middlewares...)
}

func UseToGroup(groupName string, middlewares ...interface{}) {
	routeRegister.UseToGroup(groupName, middlewares...)
}

func AddGroupNamer(namers ...func(string) string) {
	routeRegister.AddGroupNamer(namers...)
}

func Register(fn func(echo.RouteRegister)) {
	routeRegister.Register(fn)
}

func SetRootGroup(groupName string) {
	routeRegister.SetRootGroup(groupName)
}

func Host(hostName string, middlewares ...interface{}) *route.Host {
	return routeRegister.Host(hostName, middlewares...)
}

func Apply() {
	echo.PanicIf(echo.Fire(`nging.route.apply.before`))
	routeRegister.Apply()
	echo.PanicIf(echo.Fire(`nging.route.apply.after`))
}

func RegisterToGroup(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
	routeRegister.RegisterToGroup(groupName, fn, middlewares...)
}

func PublicHandler(h interface{}) echo.Handler {
	return routeRegister.MetaHandler(echo.H{`permission`: `public`}, h)
}

func GuestHandler(h interface{}) echo.Handler {
	return routeRegister.MetaHandler(echo.H{`permission`: `guest`}, h)
}
