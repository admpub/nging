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

import "github.com/webx-top/echo"

const (
	MetaKeyPermission = `permission`
	PermissionPublic  = `public` // 所有登录用户
	PermissionGuest   = `guest`  // 匿名
)

type GroupSetMetaKV interface {
	SetMetaKV(string, interface{}) *echo.Group
}

type RouteSetMetaKV interface {
	SetMetaKV(string, interface{}) echo.IRouter
}

func SetMetaPermissionPublic(s RouteSetMetaKV) echo.IRouter {
	return s.SetMetaKV(MetaKeyPermission, PermissionPublic)
}

func SetMetaPermissionGuest(s RouteSetMetaKV) echo.IRouter {
	return s.SetMetaKV(MetaKeyPermission, PermissionGuest)
}

func SetGroupMetaPermissionPublic(s GroupSetMetaKV) *echo.Group {
	return s.SetMetaKV(MetaKeyPermission, PermissionPublic)
}

func SetGroupMetaPermissionGuest(s GroupSetMetaKV) *echo.Group {
	return s.SetMetaKV(MetaKeyPermission, PermissionGuest)
}

func PublicHandler(h interface{}, meta ...echo.H) echo.Handler {
	var m echo.H
	if len(meta) > 0 && meta[0] != nil {
		m = meta[0]
	} else {
		m = echo.H{}
	}
	m.Set(MetaKeyPermission, PermissionPublic)
	return routeRegister.MetaHandler(m, h)
}

func GuestHandler(h interface{}, meta ...echo.H) echo.Handler {
	var m echo.H
	if len(meta) > 0 && meta[0] != nil {
		m = meta[0]
	} else {
		m = echo.H{}
	}
	m.Set(MetaKeyPermission, PermissionGuest)
	return routeRegister.MetaHandler(m, h)
}
