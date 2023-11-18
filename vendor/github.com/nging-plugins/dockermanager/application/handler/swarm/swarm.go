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
	"strings"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/errdefs"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
)

func detectSwarmError(ctx echo.Context, err error) error {
	if errdefs.IsUnavailable(err) {
		ctx.Set(`message`, err.Error())
		ctx.Set(`title`, ctx.T(`Swarm设置`))
		return ctx.Render(`docker/swarm/unavailable`, nil)
	}
	return err
}

func Index(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	data, err := c.SwarmInspect(ctx)
	if err != nil {
		return detectSwarmError(ctx, err)
	}
	ctx.Set(`title`, ctx.T(`Swarm信息`))
	ctx.Set(`detail`, data)
	return ctx.Render(`docker/swarm/info`, handler.Err(ctx, err))
}

func SwarmInit(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		initReq := swarm.InitRequest{}
		err = ctx.MustBind(&initReq)
		if err != nil {
			goto END
		}
		specLabelsText := ctx.Formx(`spec[labelsText]`).String()
		if len(specLabelsText) > 0 {
			specLabels := map[string]string{}
			for _, labelRow := range strings.Split(specLabelsText, "\n") {
				labelRow = strings.TrimSpace(labelRow)
				if len(labelRow) == 0 {
					continue
				}
				parts := strings.SplitN(labelRow, `:`, 2)
				for k, v := range parts {
					parts[k] = strings.TrimSpace(v)
				}
				if len(parts[0]) > 0 {
					specLabels[parts[0]] = parts[1]
				}
			}
			initReq.Spec.Labels = specLabels
		}
		var result string
		result, err = c.SwarmInit(ctx, initReq)
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, result)
		return ctx.Redirect(handler.URLFor(`/docker/swarm/index`))
	}

END:
	ctx.Set(`activeURL`, `/docker/swarm/index`)
	ctx.Set(`title`, ctx.T(`初始化Swarm集群`))
	return ctx.Render(`docker/swarm/init`, handler.Err(ctx, err))
}

func SwarmJoin(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		joinReq := swarm.JoinRequest{}
		err = ctx.MustBind(&joinReq)
		if err != nil {
			goto END
		}
		remoteAddrsText := ctx.Formx(`remoteAddrsText`).String()
		if len(remoteAddrsText) > 0 {
			remoteAddrs := []string{}
			for _, remoteAddr := range strings.Split(remoteAddrsText, "\n") {
				remoteAddr = strings.TrimSpace(remoteAddr)
				if len(remoteAddr) == 0 {
					continue
				}
				remoteAddrs = append(remoteAddrs, remoteAddr)
			}
			joinReq.RemoteAddrs = remoteAddrs
		}
		err = c.SwarmJoin(ctx, joinReq)
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, ctx.T(`加入成功`))
		return ctx.Redirect(handler.URLFor(`/docker/swarm/index`))
	}

END:
	if ctx.Get(`activeURL`) == nil {
		ctx.Set(`activeURL`, `/docker/swarm/index`)
	}
	if ctx.Get(`title`) == nil {
		ctx.Set(`title`, ctx.T(`加入Swarm集群`))
	}
	return ctx.Render(`docker/swarm/join`, handler.Err(ctx, err))
}

func SwarmLeave(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	force := ctx.Formx(`force`).Bool()
	err = c.SwarmLeave(ctx, force)
	if err != nil {
		return err
	}
	handler.SendOk(ctx, ctx.T(`离开成功`))
	return ctx.Redirect(handler.URLFor(`/docker/swarm/index`))
}
