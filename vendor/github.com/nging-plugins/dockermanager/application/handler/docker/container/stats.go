package container

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"

	"github.com/admpub/websocket"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
)

func StatsPage(ctx echo.Context) error {
	ctx.Set(`activeURL`, `/docker/base/container/index`)
	ctx.Set(`title`, ctx.T(`容器统计信息`))
	return ctx.Render(`docker/base/container/stats`, nil)
}

func Stats(conn *websocket.Conn, ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	stream := ctx.Formx(`stream`, `true`).Bool()
	sts, err := c.ContainerStats(ctx, containerID, stream)
	if err != nil {
		return err
	}
	defer sts.Body.Close()
	buf := bufio.NewReader(sts.Body)
	for {
		message, err := buf.ReadString('\n')
		if err != nil {
			return err
		}
		message = strings.TrimSuffix(message, "\n")
		message = strings.TrimSuffix(message, "\r")
		if err = conn.WriteMessage(websocket.TextMessage, []byte(message+"\r\n")); err != nil {
			return err
		}
	}
}

func StatsOneShot(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	sts, err := c.ContainerStatsOneShot(ctx, containerID)
	if err != nil {
		return err
	}
	defer sts.Body.Close()
	b, err := io.ReadAll(sts.Body)
	if err != nil {
		return err
	}
	data := echo.H{`stats`: json.RawMessage(b)}
	return ctx.JSON(ctx.Data().SetData(data))
}

func Top(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	args := ctx.FormValues(`args`)
	topData, err := c.ContainerTop(ctx, containerID, args)
	if err != nil {
		return err
	}
	return ctx.JSON(ctx.Data().SetData(topData))
}
