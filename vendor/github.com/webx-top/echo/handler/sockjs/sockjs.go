package sockjs

import (
	"strings"

	"github.com/admpub/sockjs-go/sockjs"
	"github.com/lunny/log"
	"github.com/webx-top/echo"
)

type Options struct {
	Handle   func(sockjs.Session) error
	Options  *sockjs.Options
	Validate func(echo.Context) error
	Prefix   string
}

func (o Options) Wrapper(e echo.RouteRegister) {
	if o.Options == nil {
		o.Options = &sockjs.DefaultOptions
	}
	e.Any(strings.TrimRight(o.Prefix, "/")+"/*", Websocket(o.Prefix, o.Handle, o.Validate, o.Options))
}

type Handler interface {
	Handle(sockjs.Session) error
	Options() *sockjs.Options
	Validate(echo.Context) error
	Prefix() string
}

var DefaultExecuter = func(session sockjs.Session) error {
	for {
		msg, err := session.Recv()
		if err != nil {
			return err
		}
		err = session.Send(msg)
		if err != nil {
			return err
		}
	}
}

func HanderWrapper(v interface{}) echo.Handler {
	if h, ok := v.(func(sockjs.Session) error); ok {
		return Websocket(``, h, nil)
	}
	if h, ok := v.(Handler); ok {
		return Websocket(h.Prefix(), h.Handle, h.Validate, h.Options())
	}
	if h, ok := v.(Options); ok {
		return Websocket(h.Prefix, h.Handle, h.Validate, h.Options)
	}
	if h, ok := v.(*Options); ok {
		return Websocket(h.Prefix, h.Handle, h.Validate, h.Options)
	}
	return nil
}

func Websocket(prefix string, executer func(sockjs.Session) error, validate func(echo.Context) error, opts ...*sockjs.Options) echo.HandlerFunc {
	var opt sockjs.Options
	if len(opts) > 0 && opts[0] != nil {
		opt = *opts[0]
	} else {
		opt = sockjs.DefaultOptions
	}
	if executer == nil {
		//Test mode
		executer = DefaultExecuter
	}

	handler := sockjs.NewHandler(prefix, opt, func(session sockjs.Session) {
		err := executer(session)
		if err != nil {
			log.Debug(err)
		}
		session.Close(1024, "close")
	})
	h := func(ctx echo.Context) (err error) {
		if validate != nil {
			if err = validate(ctx); err != nil {
				return
			}
		}
		w := ctx.Response().StdResponseWriter()
		r := ctx.Request().StdRequest()
		handler.ServeHTTP(w, r)
		return
	}
	return echo.HandlerFunc(h)
}
