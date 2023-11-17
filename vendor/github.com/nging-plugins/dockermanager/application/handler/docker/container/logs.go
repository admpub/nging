package container

import (
	"bufio"
	"strings"
	"time"

	"github.com/admpub/websocket"
	"github.com/docker/docker/api/types"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/library/utils"
)

func todayZeroClock() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

func timestampString(t time.Time) string {
	return param.AsString(t.Unix())
}

func Logs(conn *websocket.Conn, ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	opts := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Details:    true,
	}
	reader, err := c.ContainerLogs(ctx, containerID, opts)
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
