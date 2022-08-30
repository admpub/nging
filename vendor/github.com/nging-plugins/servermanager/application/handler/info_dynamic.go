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

package handler

import (
	"bytes"
	"strings"

	"github.com/admpub/sockjs-go/v3/sockjs"
	"github.com/admpub/websocket"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v4/application/handler"

	"github.com/nging-plugins/servermanager/application/library/system"
)

func InfoBySockJS(c sockjs.Session) error {
	send := make(chan interface{})
	//push(writer)
	go func() {
		for {
			b, err := com.JSONEncode(<-send)
			if err != nil {
				handler.WebSocketLogger.Error(`Push error: `, err.Error())
				continue
			}
			message := string(b)
			if err := c.Send(message); err != nil {
				handler.WebSocketLogger.Error(`Push error: `, err.Error())
				return
			}
		}
	}()
	//echo
	exec := func(session sockjs.Session) error {
		for {
			message, err := session.Recv()
			if err != nil {
				return err
			}
			info := strings.SplitN(message, `:`, 2)
			switch message {
			case `ping`: // Net/Memory/CPU
				var n int
				if len(info) > 1 {
					n = param.AsInt(info[1])
				}
				send <- system.RealTimeStatusObject(n)
			case `pingAll`:
				info := &system.DynamicInformation{}
				send <- info.Init()
			}
		}
	}
	err := exec(c)
	if err != nil {
		handler.WebSocketLogger.Error(err)
	}
	close(send)
	return nil
}

func InfoByWebsocket(c *websocket.Conn, ctx echo.Context) error {
	send := make(chan interface{})
	//push(writer)
	go func() {
		for {
			message := <-send
			if err := c.WriteJSON(message); err != nil {
				handler.WebSocketLogger.Error(`Push error: `, err.Error())
				return
			}
		}
	}()
	//echo
	exec := func(conn *websocket.Conn) error {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			info := bytes.SplitN(message, []byte(`:`), 2)
			switch com.Bytes2str(info[0]) {
			case `ping`: // Net/Memory/CPU
				var n int
				if len(info) > 1 {
					n = param.AsInt(string(info[1]))
				}
				send <- system.RealTimeStatusObject(n)
			case `pingAll`:
				info := &system.DynamicInformation{}
				send <- info.Init()
			}
		}
	}
	err := exec(c)
	if err != nil {
		handler.WebSocketLogger.Error(err)
	}
	close(send)
	return nil
}
