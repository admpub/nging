package captcha

import (
	"html/template"

	"github.com/webx-top/echo"
)

type ICaptcha interface {
	Render(ctx echo.Context, args ...interface{}) template.HTML
	Verify(ctx echo.Context, hostAlias string, name string, args ...string) echo.Data
	MakeData(ctx echo.Context, hostAlias string, name string) echo.H
}
