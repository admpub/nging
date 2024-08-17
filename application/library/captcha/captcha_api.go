package captcha

import (
	"html/template"
	"strconv"

	"github.com/admpub/captcha-go"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/com/formatter"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func newCaptchaAPI() ICaptcha {
	return &captchaAPI{}
}

type captchaAPI struct {
	provider  string
	endpoint  captcha.Endpoint
	siteKey   string
	verifier  *captcha.SimpleCaptchaVerifier
	jsURL     string
	captchaID string
	cfg       echo.H
}

func (c *captchaAPI) Init(opt echo.H) error {
	c.provider = opt.String(`provider`)
	c.cfg = opt.GetStore(c.provider)
	c.siteKey = c.cfg.String(`siteKey`)
	switch c.provider {
	case `turnstile`:
		c.endpoint = captcha.CloudflareTurnstile
		c.jsURL = `https://challenges.cloudflare.com/turnstile/v0/api.js`
	default:
		c.endpoint = captcha.GoogleRecaptcha
		c.jsURL = `https://www.recaptcha.net/recaptcha/api.js?render=` + c.siteKey
	}
	captchaSecret := c.cfg.String(`secret`)
	//echo.Dump(echo.H{`siteKey`: c.siteKey, `secret`: captchaSecret})
	v := captcha.NewCaptchaVerifier(c.endpoint, captchaSecret)
	c.verifier = &captcha.SimpleCaptchaVerifier{
		Verifier: *v,
	}
	return nil
}

// keysValues: key1, value1, key2, value2
func (c *captchaAPI) Render(ctx echo.Context, templatePath string, keysValues ...interface{}) template.HTML {
	options := tplfunc.MakeMap(keysValues)
	options.Set("siteKey", c.siteKey)
	options.Set("endpoint", c.endpoint)
	options.Set("provider", c.provider)
	initedKey := `CaptchaJSInited.` + c.provider
	var jsURL string
	if !ctx.Internal().Bool(initedKey) {
		ctx.Internal().Set(initedKey, true)
		jsURL = c.jsURL
	}
	options.Set("jsURL", jsURL)
	c.captchaID = com.RandomAlphanumeric(16)
	options.Set("captchaID", c.captchaID)
	if len(templatePath) == 0 {
		var htmlContent string
		switch c.endpoint {
		case captcha.CloudflareTurnstile:
			if len(jsURL) > 0 {
				htmlContent += `<script src="` + jsURL + `" async defer></script>`
			}
			locationID := `turnstile-` + c.captchaID
			htmlContent += `<input type="hidden" name="captchaId" value="` + c.captchaID + `" />`
			htmlContent += `<div class="cf-turnstile" id="` + locationID + `" data-sitekey="` + c.siteKey + `"></div>`
			htmlContent += `<input type="hidden" id="` + locationID + `-extend" disabled />`
			htmlContent += `<script>
window.addEventListener('load', function(){
	$('#` + locationID + `').closest('.input-group-addon').addClass('xxs-padding-top').prev('input').remove();
	var $form=$('#` + locationID + `').closest('form');
	$form.on('submit',function(e){
		if($('#` + locationID + `').data('lastGeneratedAt')>(new Date()).getTime()-290) {
			$('#` + locationID + `').data('lastGeneratedAt',0);
			return true;
		}
		window.setTimeout(function(){
			turnstile.reset('#` + locationID + `');
		},1000);
		$('#` + locationID + `').data('lastGeneratedAt',(new Date()).getTime());
	});
})
</script>`
		default:
			locationID := `recaptcha-` + c.captchaID
			htmlContent = `<input type="hidden" name="captchaId" value="` + c.captchaID + `" />`
			htmlContent += `<input type="hidden" id="` + locationID + `" name="g-recaptcha-response" value="" />`
			htmlContent += `<input type="hidden" id="` + locationID + `-extend" disabled />`
			if len(jsURL) > 0 {
				htmlContent += `<script src="` + jsURL + `" async defer></script>`
			}
			htmlContent += `<script>
window.addEventListener('load', function(){
	grecaptcha.ready(function() {
	  grecaptcha.execute('` + c.siteKey + `', {action: 'submit'}).then(function(token) {
		$('#` + locationID + `').val(token);
		$('#` + locationID + `').data('lastGeneratedAt',(new Date()).getTime());
	  });
	});
	var igrp=$('#` + locationID + `').closest('.input-group');
	if(igrp.length>0){
		igrp.hide();
		if(igrp.parent().hasClass('form-group')) igrp.parent().hide();
	}
	$('#` + locationID + `').closest('.input-group-addon').prev('input').remove();
	var $form=$('#` + locationID + `').closest('form');
	var $submit=$form.find(':submit');
	$submit.on('click',function(e){
		if($('#` + locationID + `').val() && $('#` + locationID + `').data('lastGeneratedAt')>(new Date()).getTime()-110) {
			$('#` + locationID + `').data('lastGeneratedAt',0);
			return true;
		}
		var $this=$(this);
		e.preventDefault();
		grecaptcha.execute('` + c.siteKey + `', {action: 'submit'}).then(function(token) {
		  $('#` + locationID + `').val(token);
		  $('#` + locationID + `').data('lastGeneratedAt',(new Date()).getTime());
		  $this.trigger('click');
		});
	});
})
</script>`
		}
		return template.HTML(htmlContent)
	}
	return RenderTemplate(ctx, TypeAPI, templatePath, options)
}

func (c *captchaAPI) Verify(ctx echo.Context, hostAlias string, _ string, _ ...string) echo.Data {
	var name string
	switch c.endpoint {
	case captcha.CloudflareTurnstile:
		name = `cf-turnstile-response`
	default:
		name = `g-recaptcha-response`
		if c.cfg.Has(`minScore`) {
			c.verifier.MinScore = c.cfg.Float32(`minScore`)
		} else {
			c.verifier.MinScore = 0.5
		}
		c.verifier.ExpectedAction = ctx.Form(`captchaAction`, `submit`)
	}
	c.captchaID = ctx.Formx(`captchaId`).String()
	if len(c.captchaID) == 0 {
		return GenCaptchaError(ctx, ErrCaptchaIdMissing, name, c.MakeData(ctx, hostAlias, name))
	}
	token := ctx.Form(name)
	if len(token) == 0 { // 为空说明没有验证码
		return ctx.Data().SetError(ErrCaptchaCodeRequired.SetMessage(ctx.T(`请先进行人机验证`)).SetZone(name))
	}
	c.verifier.ExpectedHostname = ctx.Domain()
	var clientIP string
	if c.cfg.Bool(`verifyIP`) {
		clientIP = ctx.RealIP()
	}
	resp, ok, err := c.verifier.VerifyActionWithResponse(token, clientIP, c.verifier.ExpectedAction)
	if err != nil {
		return GenCaptchaError(ctx, err, name, c.MakeData(ctx, hostAlias, name))
	}
	if !ok {
		log.Warnf(`failed to captchaAPI.Verify: %s`, formatter.AsStringer(resp))
		return GenCaptchaError(ctx, ErrCaptcha.SetMessage(ctx.T(`抱歉，未能通过人机验证`)), name, c.MakeData(ctx, hostAlias, name))
	}
	return ctx.Data().SetCode(code.Success.Int())
}

func (c *captchaAPI) MakeData(ctx echo.Context, hostAlias string, name string) echo.H {
	data := echo.H{}
	data.Set("siteKey", c.siteKey)
	data.Set("provider", c.provider)
	data.Set("jsURL", c.jsURL)
	if len(c.captchaID) == 0 {
		c.captchaID = com.RandomAlphanumeric(16)
	}
	data.Set("captchaType", TypeAPI)
	data.Set("captchaID", c.captchaID)
	var jsInit, jsCallback string
	var locationID string
	var htmlCode string
	var captchaName string
	switch c.provider {
	case `turnstile`:
		captchaName = `cf-turnstile-response`
		locationID = `turnstile-` + c.captchaID
		jsInit = `(function(){
	var f=function(){
		if(typeof(turnstile)!='undefined'){
			turnstile.render('#` + locationID + `');
		}else{
			window.setTimeout(f,200);
		}
	};
})();`
		jsCallback = `function(callback){
	callback && callback();
	window.setTimeout(function(){
		turnstile.reset('#` + locationID + `');
	},1000);
}`
		var theme string
		if ctx.Cookie().Get(`ThemeColor`) == `dark` {
			theme = `dark`
		} else {
			theme = `light`
		}
		htmlCode = `<input type="hidden" name="captchaId" value="` + c.captchaID + `" /><div class="cf-turnstile text-center" id="turnstile-` + c.captchaID + `" data-sitekey="` + c.siteKey + `" data-theme="` + theme + `"></div>`
	default:
		captchaName = `g-recaptcha-response`
		locationID = `recaptcha-` + c.captchaID
		defaultTips := strconv.Quote(ctx.T(`加载成功，请点击“提交”按钮继续`))
		jsInit = `grecaptcha.ready(function() {
	grecaptcha.execute('` + c.siteKey + `', {action: 'submit'}).then(function(token) {
		$('#` + locationID + `').val(token);
		$('#` + locationID + `').data('lastGeneratedAt',(new Date()).getTime());
		if($('#` + locationID + `-loading').length>0){
			var successTips = $('#` + locationID + `-loading').data('success-tips')||` + defaultTips + `;
			$('#` + locationID + `-loading').html('<i class="fa fa-check text-success"></i> '+successTips);
		}
	});
});`
		jsCallback = `function(callback){
	grecaptcha.execute('` + c.siteKey + `', {action: 'submit'}).then(function(token) {
		$('#` + locationID + `').val(token);
		$('#` + locationID + `').data('lastGeneratedAt',(new Date()).getTime());
		callback && callback(token);
	});
}`
		htmlCode = `<input type="hidden" name="captchaId" value="` + c.captchaID + `" /><input type="hidden" id="recaptcha-` + c.captchaID + `" name="g-recaptcha-response" value="" />`
	}
	data.Set("jsCallback", jsCallback)
	data.Set("jsInit", jsInit)
	data.Set("locationID", locationID)
	data.Set("html", htmlCode)
	data.Set("captchaIdent", `captchaId`)
	data.Set("captchaName", captchaName)
	return data
}
