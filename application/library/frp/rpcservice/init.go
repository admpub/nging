package rpcservice

import (
	"github.com/admpub/frp/client"
	"github.com/admpub/frp/server"
	"github.com/admpub/nging/application/library/frp"
)

func init() {
	frp.RegisterClientRPCService(func(s *client.Service) interface{} {
		return NewClientRPCService(s)
	})
	frp.RegisterServerRPCService(func(s *server.Service) interface{} {
		return NewServerRPCService(s)
	})
}
