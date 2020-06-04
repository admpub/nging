package defaults

import (
	"context"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/mock"
)

func MustGetContext(ctx context.Context) echo.Context {
	eCtx, ok := ctx.(echo.Context)
	if !ok {
		eCtx = NewMockContext()
		if ctx != nil {
			eCtx.SetStdContext(ctx)
		}
	}
	return eCtx
}

func NewMockContext(args ...*echo.Echo) echo.Context {
	var e *echo.Echo
	if len(args) > 0 {
		e = args[0]
	} else {
		e = Default
	}
	return echo.NewContext(mock.NewRequest(), mock.NewResponse(), e)
}
