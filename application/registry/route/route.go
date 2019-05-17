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

package route

import (
	"github.com/admpub/nging/application/library/route"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

var (
	routeRegister = route.NewRegister(defaults.Default)
)

func Echo() *echo.Echo {
	return routeRegister.Echo
}

func AddGroupNamer(namers ...func(string) string) {
	routeRegister.AddGroupNamer(namers...)
}

func Register(fn func(echo.RouteRegister)) {
	routeRegister.Register(fn)
}

func Use(groupName string, middlewares ...interface{}) {
	routeRegister.Use(groupName, middlewares...)
}

func SetRootGroup(groupName string) {
	routeRegister.RootGroup = groupName
}

func Apply() {
	routeRegister.Apply()
}

func RegisterToGroup(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
	routeRegister.RegisterToGroup(groupName, fn, middlewares...)
}
