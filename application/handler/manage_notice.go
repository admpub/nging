package handler

import (
	"errors"

	"github.com/admpub/nging/application/library/notice"
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
)

func ManageNotice(c *websocket.Conn, ctx echo.Context) error {
	user, ok := ctx.Get(`user`).(string)
	if !ok {
		return errors.New(ctx.T(`登录信息获取失败，请重新登录`))
	}
	notice.OpenClient(user)
	defer notice.CloseClient(user)
	//push(writer)
	go func() {
		for {
			message := notice.RecvJSON(user)
			WebSocketLogger.Debug(`Push message: `, string(message))
			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				WebSocketLogger.Error(`Push error: `, err.Error())
				return
			}
		}
	}()

	//echo
	var execute = func(conn *websocket.Conn) error {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			if err = conn.WriteMessage(mt, message); err != nil {
				return err
			}
		}
	}
	err := execute(c)
	if err != nil {
		WebSocketLogger.Error(err)
	}
	return nil
}
