package captcha

import (
	"github.com/webx-top/captcha"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	hdlCaptcha "github.com/webx-top/echo/handler/captcha"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/subdomains"
)

func GenCaptchaError(ctx echo.Context, err error, captchaName string, data interface{}) echo.Data {
	d := ctx.Data()
	d.SetZone(captchaName)
	d.SetData(data)
	d.SetError(err)
	return d
}

func genDefaultCaptchaError(ctx echo.Context, err error, hostAlias string, captchaName string, id string, captchaIdent ...string) echo.Data {
	return GenCaptchaError(ctx, err, captchaName, defaultCaptchaInfo(hostAlias, captchaName, id, captchaIdent...))
}

func defaultCaptchaGenAndRecordCaptchaID(ctx echo.Context, opt *hdlCaptcha.Options) string {
	cid := captcha.New()
	if len(opt.CookieName) > 0 {
		ctx.SetCookie(opt.CookieName, cid)
	}
	if len(opt.HeaderName) > 0 {
		ctx.Response().Header().Set(opt.HeaderName, cid)
	}
	return cid
}

func getHistoryOrNewCaptchaID(ctx echo.Context) string {
	opt := hdlCaptcha.DefaultOptions
	var (
		exists bool
		id     string
	)
	if len(opt.CookieName) > 0 {
		id = ctx.GetCookie(opt.CookieName)
		if len(id) > 0 {
			exists = captcha.Exists(id)
		}
	}
	if !exists && len(opt.HeaderName) > 0 {
		id = ctx.Header(opt.HeaderName)
		if len(id) > 0 {
			exists = captcha.Exists(id)
		}
	}
	if !exists {
		id = defaultCaptchaGenAndRecordCaptchaID(ctx, opt)
	}
	return id
}

func GetCaptchaID(ctx echo.Context, id string) (string, error) {
	opt := hdlCaptcha.DefaultOptions
	exists := captcha.Exists(id)
	if !exists && len(opt.CookieName) > 0 {
		id = ctx.GetCookie(opt.CookieName)
		if len(id) > 0 {
			exists = captcha.Exists(id)
		}
	}
	if !exists && len(opt.HeaderName) > 0 {
		id = ctx.Header(opt.HeaderName)
		if len(id) > 0 {
			exists = captcha.Exists(id)
		}
	}
	if !exists {
		return id, echo.ErrNotFound
	}
	return id, nil
}

// verifyDefaultCaptcha 验证码验证
func verifyDefaultCaptcha(ctx echo.Context, hostAlias string, captchaName string, captchaIdent ...string) echo.Data {
	var idGet func(name string, defaults ...string) string
	var idSet func(id string)
	if len(captchaIdent) > 0 {
		idGet = func(_ string, defaults ...string) string {
			return ctx.Form(captchaIdent[0], defaults...)
		}
		idSet = func(id string) {
			ctx.Request().Form().Set(captchaIdent[0], id)
		}
	} else {
		idGet = ctx.Form
		idSet = func(id string) {
			ctx.Request().Form().Set(`captchaId`, id)
		}
	}
	id := idGet("captchaId")
	if len(id) == 0 { // 为空说明表单没有显示验证码输入框，此时返回验证码信息供前端显示
		return genDefaultCaptchaError(ctx, ErrCaptchaIdMissing, hostAlias, captchaName, id, captchaIdent...)
	}
	vcode := ctx.Form(captchaName)
	if len(vcode) == 0 { // 为空说明没有输入验证码
		return ctx.Data().SetError(ErrCaptchaCodeRequired.SetZone(captchaName))
	}
	newId, err := GetCaptchaID(ctx, id)
	if err != nil {
		if err != echo.ErrNotFound {
			return ctx.Data().SetError(err)
		}
	} else {
		if newId != id {
			idSet(id)
		}
	}
	if !tplfunc.CaptchaVerify(vcode, idGet) {
		return genDefaultCaptchaError(ctx, ErrCaptcha, hostAlias, captchaName, defaultCaptchaGenAndRecordCaptchaID(ctx, hdlCaptcha.DefaultOptions), captchaIdent...)
	}
	return ctx.Data().SetCode(code.Success.Int())
}

// verifyAndSetDefaultCaptcha 验证码验证并设置新验证码信息
func verifyAndSetDefaultCaptcha(ctx echo.Context, hostAlias string, captchaName string, args ...string) echo.Data {
	data := verifyDefaultCaptcha(ctx, hostAlias, captchaName, args...)
	if data.GetCode() != code.CaptchaError {
		data.SetData(defaultCaptchaInfo(hostAlias, captchaName, getHistoryOrNewCaptchaID(ctx), args...))
	}
	return data
}

// defaultCaptchaInfo 新验证码信息
func defaultCaptchaInfo(hostAlias string, captchaName string, captchaID string, captchaIdent ...string) echo.H {
	if len(captchaID) == 0 {
		captchaID = captcha.New()
	}
	_captchaIdent := `captchaId`
	if len(captchaIdent) > 0 {
		_captchaIdent = captchaIdent[0]
	}
	return echo.H{
		`captchaType`:  TypeDefault,
		`captchaName`:  captchaName,
		`captchaIdent`: _captchaIdent,
		`captchaID`:    captchaID,
		`captchaURL`:   subdomains.Default.URL(`/captcha/`+captchaID+`.png`, hostAlias),
	}
}
