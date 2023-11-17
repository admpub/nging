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

	"github.com/admpub/websocket"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/library/utils"
)

func TaskIndex(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	args := filters.NewArgs()
	// id / name / node/ desired-state
	id := ctx.Form(`id`)
	if len(id) > 0 {
		args.Add(`id`, id)
	}
	node := ctx.Form(`node`)
	if len(node) > 0 {
		args.Add(`node`, node)
	}
	desiredState := ctx.Form(`desiredState`)
	if len(desiredState) > 0 {
		args.Add(`desired-state`, desiredState)
	}
	name := ctx.Form(`name`)
	if len(name) > 0 {
		args.Add(`name`, name)
	}
	service := ctx.Form(`service`)
	if len(service) > 0 {
		args.Add(`service`, service)
	}
	list, err := c.TaskList(ctx, types.TaskListOptions{Filters: args})
	if err != nil {
		return detectSwarmError(ctx, err)
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`docker/swarm/task/index`, handler.Err(ctx, err))
}

func TaskDetail(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	taskID := ctx.Param(`id`)
	data, _, err := c.TaskInspectWithRaw(ctx, taskID)
	if err != nil {
		return err
	}
	ctx.Set(`activeURL`, `/docker/swarm/task/index`)
	ctx.Set(`title`, ctx.T(`任务信息`))
	ctx.Set(`detail`, data)
	return ctx.Render(`docker/swarm/task/detail`, err)
}

func TaskLogs(conn *websocket.Conn, ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	taskID := ctx.Param(`id`)
	opts := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}
	reader, err := c.TaskLogs(ctx, taskID, opts)
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
