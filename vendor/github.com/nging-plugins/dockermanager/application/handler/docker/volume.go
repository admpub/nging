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
	"strings"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/request"
)

func VolumeIndex(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	cfg := volume.ListOptions{
		Filters: filters.Args{},
	}
	list, err := c.VolumeList(ctx, cfg)
	if err != nil {
		return err
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`docker/base/volume/index`, handler.Err(ctx, err))
}

func VolumeAdd(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.VolumeAdd)
		var result volume.Volume
		result, err = c.VolumeCreate(ctx, req.Options())
		if err != nil {
			goto END
		}
		ctx.Logger().Debugf(`VolumeCreate: %+v`, result)
		return ctx.Redirect(handler.URLFor(`/docker/base/volume/index`))
	}

END:
	ctx.Set(`activeURL`, `/docker/base/volume/index`)
	ctx.Set(`title`, ctx.T(`新建存储卷`))
	return ctx.Render(`docker/base/volume/add`, err)
}

func VolumeDetail(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	volumeID := ctx.Param(`id`)
	data, err := c.VolumeInspect(ctx, volumeID)
	if err != nil {
		return err
	}
	ctx.Set(`activeURL`, `/docker/base/volume/index`)
	ctx.Set(`title`, ctx.T(`存储卷信息`))
	ctx.Set(`detail`, data)
	return ctx.Render(`docker/base/volume/detail`, err)
}

func VolumePrune(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	args := filters.NewArgs()
	until := ctx.Formx(`until`).String()
	if len(until) > 0 {
		args.Add(`until`, until) // until = 24h (删除24小时前的镜像)
	}
	report, err := c.VolumesPrune(ctx, args)
	if err == nil {
		//ctx.Logger().Debugf(`VolumesPrune: %+v`, report)
		handler.SendOk(ctx, ctx.T(
			"操作成功。\n删除卷: %s\n收回空间: %s",
			strings.Join(report.VolumesDeleted, `, `),
			com.FormatBytes(report.SpaceReclaimed),
		))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/volume/index`))
}

func VolumeDelete(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	force := ctx.Formx(`force`).Bool()
	volumeID := ctx.Param(`id`)
	if volumeID == `0` {
		errs := common.NewErrors()
		for _, volumeID := range ctx.FormValues(`id[]`) {
			if len(volumeID) == 0 {
				continue
			}
			err = c.VolumeRemove(ctx, volumeID, force)
			if err != nil {
				errs.Add(err)
				continue
			}
		}
		err = errs.ToError()
	} else {
		err = c.VolumeRemove(ctx, volumeID, force)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/volume/index`))
}
