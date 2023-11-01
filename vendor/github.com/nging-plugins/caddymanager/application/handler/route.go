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
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/cron"
	"github.com/admpub/nging/v5/application/library/route"
	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/nging-plugins/caddymanager/application/model"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
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

	g.Route(`GET`, `/server`, metaHandler(echo.H{`name`: `引擎配置列表`}, ServerIndex))
	g.Route(`GET,POST`, `/server_add`, metaHandler(echo.H{`name`: `添加引擎配置`}, ServerAdd))
	g.Route(`GET,POST`, `/server_edit`, metaHandler(echo.H{`name`: `编辑引擎配置`}, ServerEdit))
	g.Route(`GET,POST`, `/server_delete`, metaHandler(echo.H{`name`: `删除引擎配置`}, ServerDelete))
	g.Route(`GET,POST`, `/server_renew_cert`, metaHandler(echo.H{`name`: `更新HTTPS证书`}, ServerRenewCert))

	g.Route(`GET`, `/group`, Group)
	g.Route(`GET,POST`, `/group_add`, GroupAdd)
	g.Route(`GET,POST`, `/group_edit`, GroupEdit)
	g.Route(`GET,POST`, `/group_delete`, GroupDelete)
}

func init() {
	cron.Register(`renewVhostCert`, renewVhostCertJob, `>renewVhostCert:1`, `更新虚拟机SSL证书`)
}

// renewVhostCertJob 更新虚拟机SSL证书
func renewVhostCertJob(idString string) cron.Runner {
	params := url.Values{}
	var forceObtain bool
	var id uint
	arr := strings.SplitN(idString, "?", 2) // id or id?forceObtain=1
	if len(arr) == 2 {
		var err error
		params, err = url.ParseQuery(arr[1])
		if err != nil {
			log.Error(err)
		}
		forceObtain = param.AsBool(params.Get(`forceObtain`))
		id = param.AsUint(arr[0])
	}
	return func(timeout time.Duration) (out string, runingErr string, onErr error, isTimeout bool) {
		ctx := defaults.NewMockContext()
		m := model.NewVhost(ctx)
		err := m.Get(nil, `id`, id)
		if err != nil {
			onErr = cron.ErrFailure
			runingErr = err.Error()
			return
		}
		if m.Disabled == common.BoolY {
			return
		}
		cfg, err := getServerConfig(ctx, m.ServerIdent)
		if err != nil || cfg == nil {
			return
		}
		renew, ok := cfg.(engine.CertRenewaler)
		if !ok {
			out = ctx.T(`服务器软件支持自动更新SSL证书，无需再通过自建任务来更新`)
			return
		}
		var formData url.Values
		jsonBytes := []byte(m.Setting)
		err = json.Unmarshal(jsonBytes, &formData)
		if err != nil {
			err = common.JSONBytesParseError(err, jsonBytes)
			onErr = cron.ErrFailure
			runingErr = err.Error()
			return
		}
		httpsDomains, err := renewVhostCert(ctx, renew, m.NgingVhost, forceObtain, formData)
		if err != nil {
			if !errors.Is(err, engine.ErrNotSetCertContainerDir) && !errors.Is(err, engine.ErrNotSetCertLocalDir) {
				return
			}
			onErr = cron.ErrFailure
			runingErr = err.Error()
			return
		}
		if len(httpsDomains) == 0 {
			out = ctx.T(`没有更新任何域名证书`)
			return
		}
		err = setCertPathForDomains(ctx, cfg, formData, httpsDomains)
		if err == nil {
			err = saveVhostConf(ctx, cfg, m.Id, formData)
		}
		if err != nil {
			onErr = cron.ErrFailure
			runingErr = err.Error()
			return
		}
		out = ctx.T(`成功更新了SSL证书，域名：%s`, strings.Join(httpsDomains, `, `))
		return
	}
}
