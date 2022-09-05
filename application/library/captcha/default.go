package captcha

import (
	"html/template"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
)

var dflt ICaptcha = &defaultCaptcha{}

type defaultCaptcha struct {
}

func (c *defaultCaptcha) Render(ctx echo.Context, args ...interface{}) template.HTML {
	options := tplfunc.MakeMap(args)
	options.Set("captchaId", common.GetHistoryOrNewCaptchaID(ctx))
	return tplfunc.CaptchaFormWithURLPrefix(ctx.Echo().Prefix(), options)
}

func (c *defaultCaptcha) Verify(ctx echo.Context, hostAlias string, name string, args ...string) echo.Data {
	return common.VerifyCaptcha(ctx, hostAlias, name, args...)
}

func (c *defaultCaptcha) MakeData(ctx echo.Context, hostAlias string, name string) echo.H {
	cid := common.GetHistoryOrNewCaptchaID(ctx)
	return common.CaptchaInfo(hostAlias, name, cid)
}
