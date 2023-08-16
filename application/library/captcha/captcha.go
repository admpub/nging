package captcha

import (
	"html/template"

	"github.com/webx-top/echo"
)

type ICaptcha interface {
	// keysValues: key1, value1, key2, value2
	Render(ctx echo.Context, templatePath string, keysValues ...interface{}) template.HTML
	Verify(ctx echo.Context, hostAlias string, name string, captchaIdent ...string) echo.Data
	MakeData(ctx echo.Context, hostAlias string, name string) echo.H
}
