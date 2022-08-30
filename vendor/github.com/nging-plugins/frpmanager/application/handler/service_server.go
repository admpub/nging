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

package handler

import (
	"github.com/admpub/log"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/writer"
	"github.com/nging-plugins/frpmanager/application/library/cmder"
)

func ServerRestart(ctx echo.Context) error {
	cm, err := cmder.GetServer()
	if err != nil {
		return err
	}
	data := ctx.Data()
	if err := cm.Stop(); err != nil {
		data.SetError(err)
		return ctx.JSON(data)
	}
	if err := cm.RebuildConfigFile(`frpserver`); err != nil {
		data.SetError(err)
		return ctx.JSON(data)
	}
	buf := writer.NewShadow()
	wOut := writer.NewOut(buf)
	wErr := writer.NewErr(buf)
	if err := cm.Start(wOut, wErr); err != nil {
		data.SetError(err)
		return ctx.JSON(data)
	}
	msg := ctx.T(`已经重启FRP服务端`)
	log.Info(msg)
	startResult := cleanStartResult.ReplaceAllString(buf.String(), `$1 `)
	if len(startResult) > 0 {
		msg += ":\n" + startResult
	}
	data.SetInfo(msg, 1)
	return ctx.JSON(data)
}

func ServerStop(ctx echo.Context) error {
	cm, err := cmder.GetServer()
	if err != nil {
		return err
	}
	data := ctx.Data()
	if err := cm.Stop(); err != nil {
		data.SetError(err)
		return ctx.JSON(data)
	}
	if err := cm.StopHistory(); err != nil {
		data.SetError(err)
		return ctx.JSON(data)
	}
	msg := ctx.T(`已经关闭FRP服务端`)
	log.Info(msg)
	data.SetInfo(msg, 1)
	return ctx.JSON(data)
}
