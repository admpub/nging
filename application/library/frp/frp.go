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
	"fmt"
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
	"github.com/admpub/frp/pkg/config"
	"github.com/admpub/frp/pkg/consts"
	"github.com/admpub/frp/pkg/util/log"
	"github.com/admpub/frp/pkg/util/util"
	"github.com/admpub/frp/server"
	"github.com/admpub/nging/application/dbschema"
)

func SetClientConfigFromDB(conf *dbschema.NgingFrpClient) *config.ClientCommonConf {
	c := config.GetDefaultClientConf()
	c.ServerAddr = conf.ServerAddr
	c.ServerPort = int(conf.ServerPort)
	c.User = conf.User
	c.Protocol = conf.Protocol
	c.Token = conf.Token
	c.LogLevel = conf.LogLevel
	if conf.LogWay == `console` {
		conf.LogFile = `console`
	}
	c.LogFile = conf.LogFile
	c.LogMaxDays = int64(conf.LogMaxDays)
	c.HTTPProxy = conf.HttpProxy
	c.LogWay = conf.LogWay
	c.AdminAddr = conf.AdminAddr
	c.AdminPort = int(conf.AdminPort)
	c.AdminUser = conf.AdminUser
	c.AdminPwd = conf.AdminPwd
	c.PoolCount = int(conf.PoolCount)
	c.TCPMux = conf.TcpMux == `Y`
	c.DNSServer = conf.DnsServer
	c.LoginFailExit = conf.LoginFailExit == `Y`
	c.HeartbeatInterval = conf.HeartbeatInterval
	c.HeartbeatTimeout = conf.HeartbeatTimeout
	conf.Start = strings.TrimSpace(conf.Start)
	if len(conf.Start) > 0 {
		for _, name := range strings.Split(conf.Start, `,`) {
			c.Start = append(c.Start, strings.TrimSpace(name))
		}
	}
	return &c
}

func SetServerConfigFromDB(conf *dbschema.NgingFrpServer) *config.ServerCommonConf {
	c := config.GetDefaultServerConf()
	c.BindAddr = conf.Addr
	c.BindPort = int(conf.Port)
	c.BindUDPPort = int(conf.UdpPort)
	c.KCPBindPort = int(conf.KcpPort)
	c.ProxyBindAddr = conf.ProxyAddr
	c.VhostHTTPPort = int(conf.VhostHttpPort)
	c.VhostHTTPTimeout = int64(conf.VhostHttpTimeout)
	if c.VhostHTTPTimeout < 1 {
		c.VhostHTTPTimeout = 60
	}
	c.VhostHTTPSPort = int(conf.VhostHttpsPort)

	c.DashboardAddr = conf.DashboardAddr
	c.DashboardPort = int(conf.DashboardPort)
	c.DashboardUser = conf.DashboardUser
	c.DashboardPwd = conf.DashboardPwd
	if conf.LogWay == `console` {
		conf.LogFile = `console`
	}
	c.LogFile = conf.LogFile
	c.LogWay = conf.LogWay
	c.LogLevel = conf.LogLevel
	c.LogMaxDays = int64(conf.LogMaxDays)
	c.Token = conf.Token
	c.SubDomainHost = conf.SubdomainHost
	c.MaxPortsPerClient = int64(conf.MaxPortsPerClient)
	c.TCPMux = conf.TcpMux == `Y`

	// e.g. 1000-2000,2001,2002,3000-4000
	ports, _ := util.ParseRangeNumbers(conf.AllowPorts)
	for _, port := range ports {
		c.AllowPorts[int(port)] = struct{}{}
	}
	return &c
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
		c := SetServerConfigFromDB(r)
		return StartServer(pidFile, c)
	case `.yaml`:
		r := &dbschema.NgingFrpServer{}
		err = confl.Unmarshal(b, r)
		if err != nil {
			return err
		}
		c := SetServerConfigFromDB(r)
		return StartServer(pidFile, c)
	default:
		content := string(b)
		return StartServerByConfig(content, pidFile)
	}
}

func StartServerByConfig(configContent string, pidFile string) error {
	cfg, err := config.UnmarshalServerConfFromIni(configContent)
	if err != nil {
		return err
	}
	return StartServer(pidFile, &cfg)
}

func StartServer(pidFile string, c *config.ServerCommonConf) error {
	err := c.Validate()
	if err != nil {
		return err
	}
	log.InitLog(c.LogWay,
		c.LogFile,
		c.LogLevel,
		c.LogMaxDays, true)
	if len(pidFile) > 0 {
		err := com.WritePidFile(pidFile)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	svr, err := server.NewService(*c)
	if err != nil {
		return err
	}
	log.Info("Start frps success")
	svr.Run()
	return err
}

func parseProxyConfig(c *config.ClientCommonConf, extra echo.H) (
	pxyCfgs map[string]config.ProxyConf,
	visitorCfgs map[string]config.VisitorConf,
) {
	pxyCfgs = map[string]config.ProxyConf{}
	visitorCfgs = map[string]config.VisitorConf{}
	proxyConfs := NewProxyConfig()
	proxyConfs.Visitor, _ = extra[`visitor`].(map[string]interface{})
	proxyConfs.Proxy, _ = extra[`proxy`].(map[string]interface{})
	prefix := c.User
	if len(prefix) > 0 {
		prefix += `.`
	}
	startProxy := c.Start
	startAll := true
	if len(startProxy) > 0 {
		startAll = false
	}
	for key, cfg := range proxyConfs.Proxy {
		shouldStart := com.InSlice(key, startProxy)
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
		shouldStart := com.InSlice(key, startProxy)
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
	case consts.TCPProxy:
		recv = &config.TCPProxyConf{}
	case consts.UDPProxy:
		recv = &config.UDPProxyConf{}
	case consts.HTTPProxy:
		recv = &config.HTTPProxyConf{}
	case consts.HTTPSProxy:
		recv = &config.HTTPSProxyConf{}
	case consts.STCPProxy:
		recv = &config.STCPProxyConf{}
	case consts.XTCPProxy:
		recv = &config.XTCPProxyConf{}
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
	case consts.STCPProxy:
		recv = &config.STCPVisitorConf{}
	case consts.XTCPProxy:
		recv = &config.XTCPVisitorConf{}
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
		c           *config.ClientCommonConf
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
		c = SetClientConfigFromDB(r.NgingFrpClient)
		if len(r.Extra) > 0 {
			pxyCfgs, visitorCfgs = parseProxyConfig(c, r.Extra)
		}
		filePath = ``
	case `.yaml`:
		r := NewClientConfig()
		_, err := confl.DecodeFile(filePath, r)
		if err != nil {
			return err
		}
		c = SetClientConfigFromDB(r.NgingFrpClient)
		if len(r.Extra) > 0 {
			pxyCfgs, visitorCfgs = parseProxyConfig(c, r.Extra)
		}
		filePath = ``
	default:
		content, err := config.GetRenderedConfFromFile(filePath)
		if err != nil {
			return fmt.Errorf("load frpc config file error: %w", err)
		}
		var _c config.ClientCommonConf
		_c, err = config.UnmarshalClientConfFromIni(content)
		if err != nil {
			return fmt.Errorf("load frpc common section error: %w", err)
		}
		c = &_c
		pxyCfgs, visitorCfgs, err = config.LoadAllProxyConfsFromIni(c.User, content, c.Start)
		if err != nil {
			return err
		}
	}
	return StartClient(pxyCfgs, visitorCfgs, pidFile, c, filePath)
}

func StartClientByConfig(configContent string, pidFile string) error {
	c, err := config.UnmarshalClientConfFromIni(configContent)
	if err != nil {
		return fmt.Errorf("load frpc common section error: %w", err)
	}
	pxyCfgs, visitorCfgs, err := config.LoadAllProxyConfsFromIni(c.User, configContent, c.Start)
	if err != nil {
		return err
	}
	return StartClient(pxyCfgs, visitorCfgs, pidFile, &c)
}

func StartClient(pxyCfgs map[string]config.ProxyConf, visitorCfgs map[string]config.VisitorConf,
	pidFile string, c *config.ClientCommonConf, configFileArg ...string) (err error) {
	var configFile string
	if len(configFileArg) > 0 {
		configFile = configFileArg[0]
	}
	log.InitLog(c.LogWay, c.LogFile, c.LogLevel, c.LogMaxDays, true)
	if len(c.DNSServer) > 0 {
		s := c.DNSServer
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
		echo.Dump(c)
		echo.Dump(pxyCfgs)
		echo.Dump(visitorCfgs)
	*/
	svr, err := client.NewService(*c, pxyCfgs, visitorCfgs, configFile)
	if err != nil {
		return err
	}

	err = svr.Run()

	// Capture the exit signal if we use kcp.
	if c.Protocol == "kcp" {
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
