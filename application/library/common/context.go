package common

import (
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

var (
	MustGetMyContext = defaults.MustGetContext
	MustGetContext   = defaults.MustGetContext
	NewMockContext   = defaults.NewMockContext
)

func Tx(ctx echo.Context) *factory.Transaction {
	t := ctx.Transaction().(*factory.Param).Trans()
	return t
}
