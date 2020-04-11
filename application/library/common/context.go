package common

import (
	"context"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/engine/mock"
)

func MustGetMyContext(ctx context.Context) echo.Context {
	eCtx, ok := ctx.(echo.Context)
	if !ok {
		eCtx = NewMockContext()
		if ctx != nil {
			eCtx.SetStdContext(ctx)
		}
	} 
	return eCtx
}

func NewMockContext() echo.Context {
	return echo.NewContext(mock.NewRequest(), mock.NewResponse(), defaults.Default)
}
