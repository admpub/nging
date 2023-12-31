package common

import (
	"html/template"

	"github.com/admpub/nging/v5/application/library/captcha"
	"github.com/webx-top/echo"
)

func GetCaptchaEngine(ctx echo.Context) (captcha.ICaptcha, error) {
	cfg := Setting(`captcha`)
	typ := cfg.String(`type`, `default`)
	create := captcha.Get(typ)
	if create == nil {
		return nil, echo.ErrNotImplemented
	}
	cpt := create()
	tcfg := cfg.Children(typ)
	err := cpt.Init(tcfg)
	return cpt, err
}

// VerifyCaptcha 验证码验证
func VerifyCaptcha(ctx echo.Context, hostAlias string, captchaName string, captchaIdent ...string) echo.Data {
	cpt, err := GetCaptchaEngine(ctx)
	if err != nil {
		return ctx.Data().SetError(err)
	}
	return cpt.Verify(ctx, hostAlias, captchaName, captchaIdent...)
}

// CaptchaForm 验证码输入表单
func CaptchaForm(ctx echo.Context, args ...interface{}) template.HTML {
	cpt, err := GetCaptchaEngine(ctx)
	if err != nil {
		return template.HTML(err.Error())
	}
	return cpt.Render(ctx, ``, args...)
}
