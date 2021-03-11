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
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/confl"
	_ "github.com/admpub/frp/assets/frpc/statik"
	_ "github.com/admpub/frp/assets/frps/statik"
	"github.com/admpub/frp/client"
	"github.com/admpub/frp/g"
	"github.com/admpub/frp/models/config"
	"github.com/admpub/frp/models/consts"
	"github.com/admpub/frp/server"
	"github.com/admpub/frp/utils/log"
	"github.com/admpub/frp/utils/util"
	"github.com/admpub/ini"
	"github.com/admpub/nging/application/dbschema"
)

func SetClientConfigFromDB(conf *dbschema.NgingFrpClient) *g.ClientCfg {
	g.GlbClientCfg.ServerAddr = conf.ServerAddr
	g.GlbClientCfg.ServerPort = int(conf.ServerPort)
	g.GlbClientCfg.User = conf.User
	g.GlbClientCfg.Protocol = conf.Protocol
	g.GlbClientCfg.Token = conf.Token
	g.GlbClientCfg.LogLevel = conf.LogLevel
	if conf.LogWay == `console` {
		conf.LogFile = `console`
	}
	g.GlbClientCfg.LogFile = conf.LogFile
	g.GlbClientCfg.LogMaxDays = int64(conf.LogMaxDays)
	g.GlbClientCfg.HttpProxy = conf.HttpProxy
	g.GlbClientCfg.LogWay = conf.LogWay
	g.GlbClientCfg.AdminAddr = conf.AdminAddr
	g.GlbClientCfg.AdminPort = int(conf.AdminPort)
	g.GlbClientCfg.AdminUser = conf.AdminUser
	g.GlbClientCfg.AdminPwd = conf.AdminPwd
	g.GlbClientCfg.PoolCount = int(conf.PoolCount)
	g.GlbClientCfg.TcpMux = conf.TcpMux == `Y`
	g.GlbClientCfg.DnsServer = conf.DnsServer
	g.GlbClientCfg.LoginFailExit = conf.LoginFailExit == `Y`
	g.GlbClientCfg.HeartBeatInterval = conf.HeartbeatInterval
	g.GlbClientCfg.HeartBeatTimeout = conf.HeartbeatTimeout
	conf.Start = strings.TrimSpace(conf.Start)
	if len(conf.Start) > 0 {
		for _, name := range strings.Split(conf.Start, `,`) {
			g.GlbClientCfg.Start[strings.TrimSpace(name)] = struct{}{}
		}
	}
	return g.GlbClientCfg
}

func SetServerConfigFromDB(conf *dbschema.NgingFrpServer) *g.ServerCfg {
	g.GlbServerCfg.BindAddr = conf.Addr
	g.GlbServerCfg.BindPort = int(conf.Port)
	g.GlbServerCfg.BindUdpPort = int(conf.UdpPort)
	g.GlbServerCfg.KcpBindPort = int(conf.KcpPort)
	g.GlbServerCfg.ProxyBindAddr = conf.ProxyAddr
	g.GlbServerCfg.VhostHttpPort = int(conf.VhostHttpPort)
	g.GlbServerCfg.VhostHttpTimeout = int64(conf.VhostHttpTimeout)
	if g.GlbServerCfg.VhostHttpTimeout < 1 {
		g.GlbServerCfg.VhostHttpTimeout = 60
	}
	g.GlbServerCfg.VhostHttpsPort = int(conf.VhostHttpsPort)

	g.GlbServerCfg.DashboardAddr = conf.DashboardAddr
	g.GlbServerCfg.DashboardPort = int(conf.DashboardPort)
	g.GlbServerCfg.DashboardUser = conf.DashboardUser
	g.GlbServerCfg.DashboardPwd = conf.DashboardPwd
	if conf.LogWay == `console` {
		conf.LogFile = `console`
	}
	g.GlbServerCfg.LogFile = conf.LogFile
	g.GlbServerCfg.LogWay = conf.LogWay
	g.GlbServerCfg.LogLevel = conf.LogLevel
	g.GlbServerCfg.LogMaxDays = int64(conf.LogMaxDays)
	g.GlbServerCfg.Token = conf.Token
	g.GlbServerCfg.SubDomainHost = conf.SubdomainHost
	g.GlbServerCfg.MaxPortsPerClient = int64(conf.MaxPortsPerClient)
	g.GlbServerCfg.TcpMux = conf.TcpMux == `Y`

	// e.g. 1000-2000,2001,2002,3000-4000
	ports, _ := util.ParseRangeNumbers(conf.AllowPorts)
	for _, port := range ports {
		g.GlbServerCfg.AllowPorts[int(port)] = struct{}{}
	}
	return g.GlbServerCfg
}

func StartServerByConfigFile(filePath string, pidFile string) error {
	ext := filepath.Ext(filePath)
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	switch strings.ToLower(ext) {
	case `.json`:
		r := &dbschema.NgingFrpServer{}
		err = json.Unmarshal(b, r)
		if err != nil {
			return err
		}
		SetServerConfigFromDB(r)
		return StartServer(pidFile)
	case `.yaml`:
		r := &dbschema.NgingFrpServer{}
		err = confl.Unmarshal(b, r)
		if err != nil {
			return err
		}
		SetServerConfigFromDB(r)
		return StartServer(pidFile)
	default:
		content := string(b)
		return StartServerByConfig(content, pidFile)
	}
}

func StartServerByConfig(configContent string, pidFile string) error {
	cfg, err := config.UnmarshalServerConfFromIni(&g.GlbServerCfg.ServerCommonConf, configContent)
	if err != nil {
		return err
	}
	g.GlbServerCfg.ServerCommonConf = *cfg
	return StartServer(pidFile)
}

func StartServer(pidFile string) error {
	err := g.GlbServerCfg.ServerCommonConf.Check()
	if err != nil {
		return err
	}
	config.InitServerCfg(&g.GlbServerCfg.ServerCommonConf)

	log.InitLog(g.GlbServerCfg.LogWay,
		g.GlbServerCfg.LogFile,
		g.GlbServerCfg.LogLevel,
		g.GlbServerCfg.LogMaxDays)
	if len(pidFile) > 0 {
		err := com.WritePidFile(pidFile)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	svr, err := server.NewService()
	if err != nil {
		return err
	}
	log.Info("Start frps success")
	server.ServerService = svr
	svr.Run()
	return err
}

func parseProxyConfig(extra echo.H) (pxyCfgs map[string]config.ProxyConf, visitorCfgs map[string]config.VisitorConf) {
	pxyCfgs = map[string]config.ProxyConf{}
	visitorCfgs = map[string]config.VisitorConf{}
	proxyConfs := NewProxyConfig()
	proxyConfs.Visitor, _ = extra[`visitor`].(map[string]interface{})
	proxyConfs.Proxy, _ = extra[`proxy`].(map[string]interface{})
	prefix := g.GlbClientCfg.User
	if len(prefix) > 0 {
		prefix += `.`
	}
	startProxy := g.GlbClientCfg.Start
	startAll := true
	if len(startProxy) > 0 {
		startAll = false
	}
	for key, cfg := range proxyConfs.Proxy {
		_, shouldStart := startProxy[key]
		if !startAll && !shouldStart {
			continue
		}
		_cfg, _ok := cfg.(map[string]interface{})
		if !_ok {
			continue
		}
		recv := RecvProxyConfig(_cfg)
		if recv == nil {
			continue
		}
		err := recv.CheckForCli()
		if err != nil {
			log.Error(`[frp]parseProxyConfig:`, err)
			continue
		}
		pxyCfgs[prefix+key] = recv
	}
	for key, cfg := range proxyConfs.Visitor {
		_, shouldStart := startProxy[key]
		if !startAll && !shouldStart {
			continue
		}
		_cfg, _ok := cfg.(map[string]interface{})
		if !_ok {
			continue
		}
		recv := RecvVisitorConfig(_cfg)
		if recv == nil {
			continue
		}
		err := recv.Check()
		if err != nil {
			log.Error(`[frp]parseProxyConfig:`, err)
			continue
		}
		visitorCfgs[prefix+key] = recv
	}
	return
}

func RecvProxyConfig(data map[string]interface{}) (recv config.ProxyConf) {
	proxyType, _ := data[`proxy_type`].(string)
	switch proxyType {
	case consts.TcpProxy:
		recv = &config.TcpProxyConf{}
	case consts.UdpProxy:
		recv = &config.UdpProxyConf{}
	case consts.HttpProxy:
		recv = &config.HttpProxyConf{}
	case consts.HttpsProxy:
		recv = &config.HttpsProxyConf{}
	case consts.StcpProxy:
		recv = &config.StcpProxyConf{}
	case consts.XtcpProxy:
		recv = &config.XtcpProxyConf{}
	default:
		log.Error(`[frp]Unsupported Proxy Type:`, proxyType)
		return
	}
	b, err := json.Marshal(data)
	if err == nil {
		err = json.Unmarshal(b, recv)
	}
	if err != nil {
		log.Error(`[frp]RecvProxyConfig:`, err)
		return
	}
	return
}

func RecvVisitorConfig(data map[string]interface{}) (recv config.VisitorConf) {
	proxyType, _ := data[`proxy_type`].(string)
	switch proxyType {
	case consts.StcpProxy:
		recv = &config.StcpVisitorConf{}
	case consts.XtcpProxy:
		recv = &config.XtcpVisitorConf{}
	default:
		log.Error(`[frp]Unsupported Visitor Type:`, proxyType)
		return
	}
	b, err := json.Marshal(data)
	if err == nil {
		err = json.Unmarshal(b, recv)
	}
	if err != nil {
		log.Error(`[frp]RecvVisitorConfig:`, err)
		return
	}
	return
}

func StartClientByConfigFile(filePath string, pidFile string) error {
	var (
		pxyCfgs     map[string]config.ProxyConf
		visitorCfgs map[string]config.VisitorConf
	)
	ext := filepath.Ext(filePath)
	switch strings.ToLower(ext) {
	case `.json`:
		b, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		r := NewClientConfig()
		err = json.Unmarshal(b, r)
		if err != nil {
			return err
		}
		SetClientConfigFromDB(r.NgingFrpClient)
		if len(r.Extra) > 0 {
			pxyCfgs, visitorCfgs = parseProxyConfig(r.Extra)
		}
	case `.yaml`:
		r := NewClientConfig()
		_, err := confl.DecodeFile(filePath, r)
		if err != nil {
			return err
		}
		SetClientConfigFromDB(r.NgingFrpClient)
		if len(r.Extra) > 0 {
			pxyCfgs, visitorCfgs = parseProxyConfig(r.Extra)
		}
	default:
		conf, err := ini.Load(filePath)
		if err != nil {
			return err
		}
		pxyCfgs, visitorCfgs, err = config.LoadAllConfFromIni(g.GlbClientCfg.User, conf, g.GlbClientCfg.Start)
		if err != nil {
			return err
		}
	}
	return StartClient(pxyCfgs, visitorCfgs, pidFile)
}

func StartClientByConfig(configContent string, pidFile string) error {
	conf, err := ini.LoadContent(configContent)
	if err != nil {
		return err
	}
	pxyCfgs, visitorCfgs, err := config.LoadAllConfFromIni(g.GlbClientCfg.User, conf, g.GlbClientCfg.Start)
	if err != nil {
		return err
	}
	return StartClient(pxyCfgs, visitorCfgs, pidFile)
}

func StartClient(pxyCfgs map[string]config.ProxyConf, visitorCfgs map[string]config.VisitorConf, pidFile string) (err error) {
	log.InitLog(g.GlbClientCfg.LogWay, g.GlbClientCfg.LogFile, g.GlbClientCfg.LogLevel, g.GlbClientCfg.LogMaxDays)
	if len(g.GlbClientCfg.DnsServer) > 0 {
		s := g.GlbClientCfg.DnsServer
		if !strings.Contains(s, ":") {
			s += ":53"
		}
		// Change default dns server for frpc
		net.DefaultResolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", s)
			},
		}
	}
	if len(pidFile) > 0 {
		err := com.WritePidFile(pidFile)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	/*
		echo.Dump(g.GlbClientCfg)
		echo.Dump(pxyCfgs)
		echo.Dump(visitorCfgs)
	*/
	svr, err := client.NewService(pxyCfgs, visitorCfgs)
	if err != nil {
		return err
	}

	err = svr.Run()

	// Capture the exit signal if we use kcp.
	if g.GlbClientCfg.Protocol == "kcp" {
		var kcpDoneCh = make(chan struct{})
		go handleSignal(svr, kcpDoneCh)
		<-kcpDoneCh
	}
	return
}

func handleSignal(svr *client.Service, kcpDoneCh chan struct{}) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.Close()
	time.Sleep(250 * time.Millisecond)
	close(kcpDoneCh)
	os.Exit(0)
}
