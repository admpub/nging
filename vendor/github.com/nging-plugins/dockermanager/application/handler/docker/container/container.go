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

package container

import (
	"strings"

	"github.com/admpub/log"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/nging-plugins/dockermanager/application/request"
)

func Index(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	opts := types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(), // https://docs.docker.com/engine/reference/commandline/ps/#filtering
	}
	imageID := ctx.Formx(`imageId`).String()
	if len(imageID) > 0 {
		opts.Filters.Add(`ancestor`, imageID)
	}
	list, err := c.ContainerList(ctx, opts)
	if err != nil {
		return err
	}
	if ctx.Form(`op`) == `ajaxList` {
		if ctx.Form(`type`) == `selectpage` {
			return utils.AjaxListSelectpage(ctx, list, func(v types.Container) echo.H {
				if len(v.Names) == 0 {
					return nil
				}
				return echo.H{`id`: v.ID, `name`: strings.TrimPrefix(v.Names[0], `/`)}
			})
		}
		return utils.AjaxListTypeahead(ctx, list, func(v types.Container) string {
			if len(v.Names) == 0 {
				return ``
			}
			return strings.TrimPrefix(v.Names[0], `/`)
		})
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`docker/base/container/index`, handler.Err(ctx, err))
}

func Add(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	var containerID string
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.ContainerAdd)
		var result container.CreateResponse
		result, err = c.ContainerCreate(ctx, req.Container(), req.Host(), req.Network(), req.Platform(), req.Name)
		if err != nil {
			goto END
		}
		ctx.Logger().Debugf(`ContainerCreate: %+v`, result)
		containerID = result.ID
		if req.StartNow {
			opts := types.ContainerStartOptions{}
			err = c.ContainerStart(ctx, containerID, opts)
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功，容器已经启动`))
			} else {
				handler.SendFail(ctx, ctx.T(`容器创建成功，但启动失败：%s`, err.Error()))
			}
		} else {
			handler.SendOk(ctx, ctx.T(`容器创建成功，可点击启动按钮来启动它`))
		}
		return ctx.Redirect(handler.URLFor(`/docker/base/container/index`))
	}
	containerID = ctx.Formx(`copyId`).String()
	if len(containerID) > 0 {
		var data types.ContainerJSON
		data, err = c.ContainerInspect(ctx, containerID)
		if err == nil {
			setContainerDataToFormAdd(ctx, c, data)
		}
	}

END:
	ctx.Set(`activeURL`, `/docker/base/container/index`)
	ctx.Set(`title`, ctx.T(`新建容器`))
	return ctx.Render(`docker/base/container/add`, err)
}

func Detail(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	data, err := c.ContainerInspect(ctx, containerID)
	if err != nil {
		return err
	}
	ctx.Set(`activeURL`, `/docker/base/container/index`)
	ctx.Set(`title`, ctx.T(`容器详情`))
	ctx.Set(`detail`, data)
	return ctx.Render(`docker/base/container/detail`, err)
}

func Edit(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	data, err := c.ContainerInspect(ctx, containerID)
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.ContainerEdit)
		if strings.TrimPrefix(data.ContainerJSONBase.Name, `/`) != strings.TrimPrefix(req.Name, `/`) {
			err = c.ContainerRename(ctx, containerID, req.Name)
			if err != nil {
				goto END
			}
		}
		updateOpts := container.UpdateConfig{}
		updateOpts.Resources = req.Resources()
		updateOpts.RestartPolicy.Name = req.RestartPolicy
		updateOpts.RestartPolicy.MaximumRetryCount = req.RestartMaxRetryCount
		var updateResult container.ContainerUpdateOKBody
		updateResult, err = c.ContainerUpdate(ctx, containerID, updateOpts)
		if err != nil {
			goto END
		}
		for _, warningMsg := range updateResult.Warnings {
			log.Warn(warningMsg)
		}
		return ctx.Redirect(handler.URLFor(`/docker/base/container/index`))
	}
	setContainerDataToFormEdit(ctx, data)

END:
	ctx.Set(`activeURL`, `/docker/base/container/index`)
	ctx.Set(`title`, ctx.T(`更新容器`))
	ctx.Set(`detail`, data)

	return ctx.Render(`docker/base/container/edit`, err)
}

func Prune(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	args := filters.NewArgs()
	until := ctx.Formx(`until`).String()
	if len(until) > 0 {
		args.Add(`until`, until) // until = 24h (删除已停止超过24小时的容器)
	}
	report, err := c.ContainersPrune(ctx, args)
	if err == nil {
		handler.SendOk(ctx, ctx.T(
			"操作成功。\n删除容器: %s\n收回空间: %s",
			strings.Join(report.ContainersDeleted, `, `),
			com.FormatBytes(report.SpaceReclaimed),
		))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/container/index`))
}

func Delete(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	opts := types.ContainerRemoveOptions{
		Force:         ctx.Formx(`force`).Bool(),
		RemoveLinks:   ctx.Formx(`removeLinks`).Bool(),
		RemoveVolumes: ctx.Formx(`removeVolumes`).Bool(),
	}
	all := ctx.Formx(`all`).Bool()
	if all {
		opts.RemoveLinks = true
		opts.RemoveVolumes = true
	}
	containerID := ctx.Param(`id`)
	if containerID == `0` {
		errs := common.NewErrors()
		for _, containerID := range ctx.FormValues(`id[]`) {
			if len(containerID) == 0 {
				continue
			}
			err = c.ContainerRemove(ctx, containerID, opts)
			if err != nil {
				errs.Add(err)
				continue
			}
		}
		err = errs.ToError()
	} else {
		err = c.ContainerRemove(ctx, containerID, opts)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/container/index`))
}
