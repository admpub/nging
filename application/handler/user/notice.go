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
	"github.com/coscms/webcore/dbschema"

	"github.com/admpub/websocket"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/notice"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/defaults"
)

func init() {
	notice.OnOpen(func(user string) {
		ctx := defaults.NewMockContext()
		userM := dbschema.NewNgingUser(ctx)
		err := userM.UpdateField(nil, `online`, `Y`, db.And(
			db.Cond{`username`: user},
			db.Cond{`online`: `N`},
		))
		if err != nil {
			log.Errorf(`failed to userM.UpdateField(online=Y,username=%q): %v`, user, err)
		}
	})
	notice.OnClose(func(user string) {
		ctx := defaults.NewMockContext()
		userM := dbschema.NewNgingUser(ctx)
		err := userM.UpdateField(nil, `online`, `N`, db.And(
			db.Cond{`username`: user},
			db.Cond{`online`: `Y`},
		))
		if err != nil {
			log.Errorf(`failed to userM.UpdateField(online=N,username=%q): %v`, user, err)
		}
	})
}

func Notice(c *websocket.Conn, ctx echo.Context) error {
	user := backend.User(ctx)
	if user == nil {
		return ctx.NewError(code.Unauthenticated, `登录信息获取失败，请重新登录`)
	}
	close, msgChan, err := notice.Default().MakeMessageGetter(user.Username)
	if err != nil || msgChan == nil {
		return err
	}
	if close != nil {
		defer close()
	}
	//push(writer)
	go func() {
		for {
			//message := []byte(echo.Dump(notice.NewMessageWithValue(`type`, `title`, `content:`+time.Now().Format(time.RFC1123)), false))
			//time.Sleep(time.Second)
			message, ok := <-msgChan
			if !ok || message == nil {
				return
			}
			msgBytes, err := json.Marshal(message)
			message.Release()
			if err != nil {
				backend.WebSocketLogger.Error(`Push error (json.Marshal): `, err.Error())
				return
			}
			backend.WebSocketLogger.Debug(`Push message: `, string(msgBytes))
			if err = c.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway) {
					backend.WebSocketLogger.Debug(`Push error: `, err.Error())
				} else {
					backend.WebSocketLogger.Error(`Push error: `, err.Error())
				}
				return
			}
		}
	}()

	//echo
	execute := func(conn *websocket.Conn) error {
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
	err = execute(c)
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			backend.WebSocketLogger.Debug(err.Error())
		} else {
			backend.WebSocketLogger.Error(err.Error())
		}
	}
	return nil
}
