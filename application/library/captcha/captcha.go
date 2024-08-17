package captcha

import (
	"html/template"
	"os"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type ICaptcha interface {
	Init(echo.H) error
	// keysValues: key1, value1, key2, value2
	Render(ctx echo.Context, templatePath string, keysValues ...interface{}) template.HTML
	Verify(ctx echo.Context, hostAlias string, name string, captchaIdent ...string) echo.Data
	MakeData(ctx echo.Context, hostAlias string, name string) echo.H
}

func renderTemplate(ctx echo.Context, captchaType string, templatePath string, options param.Store) template.HTML {
	tmplFile := fixTemplatePath(captchaType, templatePath)
	b, err := ctx.Fetch(tmplFile, options)
	if err != nil {
		if templatePath != `default` && os.IsNotExist(err) {
			b, err = ctx.Fetch(strings.TrimSuffix(tmplFile, templatePath)+`default`, options)
		}
		if err != nil {
			return template.HTML(err.Error())
		}
	}
	return template.HTML(com.Bytes2str(b))
}
