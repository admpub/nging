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

package frp

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/admpub/frp/models/config"
	"github.com/admpub/frp/models/consts"
	"github.com/admpub/frp/utils/util"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func Table2Config(cc *dbschema.FrpClient) (hash echo.H, err error) {
	if len(cc.Extra) == 0 {
		return
	}
	data := url.Values{}
	err = json.Unmarshal([]byte(cc.Extra), &data)
	if err != nil {
		return
	}
	visitor, proxy, erro := ProxyConfigFromForm(cc.User, data)
	if erro != nil {
		err = erro
		return
	}
	hash = echo.H{}
	hash["visitor"] = visitor
	hash["proxy"] = proxy
	return
}

func ProxyConfigFromForm(prefix string, data url.Values) (visitor echo.H, proxy echo.H, err error) {
	visitor = echo.H{}
	proxy = echo.H{}
	if len(prefix) > 0 {
		prefix += `.`
	}
	mapx := echo.NewMapx(data)
	extra := mapx.Get(`extra`)
	if extra == nil {
		return
	}
	for proxyName, m := range extra.Map {
		local := m.Value("local_port")
		remote := m.Value("remote_port")
		if strings.Contains(local, ",") || strings.Contains(remote, ",") {
			err = ParseRangePort(proxyName, m, extra.Map)
			if err != nil {
				return
			}
			delete(extra.Map, proxyName)
		}
	}
	for proxyName, m := range extra.Map {
		localC := config.LocalSvrConf{
			LocalIp:      m.Value(`local_ip`),
			LocalPort:    param.String(m.Value(`local_port`)).Int(),
			Plugin:       m.Value(`plugin`),
			PluginParams: map[string]string{},
		}
		baseC := config.BaseProxyConf{
			ProxyName:      prefix + proxyName,
			ProxyType:      m.Value(`protocol`),
			UseEncryption:  param.String(m.Value(`use_encryption`)).Bool(),
			UseCompression: param.String(m.Value(`use_compression`)).Bool(),
			Group:          m.Value(`group`),
			GroupKey:       m.Value(`group_key`),
			LocalSvrConf:   localC,
		}
		if pluginParams := m.Get(`plugin_params`); pluginParams != nil {
			for kk, vv := range pluginParams.Map {
				localC.PluginParams[kk] = vv.Value()
			}
		}
		var value interface{}
		switch baseC.ProxyType {
		case consts.TcpProxy:
			recv := &config.TcpProxyConf{
				BaseProxyConf: baseC,
			}
			recv.BindInfoConf.RemotePort = param.String(m.Value(`remote_port`)).Int()
			recv.LocalSvrConf = localC
			value = recv
		case consts.UdpProxy:
			recv := &config.UdpProxyConf{
				BaseProxyConf: baseC,
			}
			recv.BindInfoConf.RemotePort = param.String(m.Value(`remote_port`)).Int()
			value = recv
		case consts.HttpProxy:
			recv := &config.HttpProxyConf{
				BaseProxyConf: baseC,
			}
			recv.DomainConf.CustomDomains = strings.Split(m.Value(`custom_domains`), `,`)
			recv.DomainConf.SubDomain = m.Value(`subdomain`)
			recv.Locations = strings.Split(m.Value(`locations`), `,`)
			recv.HttpUser = m.Value(`http_user`)
			recv.HttpPwd = m.Value(`http_pwd`)
			recv.HostHeaderRewrite = m.Value(`host_header_rewrite`)
			recv.Headers = map[string]string{}
			hd := m.Get(`header`)
			if hd != nil {
				keys := hd.Values(`k`)
				vals := hd.Values(`v`)
				for i, k := range keys {
					if len(k) == 0 {
						continue
					}
					var v string
					if len(vals) > i {
						v = vals[i]
					}
					recv.Headers[k] = v
				}
			}
			value = recv
		case consts.HttpsProxy:
			recv := &config.HttpsProxyConf{
				BaseProxyConf: baseC,
			}
			customDomains := m.Value(`custom_domains`)
			if len(customDomains) > 0 {
				recv.DomainConf.CustomDomains = strings.Split(customDomains, `,`)
			}
			recv.DomainConf.SubDomain = m.Value(`subdomain`)
			value = recv
		case consts.StcpProxy:
			recv := &config.StcpProxyConf{
				BaseProxyConf: baseC,
			}
			recv.Role = m.Value(`role`)
			recv.Sk = m.Value(`sk`)
			value = recv
		case consts.XtcpProxy:
			recv := &config.XtcpProxyConf{
				BaseProxyConf: baseC,
			}
			recv.Role = m.Value(`role`)
			recv.Sk = m.Value(`sk`)
			value = recv
		default:
			log.Errorf(`[frp]"%s" is unsupported`, baseC.ProxyType)
			continue
		}
		if m.Value(`role`) == `visitor` {
			visitor[proxyName] = value
		} else {
			proxy[proxyName] = value
		}
	}
	return
}

func ParseRangePort(name string, m *echo.Mapx, mapx map[string]*echo.Mapx) (err error) {
	local := m.Value("local_port")
	remote := m.Value("remote_port")
	localPorts, errRet := util.ParseRangeNumbers(local)
	if errRet != nil {
		err = fmt.Errorf("[frp]Parse conf error: [%s] local_port invalid, %v", name, errRet)
		return
	}

	remotePorts, errRet := util.ParseRangeNumbers(remote)
	if errRet != nil {
		err = fmt.Errorf("[frp]ParseRangeNumbers error: [%s] remote_port invalid, %v", name, errRet)
		return
	}
	if len(localPorts) != len(remotePorts) {
		err = fmt.Errorf("[frp]ParseRangeNumbers error: [%s] local ports number should be same with remote ports number", name)
		return
	}
	if len(localPorts) == 0 {
		err = fmt.Errorf("[frp]ParseRangeNumbers error: [%s] local_port and remote_port is necessary", name)
		return
	}
	for i, port := range localPorts {
		subName := fmt.Sprintf("%s_%d", name, i)
		mCopy := m.Clone()
		mCopy.Add("local_port", []string{fmt.Sprintf("%d", port)})
		mCopy.Add("remote_port", []string{fmt.Sprintf("%d", remotePorts[i])})
		mapx[subName] = mCopy
	}
	return
}
