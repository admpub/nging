package user

import (
	"time"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/webx-top/echo"
)

func ThemeSwitch(ctx echo.Context) error {
	themeColor := ctx.Cookie().Get(`ThemeColor`)
	if len(themeColor) == 0 {
		themeColor = `light`
	}
	if themeColor == `light` {
		themeColor = `dark`
	} else {
		themeColor = `light`
	}
	ctx.Cookie().Set(`ThemeColor`, themeColor, time.Now().AddDate(1, 0, 0))
	next := ctx.Referer()
	next = echo.GetOtherURL(ctx, next)
	if len(next) == 0 {
		next = `/`
	}
	handler.SendOk(ctx, ctx.T(`已切换为%s模式`, themeColor))
	return ctx.Redirect(next)
}
