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
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/confl"
	"github.com/admpub/events"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/frpmanager/application/dbschema"

	_ "github.com/admpub/frp/assets/frps/statik"
	"github.com/admpub/frp/pkg/config"
	plugin "github.com/admpub/frp/pkg/plugin/server"
	frpLog "github.com/admpub/frp/pkg/util/log"
	"github.com/admpub/frp/pkg/util/util"
	"github.com/admpub/frp/server"
)

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

	// TODO:
	//c.EnablePrometheus = false
	//c.AssetsDir = ``

	c.LogFile = conf.LogFile
	c.LogWay = conf.LogWay
	if c.LogWay == `console` || len(c.LogFile) == 0 {
		c.LogFile = `console`
	} else {
		com.MkdirAll(filepath.Dir(c.LogFile), os.ModePerm)
	}
	c.LogLevel = conf.LogLevel
	c.LogMaxDays = int64(conf.LogMaxDays)
	c.DisableLogColor = true
	c.DetailedErrorsToClient = true
	c.Custom404Page = ``
	c.Token = conf.Token

	// TODO:
	// c.AuthenticationMethod = `token`
	// c.AuthenticateHeartBeats = false
	// c.AuthenticateNewWorkConns = false
	// c.OidcIssuer = ""
	// c.OidcAudience = ""
	// c.OidcSkipExpiryCheck = false
	// c.OidcSkipIssuerCheck = false

	c.SubDomainHost = conf.SubdomainHost
	c.MaxPoolCount = int64(conf.MaxPoolCount)
	if c.MaxPoolCount < 1 {
		c.MaxPoolCount = 5
	}
	c.MaxPortsPerClient = conf.MaxPortsPerClient
	c.TCPMux = conf.TcpMux == `Y`

	// e.g. 1000-2000,2001,2002,3000-4000
	ports, _ := util.ParseRangeNumbers(conf.AllowPorts)
	for _, port := range ports {
		c.AllowPorts[int(port)] = struct{}{}
	}

	// TODO:
	// c.TLSOnly = conf.TlsOnly
	// c.TLSCertFile = conf.TlsCertFile
	// c.TLSKeyFile = conf.TlsKeyFile
	// c.TLSTrustedCaFile = conf.TlsTrustedCaFile

	c.HeartbeatTimeout = int64(conf.HeartBeatTimeout)
	if c.HeartbeatTimeout < 1 {
		c.HeartbeatTimeout = 90
	}
	c.UserConnTimeout = int64(conf.UserConnTimeout)
	if c.UserConnTimeout < 1 {
		c.UserConnTimeout = 10
	}
	if c.HTTPPlugins == nil {
		c.HTTPPlugins = map[string]plugin.HTTPPluginOptions{}
	}
	c.UDPPacketSize = 1500
	configExtra := NewServerConfigExtra()
	configExtra.Parse(conf.Extra)
	conf.Plugins = strings.TrimSpace(conf.Plugins)
	if len(conf.Plugins) > 0 {
		plugins := strings.Split(conf.Plugins, ",")
		for _, name := range plugins {
			if options, ok := configExtra.PluginOptions[name]; ok {
				c.HTTPPlugins[name] = options
			}
		}
	}
	return &c
}

func StartServerByConfigFile(filePath string, pidFile string) error {
	ext := filepath.Ext(filePath)
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	switch strings.ToLower(ext) {
	case `.json`:
		r := dbschema.NewNgingFrpServer(nil)
		err = json.Unmarshal(b, r)
		if err != nil {
			return err
		}
		c := SetServerConfigFromDB(r)
		return StartServer(pidFile, c)
	case `.yaml`:
		r := dbschema.NewNgingFrpServer(nil)
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

var serverService *server.Service

func StartServer(pidFile string, c *config.ServerCommonConf) error {
	once.Do(onceInit)
	err := c.Validate()
	if err != nil {
		return err
	}
	frpLog.InitLog(c.LogWay, c.LogFile, c.LogLevel, c.LogMaxDays, c.DisableLogColor)
	if c.HTTPPlugins != nil {
		for name, options := range c.HTTPPlugins {
			frpLog.Info(`[frps][plugin] register %s, API URL: %s`, name, options.Addr+options.Path)
		}
	}
	if len(pidFile) > 0 {
		err := com.WritePidFile(pidFile)
		if err != nil {
			frpLog.Error(err.Error())
			return err
		}
	}

	hookData := echo.H{`serverConfig`: c}
	if err := echo.Fire(echo.NewEvent(`nging.plugins.frpmanager.server.start.before`, events.WithContext(hookData))); err != nil {
		return err
	}

	serverService, err = server.NewService(*c)
	if err != nil {
		return err
	}
	// defer svr.Close() 无此方法

	hookData[`serverService`] = serverService
	if err := echo.Fire(echo.NewEvent(`nging.plugins.frpmanager.server.start.after`, events.WithContext(hookData))); err != nil {
		return err
	}

	frpLog.Info("Start frps success")
	serverService.Run()
	return err
}
