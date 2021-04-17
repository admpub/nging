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

package model

import (
	"context"
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/rpc"
	"github.com/admpub/nging/application/model/base"
)

type FrpServerAndGroup struct {
	*dbschema.NgingFrpServer
	Group   *dbschema.NgingFrpGroup `db:"-,relation=id:group_id|gtZero"`
	Running bool
}

func NewFrpServer(ctx echo.Context) *FrpServer {
	return &FrpServer{
		NgingFrpServer: &dbschema.NgingFrpServer{},
		Base:           base.New(ctx),
	}
}

type FrpServer struct {
	*dbschema.NgingFrpServer
	*base.Base
}

func (f *FrpServer) Exists(name string, excludeIds ...uint) (bool, error) {
	cond := db.Cond{`name`: name}
	if len(excludeIds) > 0 {
		cond[`id`] = db.NotEq(excludeIds[0])
	}
	return f.NgingFrpServer.Exists(nil, cond)
}

func (f *FrpServer) check() error {
	if len(f.Name) == 0 {
		return f.NewError(code.InvalidParameter, f.T(`名称不能为空`)).SetZone(`name`)
	}
	var exists bool
	var err error
	if f.Id > 0 {
		exists, err = f.Exists(f.Name, f.Id)
	} else {
		exists, err = f.Exists(f.Name)
	}
	if err != nil {
		return err
	}
	if exists {
		return f.NewError(code.DataAlreadyExists, f.T(`名称已经存在`)).SetZone(`name`)
	}
	return err
}

func (f *FrpServer) Add() (pk interface{}, err error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	return f.NgingFrpServer.Add()
}

func (f *FrpServer) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.NgingFrpServer.Edit(mw, args...)
}

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
