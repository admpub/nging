//go:build rpc
// +build rpc

package model

import (
	"context"
	"fmt"

	"github.com/admpub/nging/v4/application/library/rpc"
)

func (f *FrpServer) CallRPC(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error {
	if f.DashboardPort > 0 {
		address := fmt.Sprintf(`%s:%d`, f.DashboardAddr, f.DashboardPort)
		rpcClient := rpc.NewClient(address, f.DashboardPwd, nil)
		if args == nil {
			args = &rpc.Empty{}
		}
		if reply == nil {
			reply = &rpc.Empty{}
		}
		return rpcClient.Call(ctx, `ServerRPCService`, serviceMethod, args, reply)
	}
	return rpc.ErrRPCServerDisabled
}
