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
package index

import (
	"strings"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/registry/navigate"
	"github.com/admpub/nging/application/registry/perm"
	"github.com/webx-top/echo"
)

func RouteList(ctx echo.Context) error {
	return ctx.JSON(handler.Echo().Routes())
}

func RouteNotin(ctx echo.Context) error {
	var unuse []string
	for _, route := range handler.Echo().Routes() {
		if strings.HasPrefix(route.Path, `/term/client/`) {
			continue
		}
		if strings.HasPrefix(route.Path, `/frp/dashboard/`) {
			continue
		}
		var exists bool
		for _, navGroup := range navigate.TopNavigate {
			for _, navItem := range navGroup.Children {
				var navRoute string
				if len(navItem.Action) > 0 {
					navRoute = `/` + navGroup.Action + `/` + navItem.Action
				} else {
					navRoute = `/` + navGroup.Action
				}
				if navRoute == route.Path {
					exists = true
					break
				}
			}
		}
		if exists {
			continue
		}
		for _, navGroup := range navigate.LeftNavigate {
			for _, navItem := range navGroup.Children {
				var navRoute string
				if len(navItem.Action) > 0 {
					navRoute = `/` + navGroup.Action + `/` + navItem.Action
				} else {
					navRoute = `/` + navGroup.Action
				}
				if navRoute == route.Path {
					exists = true
					break
				}
			}
		}
		if exists {
			continue
		}
		for _, v := range unuse {
			if v == route.Path {
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		_, exists = perm.SpecialAuths[route.Path]
		if exists {
			continue
		}
		unuse = append(unuse, route.Path)
	}

	return ctx.JSON(unuse)
}
