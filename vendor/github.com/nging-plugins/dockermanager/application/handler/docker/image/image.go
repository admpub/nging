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

package image

import (
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

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
	opts := types.ImageListOptions{
		All:     true,
		Filters: filters.NewArgs(),
		//SharedSize:     true,
		//ContainerCount: true,
	}

	if searchValue := ctx.Form(`searchValue`); len(searchValue) > 0 {
		opts.Filters.Add(`reference`, searchValue)
	} else if name := ctx.Form(`name`, ctx.Form(`q`)); len(name) > 0 {
		// a
		opts.Filters.Add(`reference`, `*`+name+`*`)

		// a/b
		opts.Filters.Add(`reference`, `*/*`+name+`*`)
		opts.Filters.Add(`reference`, `*`+name+`*/*`)

		// a/b/c
		opts.Filters.Add(`reference`, `*/*/*`+name+`*`)
		opts.Filters.Add(`reference`, `*/*`+name+`*/*`)
		opts.Filters.Add(`reference`, `*`+name+`*/*/*`)
	} else if label := ctx.Form(`label`); len(label) > 0 {
		opts.Filters.Add(`label`, label)
	}

	list, err := c.ImageList(ctx, opts)
	if err != nil {
		return err
	}
	if ctx.Form(`op`) == `ajaxList` {
		if ctx.Form(`type`) == `selectpage` {
			return utils.AjaxListSelectpage(ctx, list, func(v types.ImageSummary) echo.H {
				if len(v.RepoTags) == 0 {
					return nil
				}
				//id := v.ID
				id := v.RepoTags[0]
				return echo.H{
					`id`:   id,
					`name`: com.HTMLEncode(strings.Join(v.RepoTags, ` / `)),
				}
			})
		}
		return utils.AjaxListTypeahead(ctx, list, func(v types.ImageSummary) string {
			if len(v.RepoTags) == 0 {
				return ``
			}
			return v.RepoTags[0]
		})
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`docker/base/image/index`, handler.Err(ctx, err))
}

func Add(ctx echo.Context) error {
	user := handler.User(ctx)
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.ImageAdd)
		err = dockerclient.PullImage(ctx, user, req.Ref, c, &types.ImagePullOptions{
			All:          req.All,
			RegistryAuth: req.RegistryAuth,
			Platform:     req.Platform,
		})
		if err != nil {
			return err
		}
		data := ctx.Data()
		data.SetInfo(ctx.T(`启动成功`), code.Success.Int())
		return ctx.JSON(data)
	}
	ctx.Set(`activeURL`, `/docker/base/image/index`)
	ctx.Set(`title`, ctx.T(`拉取镜像`))
	return ctx.Render(`docker/base/image/add`, err)
}

func Detail(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	imageID := ctx.Param(`id`)
	data, _, err := c.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		return err
	}
	historyList, err := c.ImageHistory(ctx, imageID)
	if err != nil {
		return err
	}
	ctx.Set(`activeURL`, `/docker/base/image/index`)
	ctx.Set(`title`, ctx.T(`镜像信息`))
	ctx.Set(`detail`, data)
	reversedHistoryList := make([]image.HistoryResponseItem, 0, len(historyList))
	for i := len(historyList) - 1; i >= 0; i-- {
		historyList[i].CreatedBy = strings.TrimPrefix(historyList[i].CreatedBy, `/bin/sh -c #(nop) `)
		historyList[i].CreatedBy = strings.TrimPrefix(historyList[i].CreatedBy, ` `)
		reversedHistoryList = append(reversedHistoryList, historyList[i])
	}
	ctx.Set(`historyList`, reversedHistoryList)
	return ctx.Render(`docker/base/image/detail`, err)
}

func Edit(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	req := echo.GetValidated(ctx).(*request.ImageTag)
	err = c.ImageTag(ctx, req.Source, req.Target)
	data := ctx.Data()
	if err != nil {
		data.SetError(err)
	}
	return ctx.JSON(data)
}

func Pull(ctx echo.Context) error {
	user := handler.User(ctx)
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	ref := ctx.Formx(`ref`).String()
	if len(ref) == 0 {
		return ctx.NewError(code.InvalidParameter, `ref值不正确`).SetZone(`ref`)
	}
	err = dockerclient.PullImage(ctx, user, ref, c, nil)
	if err != nil {
		return err
	}
	data := ctx.Data()
	data.SetInfo(ctx.T(`启动成功`), code.Success.Int())
	return ctx.JSON(data)
}

func Prune(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	args := filters.NewArgs()
	until := ctx.Formx(`until`).String()
	if len(until) > 0 {
		args.Add(`until`, until) // until = 24h (删除24小时前的镜像)
	}
	report, err := c.ImagesPrune(ctx, args)
	if err == nil {
		deletedImages := make([]string, 0, len(report.ImagesDeleted))
		untaggedImages := make([]string, 0, len(report.ImagesDeleted))
		for _, value := range report.ImagesDeleted {
			if len(value.Deleted) > 0 {
				deletedImages = append(deletedImages, value.Deleted)
			} else if len(value.Untagged) > 0 {
				untaggedImages = append(untaggedImages, value.Untagged)
			}
		}
		handler.SendOk(ctx, ctx.T(
			"操作成功。\n删除镜像: %s\nUntagged镜像: %s\n收回空间: %s",
			strings.Join(deletedImages, `, `),
			strings.Join(untaggedImages, `, `),
			com.FormatBytes(report.SpaceReclaimed),
		))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/image/index`))
}

func Delete(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	opts := types.ImageRemoveOptions{
		Force:         ctx.Formx(`force`).Bool(),
		PruneChildren: ctx.Formx(`pruneChildren`).Bool(),
	}
	imageID := ctx.Param(`id`)
	var result []types.ImageDeleteResponseItem
	if imageID == `0` {
		errs := common.NewErrors()
		for _, imageID := range ctx.FormValues(`id[]`) {
			if len(imageID) == 0 {
				continue
			}
			result, err = c.ImageRemove(ctx, imageID, opts)
			if err != nil {
				errs.Add(err)
				continue
			}
			ctx.Logger().Debugf(`ImageRemove: %+v`, result)
		}
		err = errs.ToError()
	} else {
		result, err = c.ImageRemove(ctx, imageID, opts)
		if err == nil {
			ctx.Logger().Debugf(`ImageRemove: %+v`, result)
		}
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/image/index`))
}
