package rpcservice

import (
	"context"

	"github.com/admpub/frp/client"
	"github.com/webx-top/com"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/rpc"
)

func NewClientRPCService(s *client.Service) *ClientRPCService {
	return &ClientRPCService{s: s}
}

type ClientRPCService struct {
	s *client.Service
}

func (s *ClientRPCService) service() *client.Service {
	return s.s
}

func (s *ClientRPCService) Status(ctx context.Context, args *rpc.Empty, reply *client.StatusResp) error {
	eCtx := common.NewMockContext()
	eCtx.Response().KeepBody(true)
	err := s.service().APIStatus(eCtx)
	if err != nil {
		return err
	}
	res := eCtx.Response().Body()
	err = com.JSONDecode(res, reply)
	return err
}
