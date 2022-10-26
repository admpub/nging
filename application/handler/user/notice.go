/*
Nging is a toolbox for webmasters
Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
package user

import (
	"encoding/json"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/dbschema"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/defaults"
)

func init() {
	notice.OnOpen(func(user string) {
		ctx := defaults.NewMockContext()
		userM := dbschema.NewNgingUser(ctx)
		err := userM.UpdateField(nil, `online`, `Y`, `username`, user)
		if err != nil {
			log.Error(err)
		}
	})
	notice.OnClose(func(user string) {
		ctx := defaults.NewMockContext()
		userM := dbschema.NewNgingUser(ctx)
		err := userM.UpdateField(nil, `online`, `N`, `username`, user)
		if err != nil {
			log.Error(err)
		}
	})
}

func Notice(c *websocket.Conn, ctx echo.Context) error {
	user := handler.User(ctx)
	if user == nil {
		return ctx.NewError(code.Unauthenticated, `登录信息获取失败，请重新登录`)
	}
	oUser, clientID := notice.OpenClient(user.Username)
	defer notice.CloseClient(user.Username, clientID)
	//push(writer)
	go func(user *dbschema.NgingUser, clientID string) {
		message, err := json.Marshal(notice.NewMessage().SetMode(`-`).SetType(`clientID`).SetClientID(clientID))
		if err != nil {
			handler.WebSocketLogger.Error(`Push error: `, err.Error())
			return
		}
		handler.WebSocketLogger.Debug(`Push message: `, string(message))
		if err = c.WriteMessage(websocket.TextMessage, message); err != nil {
			handler.WebSocketLogger.Error(`Push error: `, err.Error())
			return
		}
		msgChan := oUser.Recv(clientID)
		if msgChan == nil {
			return
		}
		for {
			//message := []byte(echo.Dump(notice.NewMessageWithValue(`type`, `title`, `content:`+time.Now().Format(time.RFC1123)), false))
			//time.Sleep(time.Second)
			message := <-msgChan
			if message == nil {
				return
			}
			msgBytes, err := json.Marshal(message)
			if err != nil {
				handler.WebSocketLogger.Error(`Push error (json.Marshal): `, err.Error())
				return
			}
			handler.WebSocketLogger.Debug(`Push message: `, string(msgBytes))
			if err = c.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				handler.WebSocketLogger.Error(`Push error: `, err.Error())
				return
			}
		}
	}(user, clientID)

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
