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

func SecretIndex(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	list, err := c.SecretList(ctx, types.SecretListOptions{})
	if err != nil {
		return detectSwarmError(ctx, err)
	}
	if ctx.Form(`op`) == `ajaxList` {
		if ctx.Form(`type`) == `selectpage` {
			return utils.AjaxListSelectpage(ctx, list, func(v swarm.Secret) echo.H {
				return echo.H{`id`: v.ID, `name`: v.Spec.Name}
			})
		}
		return utils.AjaxListTypeahead(ctx, list, func(v swarm.Secret) string {
			return v.Spec.Name
		})
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`docker/swarm/secret/index`, handler.Err(ctx, err))
}

func SecretAdd(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.SwarmSecretEdit)
		_, err = c.SecretCreate(ctx, req.SecretSpec)
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, ctx.T(`创建成功`))
		return ctx.Redirect(handler.URLFor(`/docker/swarm/secret/index`))
	}

END:
	ctx.Set(`activeURL`, `/docker/swarm/secret/index`)
	ctx.Set(`title`, ctx.T(`新建服务`))
	ctx.Set(`isEdit`, false)
	return ctx.Render(`docker/swarm/secret/edit`, err)
}

func SecretEdit(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	secretID := ctx.Param(`id`)
	data, _, err := c.SecretInspectWithRaw(ctx, secretID)
	if err != nil {
		return err
	}
	var req *request.SwarmSecretEdit
	if ctx.IsPost() {
		swarmVersion := swarm.Version{
			Index: data.Version.Index,
		}
		req = echo.GetValidated(ctx).(*request.SwarmSecretEdit)
		spec := req.SecretSpec
		if len(spec.Name) == 0 {
			spec.Name = data.Spec.Name
		}
		err = c.SecretUpdate(ctx, secretID, swarmVersion, spec)
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, ctx.T(`更新成功`))
		return ctx.Redirect(handler.URLFor(`/docker/swarm/secret/index`))
	}
	req = &request.SwarmSecretEdit{SecretSpec: data.Spec}
	req.Content = com.Bytes2str(req.Data)
	echo.StructToForm(ctx, req, ``, nil)
	//echo.Dump(ctx.Forms())

END:
	ctx.Set(`activeURL`, `/docker/swarm/secret/index`)
	ctx.Set(`title`, ctx.T(`更新服务`))
	ctx.Set(`detail`, data)
	ctx.Set(`isEdit`, true)
	return ctx.Render(`docker/swarm/secret/edit`, err)
}

func SecretDetail(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	secretID := ctx.Param(`id`)
	data, _, err := c.SecretInspectWithRaw(ctx, secretID)
	if err != nil {
		return err
	}
	ctx.Set(`activeURL`, `/docker/swarm/secret/index`)
	ctx.Set(`title`, ctx.T(`密钥信息`))
	ctx.Set(`detail`, data)
	return ctx.Render(`docker/swarm/secret/detail`, err)
}

func SecretDelete(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	secretID := ctx.Param(`id`)
	if secretID == `0` {
		errs := common.NewErrors()
		for _, secretID := range ctx.FormValues(`id[]`) {
			if len(secretID) == 0 {
				continue
			}
			err = c.SecretRemove(ctx, secretID)
			if err != nil {
				errs.Add(err)
				continue
			}
			ctx.Logger().Debugf(`SecretRemove: %v`, secretID)
		}
		err = errs.ToError()
	} else {
		err = c.SecretRemove(ctx, secretID)
		if err == nil {
			ctx.Logger().Debugf(`SecretRemove: %v`, secretID)
		}
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/swarm/secret/index`))
}
