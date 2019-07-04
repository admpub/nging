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

package caddy

import (
	"strings"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/handler/tool"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/model"
	"github.com/admpub/tail"
	ua "github.com/admpub/useragent"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/caddy`, func(g echo.RouteRegister) {
		e := handler.Echo()
		g.Route(`GET,POST`, `/restart`, e.MetaHandler(echo.H{`name`: `重启Caddy服务`}, Restart))
		g.Route(`GET,POST`, `/stop`, e.MetaHandler(echo.H{`name`: `停止Caddy服务`}, Stop))
		g.Route(`GET,POST`, `/log`, e.MetaHandler(echo.H{`name`: `显示Caddy日志`}, Log))
		g.Route(`GET`, `/addon_form`, e.MetaHandler(echo.H{`name`: `Caddy配置表单`}, AddonForm))
		g.Route(`GET`, `/vhost`, e.MetaHandler(echo.H{`name`: `网站列表`}, VhostIndex))
		g.Route(`GET,POST`, `/vhost_add`, e.MetaHandler(echo.H{`name`: `添加网站`}, VhostAdd))
		g.Route(`GET,POST`, `/vhost_edit`, e.MetaHandler(echo.H{`name`: `编辑网站`}, VhostEdit))
		g.Route(`GET,POST`, `/vhost_delete`, e.MetaHandler(echo.H{`name`: `删除网站`}, VhostDelete))
		g.Route(`GET,POST`, `/vhost_file`, e.MetaHandler(echo.H{`name`: `管理网站文件`}, VhostFile))
		g.Route(`GET`, `/vhost_build`, e.MetaHandler(echo.H{`name`: `生成Caddyfile`}, Vhostbuild))
		g.Route(`GET,POST`, `/log_show`, LogShow)
		g.Route(`GET,POST`, `/vhost_log`, VhostLog)

		g.Route(`GET`, `/group`, Group)
		g.Route(`GET,POST`, `/group_add`, GroupAdd)
		g.Route(`GET,POST`, `/group_edit`, GroupEdit)
		g.Route(`GET,POST`, `/group_delete`, GroupDelete)
	})

	common.LogParsers[`access`] = func(line *tail.Line) (interface{}, error) {
		logM := model.NewAccessLog(nil)
		err := logM.Parse(line.Text)
		res := logM.AsMap()
		realIP := logM.RemoteAddr
		if len(logM.XForwardFor) > 0 {
			realIP = strings.TrimSpace(strings.SplitN(logM.XForwardFor, ",", 2)[0])
		} else if len(logM.XRealIp) > 0 {
			realIP = logM.XRealIp
		}
		delete(res, `Id`)
		delete(res, `Created`)
		delete(res, `VhostId`)
		delete(res, `RemoteAddr`)
		delete(res, `XForwardFor`)
		delete(res, `XRealIp`)
		delete(res, `Minute`)
		delete(res, `HitStatus`)
		delete(res, `Host`)
		delete(res, `LocalAddr`)
		if ipInfo, _err := tool.IPInfo(realIP); _err == nil {
			res[`ClientRegion`] = ipInfo.Country + " - " + ipInfo.Region + " - " + ipInfo.Province + " - " + ipInfo.City + " " + ipInfo.ISP + " (" + realIP + ")"
		} else {
			res[`ClientRegion`] = realIP
		}
		infoUA := ua.Parse(logM.UserAgent)
		if len(infoUA.OSVersion) > 0 {
			infoUA.OS += ` (` + infoUA.OSVersion + `)`
		}
		res[`OS`] = infoUA.OS
		if len(infoUA.Version) > 0 {
			logM.BrowerName += ` (` + infoUA.Version + `)`
		}
		res[`BrowerName`] = logM.BrowerName
		return res, err
	}
}
