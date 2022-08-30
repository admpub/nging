package rpcservice

import (
	"fmt"

	"github.com/admpub/events"
	"github.com/admpub/frp/client"
	"github.com/admpub/frp/pkg/config"
	frpLog "github.com/admpub/frp/pkg/util/log"
	"github.com/admpub/frp/server"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/rpc"
)

var (
	ClientRPCServices = []func(*client.Service) interface{}{}
	ServerRPCServices = []func(*server.Service) interface{}{}
)

func RegisterClientRPCService(r func(*client.Service) interface{}) {
	ClientRPCServices = append(ClientRPCServices, r)
}

func RegisterServerRPCService(r func(*server.Service) interface{}) {
	ServerRPCServices = append(ServerRPCServices, r)
}

func init() {
	RegisterClientRPCService(func(s *client.Service) interface{} {
		return NewClientRPCService(s)
	})
	RegisterServerRPCService(func(s *server.Service) interface{} {
		return NewServerRPCService(s)
	})

	// - client -

	echo.OnCallback(`nging.plugins.frpmanager.client.start.before`, func(data events.Event) error {
		c := data.Context.Get("clientConfig").(*config.ClientCommonConf)
		port := c.AdminPort
		if port > 0 && len(ClientRPCServices) > 0 {
			c.AdminPort = 0
		}
		data.Context.Set("port", port)
		return nil
	})
	echo.OnCallback(`nging.plugins.frpmanager.client.start.after`, func(data events.Event) error {
		port := data.Context.Int("port")
		c := data.Context.Get("clientConfig").(*config.ClientCommonConf)
		if port > 0 && len(ClientRPCServices) > 0 {
			clientService := data.Context.Get("clientService").(*client.Service)
			address := fmt.Sprintf(`%s:%d`, c.AdminAddr, port)
			rpcServer := rpc.NewServer(address, c.AdminPwd, nil)
			defer rpcServer.Close()
			for _, gen := range ClientRPCServices {
				rpcServer.Register(gen(clientService))
			}
			frpLog.Info(`[frpc] rpc server started: %s`, address)
			go frpLog.Error(`[frpc] rpc server exited: %v`, rpcServer.Start())
		}
		return nil
	})

	// - server -

	echo.OnCallback(`nging.plugins.frpmanager.server.start.before`, func(data events.Event) error {
		c := data.Context.Get("serverConfig").(*config.ServerCommonConf)
		port := c.DashboardPort
		if port > 0 && len(ServerRPCServices) > 0 {
			c.DashboardPort = 0
		}
		data.Context.Set("port", port)
		return nil
	})
	echo.OnCallback(`nging.plugins.frpmanager.server.start.after`, func(data events.Event) error {
		port := data.Context.Int("port")
		c := data.Context.Get("serverConfig").(*config.ServerCommonConf)
		if port > 0 && len(ServerRPCServices) > 0 {
			serverService := data.Context.Get("serverService").(*server.Service)
			address := fmt.Sprintf(`%s:%d`, c.DashboardAddr, port)
			rpcServer := rpc.NewServer(address, c.DashboardPwd, nil)
			defer rpcServer.Close()
			for _, gen := range ServerRPCServices {
				rpcServer.Register(gen(serverService))
			}
			frpLog.Info(`[frpc] rpc server started: %s`, address)
			go frpLog.Error(`[frpc] rpc server exited: %v`, rpcServer.Start())
		}
		return nil
	})
}
