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
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"golang.org/x/net/publicsuffix"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/nging-plugins/caddymanager/application/library/form"
	"github.com/nging-plugins/caddymanager/application/model"
)

func ServerIndex(ctx echo.Context) error {
	m := model.NewVhostServer(ctx)
	cond := db.NewCompounds()
	err := m.ListPage(cond)
	ctx.Set(`listData`, m.Objects())
	ctx.SetFunc(`engineName`, engine.Engines.Get)
	ctx.SetFunc(`environName`, engine.Environs.Get)
	return ctx.Render(`caddy/server`, handler.Err(ctx, err))
}

func ServerAdd(ctx echo.Context) error {
	var err error
	m := model.NewVhostServer(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingVhostServer, echo.ExcludeFieldName(`created`, `updated`, `configFileUpdated`))
		if err != nil {
			goto END
		}
		_, err = m.Add()
		if err != nil {
			goto END
		}
		if m.AutoModifyConfig == common.BoolY {
			var cfg engine.Configer
			var changed bool
			cfg, err = getServerConfig(ctx, m.Ident, m.NgingVhostServer)
			if err != nil {
				goto END
			}
			changed, err = engine.FixEngineConfigFile(cfg)
			if err != nil {
				goto END
			}
			if changed {
				err = m.SetConfigFileUpdated(uint(time.Now().Unix()))
				if err != nil {
					goto END
				}
			}
		}
		handler.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(handler.URLFor(`/caddy/server`))
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, db.Cond{`id`: id})
			if err == nil {
				echo.StructToForm(ctx, m.NgingVhostServer, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

END:
	ctx.Set(`activeURL`, `/caddy/server`)
	ctx.Set(`isAdd`, true)
	setServerForm(ctx)
	return ctx.Render(`caddy/server_edit`, handler.Err(ctx, err))
}

func setServerForm(ctx echo.Context) {
	thirdpartyEngines := engine.Thirdparty()
	ctx.Set(`engineList`, thirdpartyEngines)
	configDirs := map[string]string{}
	for _, eng := range thirdpartyEngines {
		configDirs[eng.K] = eng.X.(engine.Enginer).DefaultConfigDir()
	}
	ctx.Set(`configDirs`, configDirs)
	ctx.Set(`environList`, engine.Environs.Slice())
}

func ServerEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewVhostServer(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/caddy/server`))
	}
	if ctx.IsPost() {
		old := *m.NgingVhostServer
		err = ctx.MustBind(m.NgingVhostServer, echo.ExcludeFieldName(`created`, `updated`, `engine`, `ident`, `configFileUpdated`))
		if err != nil {
			goto END
		}
		m.Id = id
		err = m.Edit(nil, `id`, id)
		if err != nil {
			goto END
		}
		if old.Disabled != m.Disabled {
			if m.Disabled == `Y` {
				err = deleteCaddyfileByServer(ctx, &old, true)
			} else {
				err = vhostbuild(ctx, 0, ``, ``, m.NgingVhostServer)
			}
			if err != nil {
				ctx.Logger().Error(err)
			}
		} else if old.VhostConfigLocalDir != m.VhostConfigLocalDir {
			err = deleteCaddyfileByServer(ctx, &old, true)
			if err != nil {
				ctx.Logger().Error(err)
			}
			err = vhostbuild(ctx, 0, ``, ``, m.NgingVhostServer)
			if err != nil {
				ctx.Logger().Error(err)
			}
		}
		var cfg engine.Configer
		var changed bool
		if old.ConfigFileUpdated > 0 || m.AutoModifyConfig == common.BoolY {
			cfg, err = getServerConfig(ctx, m.Ident, m.NgingVhostServer)
			if err != nil {
				goto END
			}
		}
		if old.ConfigFileUpdated > 0 { // 以前添加过
			if old.ConfigLocalFile != m.ConfigLocalFile || old.VhostConfigLocalDir != m.VhostConfigLocalDir {
				// ==================================
				//  从旧文件中删除 || 删除旧的引用
				// ==================================
				var oldCfg engine.Configer
				var oldChanged bool
				oldCfg, err = getServerConfig(ctx, old.Ident, &old)
				if err != nil {
					goto END
				}
				oldChanged, err = engine.FixEngineConfigFile(oldCfg, true)
				if err != nil {
					goto END
				}
				_ = oldChanged

				if m.AutoModifyConfig == common.BoolY { // 添加新的
					// 在新文件中添加 || 添加新的引用
					changed, err = engine.FixEngineConfigFile(cfg)
				} else { // 没有添加，添加时间重置为0
					err = m.SetConfigFileUpdated(0)
				}
			}
		} else if m.AutoModifyConfig == common.BoolY {
			changed, err = engine.FixEngineConfigFile(cfg)
		}
		if err != nil {
			goto END
		}
		if changed {
			err = m.SetConfigFileUpdated(uint(time.Now().Unix()))
			if err != nil {
				goto END
			}
		}
		handler.SendOk(ctx, ctx.T(`修改成功`))
		return ctx.Redirect(handler.URLFor(`/caddy/server`))
	} else if ctx.IsAjax() {
		data := ctx.Data()
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			if !common.IsBoolFlag(disabled) {
				return ctx.NewError(code.InvalidParameter, ``).SetZone(`disabled`)
			}
			m.Disabled = disabled
			err = m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			if m.Disabled == `Y` {
				err = deleteCaddyfileByServer(ctx, m.NgingVhostServer, true)
			} else {
				err = vhostbuild(ctx, 0, ``, ``, m.NgingVhostServer)
			}
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
		}
		return ctx.JSON(data)
	}
	echo.StructToForm(ctx, m.NgingVhostServer, ``, echo.LowerCaseFirstLetter)

END:
	ctx.Set(`activeURL`, `/caddy/server`)
	ctx.Set(`isAdd`, false)
	setServerForm(ctx)
	return ctx.Render(`caddy/server_edit`, handler.Err(ctx, err))
}

func ServerDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewVhostServer(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		goto END
	}
	if m.AutoModifyConfig == common.BoolY && m.ConfigFileUpdated > 0 {
		var cfg engine.Configer
		cfg, err = getServerConfig(ctx, m.Ident, m.NgingVhostServer)
		if err != nil {
			goto END
		}
		_, err = engine.FixEngineConfigFile(cfg, true)
		if err != nil {
			goto END
		}
	}
	err = m.Delete(nil, db.Cond{`id`: id})
	if err != nil {
		goto END
	}
	err = deleteCaddyfileByServer(ctx, m.NgingVhostServer, true)
	if err != nil {
		goto END
	}

END:
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}
	return ctx.Redirect(handler.URLFor(`/caddy/server`))
}

func ServerRenewalCert(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		return ctx.String(code.InvalidParameter.String(), http.StatusBadRequest)
	}
	m := model.NewVhostServer(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return ctx.String(err.Error(), http.StatusNotFound)
		}
		return ctx.String(err.Error(), http.StatusInternalServerError)
	}
	cfg, err := getServerConfig(ctx, m.Ident, m.NgingVhostServer)
	if err != nil {
		return ctx.String(err.Error(), http.StatusInternalServerError)
	}
	renew, ok := cfg.(engine.CertRenewaler)
	if !ok {
		return ctx.String(code.Unsupported.String(), http.StatusNotImplemented)
	}
	vhostM := model.NewVhost(ctx)
	conds := db.And(
		db.Cond{`server_ident`: m.Ident},
		db.Cond{`disabled`: common.BoolN},
	)
	_, err = vhostM.ListByOffset(nil, nil, 0, -1, conds)
	if err != nil {
		return ctx.String(err.Error(), http.StatusInternalServerError)
	}
	var updateCount int
	for _, row := range vhostM.Objects() {
		updatedDomains, err := renewalVhostCert(ctx, renew, row)
		if err != nil {
			ctx.Logger().Error(err.Error())
		}
		updateCount += len(updatedDomains)
	}
	if updateCount > 0 {
		item := engine.Engines.GetItem(cfg.GetEngine())
		if item != nil {
			err = item.X.(engine.Enginer).ReloadServer(ctx, cfg)
		}
	}
	if err != nil {
		return ctx.String(err.Error(), http.StatusInternalServerError)
	}
	return ctx.String(ctx.T(`更新成功: %d 个`, updateCount))
}

type DomainData struct {
	ID    uint
	Email string
}

type RequestCertUpdate struct {
	Domains []string
	Email   string
}

func supportedAutoSSL(formData url.Values) bool {
	enabledTLS := formData.Get(`tls`) == `1`
	if !enabledTLS {
		return false
	}
	tlsEmail := formData.Get(`tls_email`)
	if len(tlsEmail) == 0 {
		return false
	}
	if len(formData.Get(`tls_cert`)) > 0 && len(formData.Get(`tls_key`)) > 0 {
		return false
	}
	return true
}

func renewalVhostCert(ctx echo.Context, renew engine.CertRenewaler, row *dbschema.NgingVhost, formData ...url.Values) (updatedDomains []string, err error) {
	var idDomains map[uint]*RequestCertUpdate
	idDomains, err = parseVhostEnabledHTTPSDomains(ctx, row, formData...)
	if err != nil {
		return
	}
	for id, req := range idDomains {
		err = renew.RenewalCert(ctx, id, req.Domains, req.Email)
		if err != nil {
			ctx.Logger().Error(err.Error())
		} else {
			updatedDomains = append(updatedDomains, req.Domains...)
		}
	}
	return
}

func parseVhostEnabledHTTPSDomains(ctx echo.Context, row *dbschema.NgingVhost, _formData ...url.Values) (idDomains map[uint]*RequestCertUpdate, err error) {
	var formData url.Values
	if len(_formData) > 0 {
		formData = _formData[0]
	} else {
		jsonBytes := []byte(row.Setting)
		err = json.Unmarshal(jsonBytes, &formData)
		if err != nil {
			return nil, common.JSONBytesParseError(err, jsonBytes)
		}
	}
	if !supportedAutoSSL(formData) {
		return nil, nil
	}
	tlsEmail := formData.Get(`tls_email`)
	domainAndEmails := map[string]DomainData{}
	for _, domain := range form.SplitBySpace(row.Domain) {
		domain = com.ParseEnvVar(domain)
		parts := strings.SplitN(domain, `://`, 2)
		var scheme, host, port string
		if len(parts) == 2 {
			scheme = parts[1]
			host, port = com.SplitHostPort(parts[1])
		} else {
			host, port = com.SplitHostPort(domain)
		}
		if len(host) == 0 || com.IsLocalhost(host) {
			continue
		}
		if _, derr := publicsuffix.EffectiveTLDPlusOne(host); derr != nil {
			ctx.Logger().Error(derr.Error())
			continue
		}
		if port == `443` || scheme == `https` {
			domainAndEmails[host] = DomainData{ID: row.Id, Email: tlsEmail}
		}
	}
	idDomains = map[uint]*RequestCertUpdate{}
	for domain, info := range domainAndEmails {
		if _, ok := idDomains[info.ID]; !ok {
			idDomains[info.ID] = &RequestCertUpdate{
				Email: info.Email,
			}
		}
		idDomains[info.ID].Domains = append(idDomains[info.ID].Domains, domain)
	}
	return
}
