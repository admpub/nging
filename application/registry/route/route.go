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
)

var (
	routeRegister = route.NewRegister()
)

func Register(fn func(*echo.Echo)) {
	routeRegister.Register(fn)
}

func Use(groupName string, middlewares ...interface{}) {
	routeRegister.Use(groupName, middlewares...)
}

func Apply(e *echo.Echo) {
	routeRegister.Apply(e)
}

func RegisterToGroup(groupName string, fn func(*echo.Group), middlewares ...interface{}) {
	routeRegister.RegisterToGroup(groupName, fn, middlewares...)
}
