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
	"strings"

	"github.com/webx-top/com"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/nerrors"
	premLib "github.com/coscms/webcore/library/perm"
	"github.com/coscms/webcore/library/role"
	"github.com/coscms/webcore/registry/navigate"
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func RouteList(ctx echo.Context) error {
	return ctx.JSON(route.IRegister().Routes())
}

func NavTree(ctx echo.Context) error {
	return ctx.JSON(premLib.NavTreeCached())
}

func Headers(ctx echo.Context) error {
	user := backend.User(ctx)
	if user == nil {
		return nerrors.ErrUserNotLoggedIn
	}
	if !role.IsFounder(user) {
		return nerrors.ErrUserNoPerm.SetMessage(ctx.T(`此功能仅供网站创始人查看`))
	}
	headers := ctx.Request().Header().Std()
	headers.Del(`Cookie`)
	return ctx.JSON(headers)
}

// UnlimitedURLs 不用采用权限验证的路由前缀
var UnlimitedURLPrefixes = []string{
	`/user/`,
}

// UnlimitedURLs 不用采用权限验证的路由
var UnlimitedURLs = []string{
	`/public/upload/:subdir/*`, //查看上传后的文件
}

var HandlerPermissions = []string{
	httpserver.PermissionGuest,  // 游客可浏览
	httpserver.PermissionPublic, // 任意登录用户可浏览
}

func RouteNotin(ctx echo.Context) error {
	var unuse []string
	prefix := ctx.Echo().Prefix()
	prefixLen := len(prefix)
	hasPre := prefixLen > 0
	for _, route := range route.IRegister().Routes() {
		urlPath := route.Path
		length := len(urlPath)
		if hasPre && length > prefixLen {
			urlPath = `/` + strings.TrimPrefix(urlPath, prefix+`/`)
		}
		if length <= 1 {
			continue
		}
		if !strings.Contains(urlPath[1:], `/`) {
			continue
		}
		if com.InSlice(urlPath, UnlimitedURLs) {
			continue
		}
		if com.InSlice(route.String(`permission`), HandlerPermissions) {
			continue
		}
		var found bool
		for _, prefix := range UnlimitedURLPrefixes {
			if strings.HasPrefix(urlPath, prefix) {
				found = true
				break
			}
		}
		if found {
			continue
		}
		if _, ok := navigate.TopNavURLs()[strings.TrimPrefix(urlPath, `/`)]; ok {
			continue
		}
		if urlPath == `/download/` {
			urlPath = `/download/index.html`
		}
		if ident := navigate.ProjectIdent(urlPath); len(ident) > 0 {
			continue
		}

		for _, _urlPath := range role.GetDependency(urlPath) {
			if _, ok := navigate.TopNavURLs()[strings.TrimPrefix(_urlPath, `/`)]; ok {
				found = true
				break
			}
			if ident := navigate.ProjectIdent(_urlPath); len(ident) > 0 {
				found = true
				break
			}
		}
		if found {
			continue
		}

		if com.InSlice(urlPath, unuse) {
			continue
		}
		if _, ok := role.SpecialAuths[urlPath]; ok {
			continue
		}
		unuse = append(unuse, urlPath)
	}

	return ctx.JSON(unuse)
}
