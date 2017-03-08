/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package user

import (
	"errors"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/notice"
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
	ws "github.com/webx-top/echo/handler/websocket"
)

func init() {
	handler.RegisterToGroup(`/manage`, func(g *echo.Group) {
		wsOpts := ws.Options{
			Handle: Notice,
			Prefix: "/notice",
		}
		wsOpts.Wrapper(g)
	})
}

func Notice(c *websocket.Conn, ctx echo.Context) error {
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
			handler.WebSocketLogger.Debug(`Push message: `, string(message))
			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				handler.WebSocketLogger.Error(`Push error: `, err.Error())
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
		handler.WebSocketLogger.Error(err)
	}
	return nil
}
