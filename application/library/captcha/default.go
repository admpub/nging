package captcha

import (
	"html/template"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
)

var dflt ICaptcha = &defaultCaptcha{}

type defaultCaptcha struct {
}

// keysValues: key1, value1, key2, value2
func (c *defaultCaptcha) Render(ctx echo.Context, keysValues ...interface{}) template.HTML {
	options := tplfunc.MakeMap(keysValues)
	options.Set("captchaId", common.GetHistoryOrNewCaptchaID(ctx))
	return tplfunc.CaptchaFormWithURLPrefix(ctx.Echo().Prefix(), options)
}

func (c *defaultCaptcha) Verify(ctx echo.Context, hostAlias string, name string, captchaIdent ...string) echo.Data {
	return common.VerifyCaptcha(ctx, hostAlias, name, captchaIdent...)
}

func (c *defaultCaptcha) MakeData(ctx echo.Context, hostAlias string, name string) echo.H {
	cid := common.GetHistoryOrNewCaptchaID(ctx)
	return common.CaptchaInfo(hostAlias, name, cid)
}
