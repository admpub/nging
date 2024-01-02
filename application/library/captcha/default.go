package captcha

import (
	"html/template"
	"path"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
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
	options.Set("captchaId", getHistoryOrNewCaptchaID(ctx))
	if len(templatePath) == 0 {
		return tplfunc.CaptchaFormWithURLPrefix(ctx.Echo().Prefix(), options)
	}
	options.Set("captchaImage", tplfunc.CaptchaFormWithURLPrefix(ctx.Echo().Prefix(), options))
	b, err := ctx.Fetch(path.Join(`captcha/default`, templatePath), options)
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(com.Bytes2str(b))
}

func (c *defaultCaptcha) Verify(ctx echo.Context, hostAlias string, name string, captchaIdent ...string) echo.Data {
	return verifyAndSetDefaultCaptcha(ctx, hostAlias, name, captchaIdent...)
}

func (c *defaultCaptcha) MakeData(ctx echo.Context, hostAlias string, name string) echo.H {
	cid := getHistoryOrNewCaptchaID(ctx)
	return defaultCaptchaInfo(hostAlias, name, cid)
}
