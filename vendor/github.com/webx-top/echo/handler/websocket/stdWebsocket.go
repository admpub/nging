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
package websocket

import (
	"github.com/admpub/log"
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
)

type StdOptions struct {
	Handle   func(*websocket.Conn, echo.Context) error
	Upgrader *websocket.Upgrader
	Validate func(echo.Context) error
	Prefix   string
}

func (o StdOptions) Wrapper(e *echo.Echo) {
	e.Any(o.Prefix, StdWebsocket(o.Handle, o.Validate, o.Upgrader))
}

type StdHandler interface {
	Handle(*websocket.Conn, echo.Context) error
	Upgrader() *websocket.Upgrader
	Validate(echo.Context) error
}

var (
	DefaultStdUpgrader = &websocket.Upgrader{}
	DefaultExecuter    = func(c *websocket.Conn, ctx echo.Context) (err error) {
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				return err
			}
			log.Infof("Websocket recv: %s", message)

			if err = c.WriteMessage(mt, message); err != nil {
				return err
			}
		}
	}
)

func StdHanderWrapper(v interface{}) echo.Handler {
	if h, ok := v.(func(*websocket.Conn, echo.Context) error); ok {
		return StdWebsocket(h, nil)
	}
	if h, ok := v.(StdHandler); ok {
		return StdWebsocket(h.Handle, h.Validate, h.Upgrader())
	}
	if h, ok := v.(StdOptions); ok {
		return StdWebsocket(h.Handle, h.Validate, h.Upgrader)
	}
	if h, ok := v.(*StdOptions); ok {
		return StdWebsocket(h.Handle, h.Validate, h.Upgrader)
	}
	return nil
}

func StdWebsocket(executer func(*websocket.Conn, echo.Context) error, validate func(echo.Context) error, opts ...*websocket.Upgrader) echo.HandlerFunc {
	var opt *websocket.Upgrader
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt == nil {
		opt = DefaultStdUpgrader
	}
	if executer == nil {
		//Test mode
		executer = DefaultExecuter
	}
	h := func(ctx echo.Context) error {
		if validate != nil {
			if err := validate(ctx); err != nil {
				return err
			}
		}
		w := ctx.Response().StdResponseWriter()
		r := ctx.Request().StdRequest()
		c, err := opt.Upgrade(w, r, nil)
		if err != nil {
			return err
		}
		defer c.Close()

		return executer(c, ctx)
	}
	return echo.HandlerFunc(h)
}
