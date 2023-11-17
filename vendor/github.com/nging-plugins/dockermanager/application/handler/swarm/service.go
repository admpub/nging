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

package swarm

import (
	"bufio"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/websocket"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/nging-plugins/dockermanager/application/request"
)

func ServiceIndex(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	args := filters.NewArgs()
	// id / label / mode / name
	id := ctx.Form(`id`)
	if len(id) > 0 {
		args.Add(`id`, id)
	}
	label := ctx.Form(`label`)
	if len(label) > 0 {
		args.Add(`label`, label)
	}
	mode := ctx.Form(`mode`)
	if len(mode) > 0 {
		args.Add(`mode`, mode)
	}
	name := ctx.Form(`name`)
	if len(name) > 0 {
		args.Add(`name`, name)
	}
	list, err := c.ServiceList(ctx, types.ServiceListOptions{Filters: args, Status: true})
	if err != nil {
		return detectSwarmError(ctx, err)
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`docker/swarm/service/index`, handler.Err(ctx, err))
}

func ServiceAdd(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.SwarmServiceAdd)
		opts := req.ServiceCreateOptions
		spec := req.ServiceSpec
		//echo.Dump(spec)
		var result types.ServiceCreateResponse
		result, err = c.ServiceCreate(ctx, spec, opts)
		if err != nil {
			goto END
		}
		for _, warningMsg := range result.Warnings {
			log.Warn(warningMsg)
		}
		handler.SendOk(ctx, ctx.T(`创建成功`))
		return ctx.Redirect(handler.URLFor(`/docker/swarm/service/index`))
	}

END:
	ctx.Set(`activeURL`, `/docker/swarm/service/index`)
	ctx.Set(`title`, ctx.T(`新建服务`))
	ctx.Set(`isEdit`, false)
	return ctx.Render(`docker/swarm/service/edit`, err)
}

func ServiceEdit(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	serviceID := ctx.Param(`id`)
	data, _, err := c.ServiceInspectWithRaw(ctx, serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return err
	}
	var req *request.SwarmServiceEdit
	if ctx.IsPost() {
		req = echo.GetValidated(ctx).(*request.SwarmServiceEdit)
		updateOpts := req.ServiceUpdateOptions
		swarmVersion := swarm.Version{
			Index: data.Version.Index,
		}
		spec := req.ServiceSpec
		var updateResult types.ServiceUpdateResponse
		updateResult, err = c.ServiceUpdate(ctx, serviceID, swarmVersion, spec, updateOpts)
		if err != nil {
			goto END
		}
		for _, warningMsg := range updateResult.Warnings {
			log.Warn(warningMsg)
		}
		handler.SendOk(ctx, ctx.T(`更新成功`))
		return ctx.Redirect(handler.URLFor(`/docker/swarm/service/index`))
	}
	req = &request.SwarmServiceEdit{ServiceSpec: data.Spec}
	echo.StructToForm(ctx, req, ``, nil)
	if req.Mode.Global != nil {
		ctx.Request().Form().Set(`mode`, `global`)
	} else if req.Mode.GlobalJob != nil {
		ctx.Request().Form().Set(`mode`, `globalJob`)
	} else if req.Mode.Replicated != nil {
		ctx.Request().Form().Set(`mode`, `replicated`)
	} else if req.Mode.ReplicatedJob != nil {
		ctx.Request().Form().Set(`mode`, `replicatedJob`)
	}
	//echo.Dump(ctx.Forms())

END:
	ctx.Set(`activeURL`, `/docker/swarm/service/index`)
	ctx.Set(`title`, ctx.T(`更新服务`))
	ctx.Set(`detail`, data)
	ctx.Set(`isEdit`, true)
	return ctx.Render(`docker/swarm/service/edit`, err)
}

func ServiceRollback(ctx echo.Context) error {
	data := ctx.Data()
	c, err := dockerclient.Client()
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	serviceID := ctx.Param(`id`)
	svc, _, err := c.ServiceInspectWithRaw(ctx, serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	req := &request.SwarmServiceEdit{}
	updateOpts := req.ServiceUpdateOptions
	updateOpts.Rollback = `previous`
	swarmVersion := swarm.Version{
		Index: svc.Version.Index,
	}
	spec := req.ServiceSpec
	var updateResult types.ServiceUpdateResponse
	updateResult, err = c.ServiceUpdate(ctx, serviceID, swarmVersion, spec, updateOpts)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	for _, warningMsg := range updateResult.Warnings {
		log.Warn(warningMsg)
	}
	return ctx.JSON(data.SetInfo(ctx.T(`操作成功`, 1)))
}

func ServiceDetail(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	serviceID := ctx.Param(`id`)
	data, _, err := c.ServiceInspectWithRaw(ctx, serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return err
	}
	ctx.Set(`activeURL`, `/docker/swarm/service/index`)
	ctx.Set(`title`, ctx.T(`服务信息`))
	ctx.Set(`detail`, data)
	return ctx.Render(`docker/swarm/service/detail`, err)
}

func ServiceDelete(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	serviceID := ctx.Param(`id`)
	if serviceID == `0` {
		errs := common.NewErrors()
		for _, serviceID := range ctx.FormValues(`id[]`) {
			if len(serviceID) == 0 {
				continue
			}
			err = c.ServiceRemove(ctx, serviceID)
			if err != nil {
				errs.Add(err)
				continue
			}
			ctx.Logger().Debugf(`ServiceRemove: %v`, serviceID)
		}
		err = errs.ToError()
	} else {
		err = c.ServiceRemove(ctx, serviceID)
		if err == nil {
			ctx.Logger().Debugf(`ServiceRemove: %v`, serviceID)
		}
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/swarm/service/index`))
}

func ServiceLogs(conn *websocket.Conn, ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	serviceID := ctx.Param(`id`)
	opts := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}
	reader, err := c.ServiceLogs(ctx, serviceID, opts)
	if err != nil {
		return err
	}
	defer reader.Close()
	buf := bufio.NewReader(reader)
	for {
		message, err := buf.ReadString('\n')
		if err != nil {
			return err
		}
		message = strings.TrimSuffix(message, "\n")
		message = strings.TrimSuffix(message, "\r")
		message = utils.TrimHeader(message)
		if err = conn.WriteMessage(websocket.BinaryMessage, []byte(message+"\r\n")); err != nil {
			return err
		}
	}
}
