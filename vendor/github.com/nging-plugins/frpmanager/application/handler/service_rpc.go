//go:build rpc
// +build rpc

package handler

import (
	"context"

	"github.com/admpub/frp/client"
	"github.com/admpub/log"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/rpc"
	"github.com/admpub/nging/v4/application/model"

	_ "github.com/nging-plugins/frpmanager/application/library/frp/rpcservice"
)

func clientTestRPC(ctx echo.Context) {
	m := model.NewFrpClient(ctx)
	err := m.Get(nil, db.Cond{`disabled`: `N`})
	if err != nil {
		panic(err)
	}
	err = m.CallRPC(context.Background(), `Status`, &rpc.Empty{}, &client.StatusResp{})
	if err != nil {
		log.Error(err)
	}
}

func serverTestRPC(ctx echo.Context) {
	m := model.NewFrpServer(ctx)
	err := m.Get(nil, db.Cond{`disabled`: `N`})
	if err != nil {
		panic(err)
	}
	res := &echo.H{}
	err = m.CallRPC(context.Background(), `ServerInfo`, &rpc.Empty{}, res)
	if err != nil {
		log.Error(err)
	}
	echo.Dump(res)
}
