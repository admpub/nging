package common

import (
	"github.com/webx-top/captcha"
	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"
	hdlCaptcha "github.com/webx-top/echo/handler/captcha"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/subdomains"
)

func GenCaptchaError(ctx echo.Context, hostAlias string, captchaName string, id string, args ...string) echo.Data {
	data := ctx.Data()
	data.SetZone(captchaName)
	data.SetData(CaptchaInfo(hostAlias, captchaName, id, args...))
	data.SetError(ErrCaptcha)
	return data
}

func GenAndRecordCaptchaID(ctx echo.Context, opt *hdlCaptcha.Options) string {
	cid := captcha.New()
	if len(opt.CookieName) > 0 {
		ctx.SetCookie(opt.CookieName, cid)
	}
	if len(opt.HeaderName) > 0 {
		ctx.Response().Header().Set(opt.HeaderName, cid)
	}
	return cid
}

func GetHistoryOrNewCaptchaId(ctx echo.Context) string {
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
		id = GenAndRecordCaptchaID(ctx, opt)
	}
	return id
}

func GetCaptchaId(ctx echo.Context, id string) (string, error) {
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

// VerifyCaptcha 验证码验证
func VerifyCaptcha(ctx echo.Context, hostAlias string, captchaName string, args ...string) echo.Data {
	idGet := ctx.Form
	idSet := func(id string) {
		ctx.Request().Form().Set(`captchaId`, id)
	}
	if len(args) > 0 {
		idGet = func(_ string, defaults ...string) string {
			return ctx.Form(args[0], defaults...)
		}
		idSet = func(id string) {
			ctx.Request().Form().Set(args[0], id)
		}
	}
	code := ctx.Form(captchaName)
	id := idGet("captchaId")
	if len(id) == 0 {
		return GenCaptchaError(ctx, hostAlias, captchaName, id, args...)
	}
	if len(code) == 0 {
		return ctx.Data().SetError(ErrCaptchaIdMissing)
	}
	newId, err := GetCaptchaId(ctx, id)
	if err != nil {
		if err != echo.ErrNotFound {
			return ctx.Data().SetError(err)
		}
	} else {
		if newId != id {
			idSet(id)
		}
	}
	if !tplfunc.CaptchaVerify(code, idGet) {
		return GenCaptchaError(ctx, hostAlias, captchaName, GenAndRecordCaptchaID(ctx, hdlCaptcha.DefaultOptions), args...)
	}
	return ctx.Data().SetCode(stdCode.Success.Int())
}

// VerifyAndSetCaptcha 验证码验证并设置新验证码信息
func VerifyAndSetCaptcha(ctx echo.Context, hostAlias string, captchaName string, args ...string) echo.Data {
	data := VerifyCaptcha(ctx, hostAlias, captchaName, args...)
	if data.GetCode() != stdCode.CaptchaError {
		id := GetHistoryOrNewCaptchaId(ctx)
		data.SetData(CaptchaInfo(hostAlias, captchaName, id, args...))
	}
	return data
}

// CaptchaInfo 新验证码信息
func CaptchaInfo(hostAlias string, captchaName string, captchaID string, args ...string) echo.H {
	if len(captchaID) == 0 {
		captchaID = captcha.New()
	}
	captchaIdent := `captchaId`
	if len(args) > 0 {
		captchaIdent = args[0]
	}
	return echo.H{
		`captchaName`:  captchaName,
		`captchaIdent`: captchaIdent,
		`captchaID`:    captchaID,
		`captchaURL`:   subdomains.Default.URL(`/captcha/`+captchaID+`.png`, hostAlias),
	}
}
