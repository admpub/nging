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
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
)

func New(prefix string, handler func(*websocket.Conn, echo.Context) error) *Options {
	return &Options{Handle: handler, Prefix: prefix}
}

type Options struct {
	Handle   func(*websocket.Conn, echo.Context) error
	Upgrader *websocket.EchoUpgrader
	Validate func(echo.Context) error
	Prefix   string
}

func (o *Options) SetPrefix(prefix string) *Options {
	o.Prefix = prefix
	return o
}

func (o *Options) SetHandler(handler func(*websocket.Conn, echo.Context) error) *Options {
	o.Handle = handler
	return o
}

func (o *Options) SetValidator(validator func(echo.Context) error) *Options {
	o.Validate = validator
	return o
}

func (o *Options) SetUpgrader(upgrader *websocket.EchoUpgrader) *Options {
	o.Upgrader = upgrader
	return o
}

func (o Options) Wrapper(e echo.RouteRegister) echo.IRouter {
	if o.Upgrader == nil {
		o.Upgrader = DefaultUpgrader
	}
	return e.Any(o.Prefix, Websocket(o.Handle, o.Validate, o.Upgrader))
}

type Handler interface {
	Handle(*websocket.Conn, echo.Context) error
	Upgrader() *websocket.EchoUpgrader
	Validate(echo.Context) error
}

var (
	DefaultUpgrader = &websocket.EchoUpgrader{}
)

func HanderWrapper(v interface{}) echo.Handler {
	if h, ok := v.(func(*websocket.Conn, echo.Context) error); ok {
		return Websocket(h, nil)
	}
	if h, ok := v.(Handler); ok {
		return Websocket(h.Handle, h.Validate, h.Upgrader())
	}
	if h, ok := v.(Options); ok {
		return Websocket(h.Handle, h.Validate, h.Upgrader)
	}
	if h, ok := v.(*Options); ok {
		return Websocket(h.Handle, h.Validate, h.Upgrader)
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

func Websocket(executer func(*websocket.Conn, echo.Context) error, validate func(echo.Context) error, opts ...*websocket.EchoUpgrader) echo.HandlerFunc {
	var opt *websocket.EchoUpgrader
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt == nil {
		opt = DefaultUpgrader
	}
	if executer == nil {
		//Test mode
		executer = DefaultExecuter
	}
	h := func(ctx echo.Context) (err error) {
		if validate != nil {
			if err = validate(ctx); err != nil {
				return
			}
		}
		return opt.Upgrade(ctx, func(conn *websocket.Conn) error {
			defer conn.Close()
			return executer(conn, ctx)
		}, nil)
	}
	return echo.HandlerFunc(h)
}
