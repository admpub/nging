package rpcservice

import (
	"context"

	"github.com/admpub/frp/server"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/rpc"
)

func NewServerRPCService(s *server.Service) *ServerRPCService {
	return &ServerRPCService{s: s}
}

type ServerRPCService struct {
	s *server.Service
}

func (s *ServerRPCService) service() *server.Service {
	return s.s
}

func (s *ServerRPCService) ServerInfo(ctx context.Context, args *rpc.Empty, reply *echo.H) error {
	eCtx := common.NewMockContext()
	eCtx.Response().KeepBody(true)
	err := s.service().APIServerInfo(eCtx)
	if err != nil {
		return err
	}
	res := eCtx.Response().Body()
	err = com.JSONDecode(res, reply)
	return err
}

func (s *ServerRPCService) ProxyByType(ctx context.Context, args *rpc.Empty, reply *echo.H) error {
	eCtx := common.NewMockContext()
	eCtx.Response().KeepBody(true)
	err := s.service().APIProxyByType(eCtx)
	if err != nil {
		return err
	}
	res := eCtx.Response().Body()
	err = com.JSONDecode(res, reply)
	return err
}

func (s *ServerRPCService) ProxyByTypeAndName(ctx context.Context, args *rpc.Empty, reply *echo.H) error {
	eCtx := common.NewMockContext()
	eCtx.Response().KeepBody(true)
	err := s.service().APIProxyByTypeAndName(eCtx)
	if err != nil {
		return err
	}
	res := eCtx.Response().Body()
	err = com.JSONDecode(res, reply)
	return err
}

func (s *ServerRPCService) ProxyTraffic(ctx context.Context, args *rpc.Empty, reply *echo.H) error {
	eCtx := common.NewMockContext()
	eCtx.Response().KeepBody(true)
	err := s.service().APIProxyTraffic(eCtx)
	if err != nil {
		return err
	}
	res := eCtx.Response().Body()
	err = com.JSONDecode(res, reply)
	return err
}
