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

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/nging-plugins/dockermanager/application/request"
)

func ConfigIndex(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	list, err := c.ConfigList(ctx, types.ConfigListOptions{})
	if err != nil {
		return detectSwarmError(ctx, err)
	}
	if ctx.Form(`op`) == `ajaxList` {
		if ctx.Form(`type`) == `selectpage` {
			return utils.AjaxListSelectpage(ctx, list, func(v swarm.Config) echo.H {
				return echo.H{`id`: v.ID, `name`: v.Spec.Name}
			})
		}
		return utils.AjaxListTypeahead(ctx, list, func(v swarm.Config) string {
			return v.Spec.Name
		})
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`docker/swarm/config/index`, handler.Err(ctx, err))
}

func ConfigAdd(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.SwarmConfigEdit)
		_, err = c.ConfigCreate(ctx, req.ConfigSpec)
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, ctx.T(`创建成功`))
		return ctx.Redirect(handler.URLFor(`/docker/swarm/config/index`))
	}

END:
	ctx.Set(`activeURL`, `/docker/swarm/config/index`)
	ctx.Set(`title`, ctx.T(`新建配置`))
	ctx.Set(`isEdit`, false)
	return ctx.Render(`docker/swarm/config/edit`, err)
}

func ConfigEdit(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	configID := ctx.Param(`id`)
	data, _, err := c.ConfigInspectWithRaw(ctx, configID)
	if err != nil {
		return err
	}
	var req *request.SwarmConfigEdit
	if ctx.IsPost() {
		swarmVersion := swarm.Version{
			Index: data.Version.Index,
		}
		req = echo.GetValidated(ctx).(*request.SwarmConfigEdit)
		spec := req.ConfigSpec
		if len(spec.Name) == 0 {
			spec.Name = data.Spec.Name
		}
		err = c.ConfigUpdate(ctx, configID, swarmVersion, spec)
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, ctx.T(`更新成功`))
		return ctx.Redirect(handler.URLFor(`/docker/swarm/config/index`))
	}
	req = &request.SwarmConfigEdit{ConfigSpec: data.Spec}
	req.Content = com.Bytes2str(req.Data)
	echo.StructToForm(ctx, req, ``, nil)

END:
	ctx.Set(`activeURL`, `/docker/swarm/config/index`)
	ctx.Set(`title`, ctx.T(`更新配置`))
	ctx.Set(`detail`, data)
	ctx.Set(`isEdit`, true)
	return ctx.Render(`docker/swarm/config/edit`, err)
}

func ConfigDetail(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	configID := ctx.Param(`id`)
	data, _, err := c.ConfigInspectWithRaw(ctx, configID)
	if err != nil {
		return err
	}
	ctx.Set(`activeURL`, `/docker/swarm/config/index`)
	ctx.Set(`title`, ctx.T(`配置信息`))
	ctx.Set(`detail`, data)
	return ctx.Render(`docker/swarm/config/detail`, err)
}

func ConfigDelete(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	configID := ctx.Param(`id`)
	if configID == `0` {
		errs := common.NewErrors()
		for _, configID := range ctx.FormValues(`id[]`) {
			if len(configID) == 0 {
				continue
			}
			err = c.ConfigRemove(ctx, configID)
			if err != nil {
				errs.Add(err)
				continue
			}
			ctx.Logger().Debugf(`ConfigRemove: %v`, configID)
		}
		err = errs.ToError()
	} else {
		err = c.ConfigRemove(ctx, configID)
		if err == nil {
			ctx.Logger().Debugf(`ConfigRemove: %v`, configID)
		}
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/swarm/config/index`))
}
