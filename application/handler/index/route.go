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

	"github.com/webx-top/com"

	"github.com/admpub/nging/v3/application/handler"
	premLib "github.com/admpub/nging/v3/application/library/perm"
	"github.com/admpub/nging/v3/application/registry/navigate"
	"github.com/admpub/nging/v3/application/registry/perm"
	"github.com/webx-top/echo"
)

func RouteList(ctx echo.Context) error {
	return ctx.JSON(handler.IRegister().Routes())
}

func NavTree(ctx echo.Context) error {
	return ctx.JSON(premLib.NavTreeCached())
}

// UnlimitedURLs 不用采用权限验证的路由
var UnlimitedURLs = []string{
	`/favicon.ico`,
	`/captcha/*`,
	`/setup`,
	`/progress`,
	`/license`,
	``,
	`/`,
	`/project/:ident`,
	`/index`,
	`/login`,
	`/register`,
	`/logout`,
	`/donation`,
	`/icon`,
	`/routeList`,
	`/routeNotin`,
	`/navTree`,
	`/gauth_check`,
	`/qrcode`,
	`/server/dynamic`,
	`/public/upload/:subdir/*`, //查看上传后的文件
	`/finder`,
	`/donation/:type`,
}

var HandlerPermissions = []string{
	`guest`,  // 游客可浏览
	`public`, // 任意登录用户可浏览
}

func RouteNotin(ctx echo.Context) error {
	var unuse []string
	for _, route := range handler.IRegister().Routes() {
		urlPath := route.Path
		if com.InSlice(urlPath, UnlimitedURLs) {
			continue
		}
		if com.InSlice(route.String(`permission`), HandlerPermissions) {
			continue
		}
		if strings.HasPrefix(urlPath, `/term/client/`) {
			continue
		}
		if strings.HasPrefix(urlPath, `/frp/dashboard/`) {
			continue
		}
		if strings.HasPrefix(urlPath, `/debug/`) {
			continue
		}
		if strings.HasPrefix(urlPath, `/captcha/`) {
			continue
		}
		if strings.HasPrefix(urlPath, `/user/`) {
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
		if com.InSlice(urlPath, unuse) {
			continue
		}
		if _, ok := perm.SpecialAuths[urlPath]; ok {
			continue
		}
		unuse = append(unuse, urlPath)
	}

	return ctx.JSON(unuse)
}
