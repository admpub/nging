package frp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/admpub/confl"
	"github.com/admpub/events"
	_ "github.com/admpub/frp/assets/frpc/statik"
	"github.com/admpub/frp/client"
	"github.com/admpub/frp/pkg/config"
	frpLog "github.com/admpub/frp/pkg/util/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/frpmanager/application/dbschema"
)

func SetClientConfigFromDB(conf *dbschema.NgingFrpClient) *config.ClientCommonConf {
	c := config.GetDefaultClientConf()
	c.ServerAddr = conf.ServerAddr
	c.ServerPort = int(conf.ServerPort)
	c.Protocol = conf.Protocol
	c.User = conf.User
	c.Token = conf.Token

	// TODO:
	// c.AuthenticationMethod = `token`
	// c.AuthenticateHeartBeats = false
	// c.AuthenticateNewWorkConns = false
	// c.OidcClientID = ""
	// c.OidcClientSecret = ""
	// c.OidcAudience = ""
	// c.OidcTokenEndpointURL = ""
	// c.AssetsDir = ``
	// TODO:
	// c.TLSEnable = conf.TlsEnable
	// c.TLSCertFile = conf.TlsCertFile
	// c.TLSKeyFile = conf.TlsKeyFile
	// c.TLSTrustedCaFile = conf.TlsTrustedCaFile
	// c.TLSServerName = conf.TlsServerName

	c.DisableLogColor = true

	c.LogLevel = conf.LogLevel
	c.LogFile = conf.LogFile
	c.LogMaxDays = int64(conf.LogMaxDays)
	c.HTTPProxy = conf.HttpProxy
	c.LogWay = conf.LogWay
	if c.LogWay == `console` || len(c.LogFile) == 0 {
		c.LogFile = `console`
	} else {
		com.MkdirAll(filepath.Dir(c.LogFile), os.ModePerm)
	}
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
	if len(conf.Metas) > 0 {
		c.Metas = ParseMetas(conf.Metas)
	}
	c.UDPPacketSize = 1500
	return &c
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
			frpLog.Error(`[frp]parseProxyConfig: %v`, err.Error())
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
			frpLog.Error(`[frp]parseProxyConfig: %v`, err.Error())
			continue
		}
		visitorCfgs[prefix+key] = recv
	}
	return
}

func RecvProxyConfig(data map[string]interface{}) (recv config.ProxyConf) {
	proxyType, _ := data[`type`].(string)
	recv = config.DefaultProxyConf(proxyType)
	if recv == nil {
		frpLog.Error(`[frp]Unsupported Proxy Type: %v`, proxyType)
		return
	}
	b, err := json.Marshal(data)
	if err == nil {
		err = json.Unmarshal(b, recv)
	}
	if err != nil {
		frpLog.Error(`[frp]RecvProxyConfig: %v`, err.Error())
		return
	}
	return
}

func RecvVisitorConfig(data map[string]interface{}) (recv config.VisitorConf) {
	proxyType, _ := data[`type`].(string)
	recv = config.DefaultVisitorConf(proxyType)
	if recv == nil {
		frpLog.Error(`[frp]Unsupported Visitor Type: %v`, proxyType)
		return
	}
	b, err := json.Marshal(data)
	if err == nil {
		err = json.Unmarshal(b, recv)
	}
	if err != nil {
		frpLog.Error(`[frp]RecvVisitorConfig: %v`, err.Error())
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
		b, err := config.GetRenderedConfFromFile(filePath)
		if err != nil {
			return fmt.Errorf("load frpc config file error: %w", err)
		}
		r := NewClientConfig()
		err = json.Unmarshal(b, r)
		if err != nil {
			return fmt.Errorf("load frpc config file unmarshal error: %w", err)
		}
		c = SetClientConfigFromDB(r.NgingFrpClient)
		if len(r.Extra) > 0 {
			pxyCfgs, visitorCfgs = parseProxyConfig(c, r.Extra)
		}
		filePath = ``
	case `.yaml`:
		b, err := config.GetRenderedConfFromFile(filePath)
		if err != nil {
			return fmt.Errorf("load frpc config file error: %w", err)
		}
		r := NewClientConfig()
		_, err = confl.Decode(string(b), r)
		if err != nil {
			return fmt.Errorf("load frpc config file decode error: %w", err)
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

var clientService *client.Service

func StartClient(pxyCfgs map[string]config.ProxyConf, visitorCfgs map[string]config.VisitorConf,
	pidFile string, c *config.ClientCommonConf, configFileArg ...string) (err error) {
	once.Do(onceInit)
	var configFile string
	if len(configFileArg) > 0 {
		configFile = configFileArg[0]
	}
	frpLog.InitLog(c.LogWay, c.LogFile, c.LogLevel, c.LogMaxDays, c.DisableLogColor)
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
			frpLog.Error(err.Error())
			return err
		}
	}
	/*
		echo.Dump(c)
		echo.Dump(pxyCfgs)
		echo.Dump(visitorCfgs)
	*/

	hookData := echo.H{`clientConfig`: c}
	if err := echo.Fire(echo.NewEvent(`nging.plugins.frpmanager.client.start.before`, events.WithContext(hookData))); err != nil {
		return err
	}

	if clientService != nil {
		clientService.Close()
	}
	clientService, err = client.NewService(*c, pxyCfgs, visitorCfgs, configFile)
	if err != nil {
		return err
	}
	defer clientService.Close()

	hookData[`clientService`] = clientService
	if err := echo.Fire(echo.NewEvent(`nging.plugins.frpmanager.client.start.after`, events.WithContext(hookData))); err != nil {
		return err
	}

	if c.Protocol == "kcp" {
		kcpDoneCh = make(chan struct{})
		// Capture the exit signal if we use kcp.
		go handleSignal(clientService, kcpDoneCh)
	}

	if err = clientService.Run(); err != nil {
		return
	}

	if c.Protocol == "kcp" {
		<-kcpDoneCh
	}
	return
}

func handleSignal(svr *client.Service, kcpDoneCh chan struct{}) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.Close()
	time.Sleep(250 * time.Millisecond)
	close(kcpDoneCh)
	os.Exit(0)
}
