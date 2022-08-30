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

	"github.com/admpub/frp/pkg/config"
	"github.com/admpub/frp/pkg/consts"
	"github.com/admpub/frp/pkg/util/util"
	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/nging-plugins/frpmanager/application/dbschema"
)

// Table2Config 数据库表中的数据转为frp配置文件数据格式
func Table2Config(cc *dbschema.NgingFrpClient) (hash echo.H, err error) {
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

func ParseMetas(metas string) map[string]string {
	r := map[string]string{}
	metas = strings.TrimSpace(metas)
	for _, meta := range strings.Split(metas, "\n") {
		meta = strings.TrimSpace(meta)
		if len(meta) == 0 {
			continue
		}
		vs := strings.SplitN(meta, "=", 2)
		if len(vs) != 2 {
			continue
		}
		vs[0] = strings.TrimSpace(vs[0])
		vs[1] = strings.TrimSpace(vs[1])
		if len(vs[0]) == 0 || len(vs[1]) == 0 {
			continue
		}
		r[vs[0]] = vs[1]
	}
	return r
}

// ProxyConfigFromForm 将表单数据转为frpc代理参数
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
			LocalIP:      m.Value(`local_ip`),
			LocalPort:    param.String(m.Value(`local_port`)).Int(),
			Plugin:       m.Value(`plugin`),
			PluginParams: map[string]string{},
		}
		if len(localC.LocalIP) == 0 {
			localC.LocalIP = `127.0.0.1`
		}
		//TODO:
		healthCheckC := config.HealthCheckConf{
			HealthCheckType:      ``, // tcp/http 留空表示禁用
			HealthCheckTimeoutS:  3,
			HealthCheckMaxFailed: 1,
			HealthCheckIntervalS: 10,
			HealthCheckURL:       ``, // for http
			HealthCheckAddr:      ``, // for tcp
		}
		//TODO:
		bandwidthC, _ := config.NewBandwidthQuantity(m.Value(`bandwidth_quantity`)) // unit: KB/MB
		baseC := config.BaseProxyConf{
			ProxyName:            prefix + proxyName,
			ProxyType:            m.Value(`protocol`),
			UseEncryption:        param.String(m.Value(`use_encryption`)).Bool(),
			UseCompression:       param.String(m.Value(`use_compression`)).Bool(),
			Group:                m.Value(`group`),
			GroupKey:             m.Value(`group_key`),
			ProxyProtocolVersion: ``,
			BandwidthLimit:       bandwidthC,
			Metas:                map[string]string{},
			LocalSvrConf:         localC,
			HealthCheckConf:      healthCheckC,
		}
		if pluginParams := m.Get(`plugin_params`); pluginParams != nil {
			for kk, vv := range pluginParams.Map {
				baseC.LocalSvrConf.PluginParams[kk] = vv.Value()
			}
		}
		metas := m.Value(`metas`)
		if len(metas) > 0 {
			baseC.Metas = ParseMetas(metas)
		}
		var value interface{}
		switch baseC.ProxyType {
		case consts.TCPProxy:
			recv := &config.TCPProxyConf{
				BaseProxyConf: baseC,
			}
			recv.RemotePort = param.String(m.Value(`remote_port`)).Int()
			value = recv
		case consts.TCPMuxProxy: //TODO:
			recv := &config.TCPMuxProxyConf{
				BaseProxyConf: baseC,
			}
			recv.DomainConf.CustomDomains = strings.Split(m.Value(`custom_domains`), `,`)
			recv.DomainConf.SubDomain = m.Value(`subdomain`)
			recv.Multiplexer = m.Value(`multiplexer`)
			value = recv
		case consts.UDPProxy:
			recv := &config.UDPProxyConf{
				BaseProxyConf: baseC,
			}
			recv.RemotePort = param.String(m.Value(`remote_port`)).Int()
			value = recv
		case consts.HTTPProxy:
			recv := &config.HTTPProxyConf{
				BaseProxyConf: baseC,
			}
			recv.DomainConf.CustomDomains = strings.Split(m.Value(`custom_domains`), `,`)
			recv.DomainConf.SubDomain = m.Value(`subdomain`)
			recv.Locations = strings.Split(m.Value(`locations`), `,`)
			recv.HTTPUser = m.Value(`http_user`)
			recv.HTTPPwd = m.Value(`http_pwd`)
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
		case consts.HTTPSProxy:
			recv := &config.HTTPSProxyConf{
				BaseProxyConf: baseC,
			}
			customDomains := m.Value(`custom_domains`)
			if len(customDomains) > 0 {
				recv.DomainConf.CustomDomains = strings.Split(customDomains, `,`)
			}
			recv.DomainConf.SubDomain = m.Value(`subdomain`)
			value = recv
		case consts.STCPProxy:
			recv := &config.STCPProxyConf{
				BaseProxyConf: baseC,
			}
			recv.Role = m.Value(`role`)
			recv.Sk = m.Value(`sk`)
			value = recv
		case consts.XTCPProxy:
			recv := &config.XTCPProxyConf{
				BaseProxyConf: baseC,
			}
			recv.Role = m.Value(`role`)
			recv.Sk = m.Value(`sk`)
			value = recv
		case consts.SUDPProxy:
			recv := &config.SUDPProxyConf{
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
