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

package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/nging-plugins/dockermanager/application/request"
)

func NetworkIndex(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	list, err := c.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return err
	}
	if ctx.Form(`op`) == `ajaxList` {
		return utils.AjaxListTypeahead(ctx, list, func(v types.NetworkResource) string {
			return v.Name
		})
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`docker/base/network/index`, handler.Err(ctx, err))
}

func NetworkAdd(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.NetworkAdd)
		var result types.NetworkCreateResponse
		result, err = c.NetworkCreate(ctx, req.Name, req.NetworkCreate)
		if err != nil {
			goto END
		}
		ctx.Logger().Debugf(`NetworkCreate: %+v`, result)
		return ctx.Redirect(handler.URLFor(`/docker/base/network/index`))
	}

END:
	ctx.Set(`activeURL`, `/docker/base/network/index`)
	ctx.Set(`title`, ctx.T(`新建网络`))
	return ctx.Render(`docker/base/network/add`, err)
}

func NetworkConnect(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.NetworkConnect)
		err = c.NetworkConnect(ctx, req.NetworkID, req.ContainerID, &req.EndpointSettings)
		if err != nil {
			goto END
		}
		return ctx.Redirect(handler.URLFor(`/docker/base/network/index`))
	}

END:
	ctx.Set(`activeURL`, `/docker/base/network/index`)
	ctx.Set(`title`, ctx.T(`连接网络`))
	return ctx.Render(`docker/base/network/edit`, err)
}

func NetworkDisconnect(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.NetworkDisconnect)
		err = c.NetworkDisconnect(ctx, req.NetworkID, req.ContainerID, req.Force)
		if err != nil {
			goto END
		}
		return ctx.Redirect(handler.URLFor(`/docker/base/network/index`))
	}

END:
	ctx.Set(`activeURL`, `/docker/base/network/index`)
	ctx.Set(`title`, ctx.T(`取消网络连接`))
	return ctx.Render(`docker/base/network/edit`, err)
}

func NetworkPrune(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	args := filters.NewArgs()
	until := ctx.Formx(`until`).String()
	if len(until) > 0 {
		args.Add(`until`, until) // until = 24h (删除24小时前的镜像)
	}
	var result types.NetworksPruneReport
	result, err = c.NetworksPrune(ctx, args)
	if err == nil {
		ctx.Logger().Debugf(`NetworksPrune: %+v`, result)
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}
	return ctx.Redirect(handler.URLFor(`/docker/base/network/index`))
}

func NetworkDetail(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	networkID := ctx.Param(`id`)
	opts := types.NetworkInspectOptions{
		Scope:   ctx.Formx(`scope`).String(),
		Verbose: ctx.Formx(`verbose`).Bool(),
	}
	data, err := c.NetworkInspect(ctx, networkID, opts)
	if err != nil {
		return err
	}
	ctx.Set(`activeURL`, `/docker/base/network/index`)
	ctx.Set(`title`, ctx.T(`网络信息`))
	ctx.Set(`detail`, data)
	return ctx.Render(`docker/base/network/detail`, err)
}

func NetworkDelete(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	networkID := ctx.Param(`id`)
	if networkID == `0` {
		errs := common.NewErrors()
		for _, networkID := range ctx.FormValues(`id[]`) {
			if len(networkID) == 0 {
				continue
			}
			err = c.NetworkRemove(ctx, networkID)
			if err != nil {
				errs.Add(err)
				continue
			}
		}
		err = errs.ToError()
	} else {
		err = c.NetworkRemove(ctx, networkID)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/network/index`))
}
