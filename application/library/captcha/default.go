package captcha

import (
	"html/template"

	"github.com/webx-top/echo"
	hdlCaptcha "github.com/webx-top/echo/handler/captcha"
	"github.com/webx-top/echo/middleware/tplfunc"
)

var dflt ICaptcha = &defaultCaptcha{}

type defaultCaptcha struct {
}

func (c *defaultCaptcha) Init(_ echo.H) error {
	return nil
}

// keysValues: key1, value1, key2, value2
func (c *defaultCaptcha) Render(ctx echo.Context, templatePath string, keysValues ...interface{}) template.HTML {
	options := tplfunc.MakeMap(keysValues)
	options.Set("captchaId", GetHistoryOrNewCaptchaID(ctx, hdlCaptcha.DefaultOptions))
	if len(templatePath) == 0 {
		return tplfunc.CaptchaFormWithURLPrefix(ctx.Echo().Prefix(), options)
	}
	options.Set("captchaImage", tplfunc.CaptchaFormWithURLPrefix(ctx.Echo().Prefix(), options))
	return RenderTemplate(ctx, TypeDefault, templatePath, options)
}

func (c *defaultCaptcha) Verify(ctx echo.Context, hostAlias string, name string, captchaIdent ...string) echo.Data {
	return verifyAndSetDefaultCaptcha(ctx, hostAlias, name, captchaIdent...)
}

func (c *defaultCaptcha) MakeData(ctx echo.Context, hostAlias string, name string) echo.H {
	cid := GetHistoryOrNewCaptchaID(ctx, hdlCaptcha.DefaultOptions)
	return defaultCaptchaInfo(hostAlias, name, cid)
}
