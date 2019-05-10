package server

import (
	"github.com/admpub/frp/assets"
	"github.com/admpub/frp/g"
	"github.com/admpub/frp/utils/version"
	"github.com/webx-top/echo"
)

func APIServerInfo(c echo.Context) error {
	cfg := &g.GlbServerCfg.ServerCommonConf
	serverStats := StatsGetServer()
	res := ServerInfoResp{
		Version:           version.Full(),
		BindPort:          cfg.BindPort,
		BindUdpPort:       cfg.BindUdpPort,
		VhostHttpPort:     cfg.VhostHttpPort,
		VhostHttpsPort:    cfg.VhostHttpsPort,
		KcpBindPort:       cfg.KcpBindPort,
		AuthTimeout:       cfg.AuthTimeout,
		SubdomainHost:     cfg.SubDomainHost,
		MaxPoolCount:      cfg.MaxPoolCount,
		MaxPortsPerClient: cfg.MaxPortsPerClient,
		HeartBeatTimeout:  cfg.HeartBeatTimeout,

		TotalTrafficIn:  serverStats.TotalTrafficIn,
		TotalTrafficOut: serverStats.TotalTrafficOut,
		CurConns:        serverStats.CurConns,
		ClientCounts:    serverStats.ClientCounts,
		ProxyTypeCounts: serverStats.ProxyTypeCounts,
	}
	return c.JSON(res)
}

func APIProxyByType(c echo.Context) error {
	var res GetProxyInfoResp
	proxyType := c.Param(`type`)
	res.Proxies = getProxyStatsByType(proxyType)
	return c.JSON(res)
}

func APIProxyByTypeAndName(c echo.Context) error {
	proxyType := c.Param(`type`)
	name := c.Param(`name`)
	res := getProxyStatsByTypeAndName(proxyType, name)
	return c.JSON(res)
}

func APIProxyTraffic(c echo.Context) error {
	var res GetProxyTrafficResp
	res.Name = c.Param(`name`)
	proxyTrafficInfo := StatsGetProxyTraffic(res.Name)
	if proxyTrafficInfo == nil {
		res.Code = 1
		res.Msg = "no proxy info found"
	} else {
		res.TrafficIn = proxyTrafficInfo.TrafficIn
		res.TrafficOut = proxyTrafficInfo.TrafficOut
	}

	return c.JSON(res)
}

// RegisterTo 为echo框架创建路由
func RegisterTo(router echo.RouteRegister) {
	// api
	router.Get("/api/serverinfo", APIServerInfo)
	router.Get("/api/proxy/:type", APIProxyByType)
	router.Get("/api/proxy/:type/:name", APIProxyByTypeAndName)
	router.Get("/api/traffic/:name", APIProxyTraffic)

	// view
	router.Get("/", func(c echo.Context) error {
		return c.Redirect("./static/")
	})
	cfg := &g.GlbServerCfg.ServerCommonConf
	//cfg.AssetsDir = `/Users/hank/go/src/github.com/admpub/frp/assets/static`
	err := assets.Load(cfg.AssetsDir)
	if err != nil {
		panic(err)
	}
	router.Get("/static*", func(c echo.Context) error {
		file := c.Param(`*`)
		return c.File(file, assets.FileSystem)
	})
}
