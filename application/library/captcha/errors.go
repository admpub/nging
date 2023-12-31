package captcha

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

var (
	//ErrCaptcha 验证码错误
	ErrCaptcha = echo.NewError(`Captcha is incorrect`, code.CaptchaError)
	//ErrCaptchaIdMissing 缺少captchaId
	ErrCaptchaIdMissing = echo.NewError(`Missing captchaId`, code.CaptchaIdMissing).SetZone(`captchaId`)
	//ErrCaptchaCodeRequired 验证码不能为空
	ErrCaptchaCodeRequired = echo.NewError(`Captcha code is required`, code.CaptchaCodeRequired).SetZone(`code`)
)
