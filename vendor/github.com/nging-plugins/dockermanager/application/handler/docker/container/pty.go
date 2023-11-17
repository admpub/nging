package container

import (
	"io"

	"github.com/admpub/websocket"
	"github.com/docker/docker/api/types"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/request"
)

func Resize(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	req := echo.GetValidated(ctx).(*request.ContainerResize)
	opts := types.ResizeOptions{
		Width:  req.Width,
		Height: req.Height,
	}
	err = c.ContainerResize(ctx, containerID, opts)
	if err != nil {
		return err
	}
	return ctx.JSON(ctx.Data().SetInfo(ctx.T(`操作成功`), code.Success.Int()))
}

func Pty(conn *websocket.Conn, ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	opts := types.ExecConfig{
		Tty:          true,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"/bin/sh"},
	}
	idResp, err := c.ContainerExecCreate(ctx, containerID, opts)
	if err != nil {
		return err
	}
	startOpts := types.ExecStartCheck{
		Detach: false,
		Tty:    true,
	}
	resp, err := c.ContainerExecAttach(ctx, idResp.ID, startOpts)
	if err != nil {
		return err
	}
	defer resp.Close()

	go io.Copy(NewWSWriter(conn), resp.Conn)
	for {
		mt, message, err := conn.ReadMessage()
		if mt == -1 || err != nil {
			return err
		}
		_, err = resp.Conn.Write(message)
		if err != nil {
			return err
		}
	}
}

func NewWSWriter(ws *websocket.Conn, msgTypes ...int) *wsWriter {
	var msgType int
	if len(msgTypes) > 0 {
		msgType = msgTypes[0]
	}
	if msgType <= 0 {
		msgType = websocket.BinaryMessage
	}
	return &wsWriter{ws: ws, msgType: msgType}
}

type wsWriter struct {
	ws      *websocket.Conn
	msgType int
}

func (w *wsWriter) Write(p []byte) (int, error) {
	err := w.ws.WriteMessage(w.msgType, p)
	return len(p), err
}
