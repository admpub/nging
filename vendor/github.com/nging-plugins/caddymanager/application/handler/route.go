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

package handler

import (
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/route"
	"github.com/webx-top/echo"
)

func RegisterRoute(r *route.Collection) {
	r.Backend.RegisterToGroup(`/caddy`, registerRoute)
}

func registerRoute(g echo.RouteRegister) {
	metaHandler := handler.IRegister().MetaHandler
	g.Route(`GET,POST`, `/restart`, metaHandler(echo.H{`name`: `重启Caddy服务`}, Restart))
	g.Route(`GET,POST`, `/stop`, metaHandler(echo.H{`name`: `停止Caddy服务`}, Stop))
	g.Route(`GET,POST`, `/log`, metaHandler(echo.H{`name`: `显示Caddy日志`}, Log))
	g.Route(`GET`, `/addon_form`, metaHandler(echo.H{`name`: `Caddy配置表单`}, AddonForm))
	g.Route(`GET`, `/vhost`, metaHandler(echo.H{`name`: `网站列表`}, VhostIndex))
	g.Route(`GET,POST`, `/vhost_add`, metaHandler(echo.H{`name`: `添加网站`}, VhostAdd))
	g.Route(`GET,POST`, `/vhost_edit`, metaHandler(echo.H{`name`: `编辑网站`}, VhostEdit))
	g.Route(`GET,POST`, `/vhost_delete`, metaHandler(echo.H{`name`: `删除网站`}, VhostDelete))
	g.Route(`GET,POST`, `/vhost_file`, metaHandler(echo.H{`name`: `管理网站文件`}, VhostFile))
	g.Route(`GET`, `/vhost_build`, metaHandler(echo.H{`name`: `生成Caddyfile`}, Vhostbuild))
	g.Route(`GET,POST`, `/log_show`, LogShow)
	g.Route(`GET,POST`, `/vhost_log`, VhostLog)

	g.Route(`GET`, `/group`, Group)
	g.Route(`GET,POST`, `/group_add`, GroupAdd)
	g.Route(`GET,POST`, `/group_edit`, GroupEdit)
	g.Route(`GET,POST`, `/group_delete`, GroupDelete)
}
