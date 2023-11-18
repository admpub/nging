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
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/request"
)

func NodeIndex(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	list, err := c.NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		return detectSwarmError(ctx, err)
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`docker/swarm/node/index`, handler.Err(ctx, err))
}

func NodeAdd(ctx echo.Context) error {
	ctx.Set(`activeURL`, `/docker/swarm/node/index`)
	ctx.Set(`title`, ctx.T(`添加节点`))
	return SwarmJoin(ctx)
}

func NodeEdit(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	nodeID := ctx.Param(`id`)
	data, _, err := c.NodeInspectWithRaw(ctx, nodeID)
	if err != nil {
		return err
	}
	var req *request.SwarmNodeEdit
	if ctx.IsPost() {
		swarmVersion := swarm.Version{
			Index: data.Version.Index,
		}
		req = echo.GetValidated(ctx).(*request.SwarmNodeEdit)
		err = c.NodeUpdate(ctx, nodeID, swarmVersion, req.NodeSpec)
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, ctx.T(`更新成功`))
		return ctx.Redirect(handler.URLFor(`/docker/swarm/node/index`))
	}
	req = &request.SwarmNodeEdit{NodeSpec: data.Spec}
	echo.StructToForm(ctx, req, ``, echo.LowerCaseFirstLetter, param.StringerMapStart().AddFunc(`labels`, func(v interface{}) string {
		return com.JoinKVRows(v)
	}))

END:
	ctx.Set(`activeURL`, `/docker/swarm/node/index`)
	ctx.Set(`title`, ctx.T(`更新节点`))
	ctx.Set(`detail`, data)
	return ctx.Render(`docker/swarm/node/edit`, err)
}

func NodeDetail(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	nodeID := ctx.Param(`id`)
	data, _, err := c.NodeInspectWithRaw(ctx, nodeID)
	if err != nil {
		return detectSwarmError(ctx, err)
	}
	ctx.Set(`activeURL`, `/docker/swarm/node/index`)
	ctx.Set(`title`, ctx.T(`节点信息`))
	ctx.Set(`detail`, data)
	return ctx.Render(`docker/swarm/node/detail`, handler.Err(ctx, err))
}

func NodeDelete(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	nodeID := ctx.Param(`id`)
	force := ctx.Formx(`force`).Bool()
	opts := types.NodeRemoveOptions{Force: force}
	if nodeID == `0` {
		errs := common.NewErrors()
		for _, nodeID := range ctx.FormValues(`id[]`) {
			if len(nodeID) == 0 {
				continue
			}
			err = c.NodeRemove(ctx, nodeID, opts)
			if err != nil {
				errs.Add(err)
				continue
			}
			ctx.Logger().Debugf(`NodeRemove: %v`, nodeID)
		}
		err = errs.ToError()
	} else {
		err = c.NodeRemove(ctx, nodeID, opts)
		if err == nil {
			ctx.Logger().Debugf(`NodeRemove: %v`, nodeID)
		}
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/swarm/node/index`))
}
