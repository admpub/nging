package captcha

import (
	"html/template"

	"github.com/admpub/captcha-go"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/com/formatter"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func newRecaptcha() ICaptcha {
	return &reCaptcha{}
}

type reCaptcha struct {
	endpoint captcha.Endpoint
	siteKey  string
	verifier *captcha.SimpleCaptchaVerifier
}

const HCaptcha captcha.Endpoint = `https://api.hcaptcha.com/siteverify`

func (c *reCaptcha) Init(opt echo.H) error {
	var endpoint captcha.Endpoint
	switch opt.String(`provider`) {
	case `cloudflare`:
		endpoint = captcha.CloudflareTurnstile
	case `hcaptcha`:
		endpoint = HCaptcha
	default:
		endpoint = captcha.GoogleRecaptcha
	}
	c.endpoint = endpoint
	c.siteKey = opt.String(`siteKey`)
	captchaSecret := opt.String(`secret`)
	//echo.Dump(echo.H{`siteKey`: c.siteKey, `secret`: captchaSecret})
	v := captcha.NewCaptchaVerifier(c.endpoint, captchaSecret)
	c.verifier = &captcha.SimpleCaptchaVerifier{
		Verifier: *v,
	}
	return nil
}

// keysValues: key1, value1, key2, value2
func (c *reCaptcha) Render(ctx echo.Context, templatePath string, keysValues ...interface{}) template.HTML {
	options := tplfunc.MakeMap(keysValues)
	options.Set("siteKey", c.siteKey)
	options.Set("endpoint", c.endpoint)
	if len(templatePath) == 0 {
		var htmlContent string
		switch c.endpoint {
		case captcha.CloudflareTurnstile:
			if !ctx.Internal().Bool(`CaptchaJSInited.CloudflareTurnstile`) {
				ctx.Internal().Set(`CaptchaJSInited.CloudflareTurnstile`, true)
				htmlContent += `<script src="https://challenges.cloudflare.com/turnstile/v0/api.js"></script>`
			}
			locationID := `turnstile-` + com.RandomAlphanumeric(16)
			htmlContent += `<div class="cf-turnstile" id="` + locationID + `" data-sitekey="` + c.siteKey + `"></div><input type="hidden" id="` + locationID + `-extend" disabled />`
			htmlContent += `<script>
			window.onload=function(){
				$('#` + locationID + `').closest('.input-group-addon').addClass('xxs-padding-top').prev('input').remove();
				turnstile.ready(function(){$('#` + locationID + `').data('lastGeneratedAt',(new Date()).getTime());});
				var $submit=$('#` + locationID + `').closest('form').find(':submit');
				$submit.on('click',function(e){
					if($('#` + locationID + `').data('lastGeneratedAt')>(new Date()).getTime()-290) {
						$('#` + locationID + `').data('lastGeneratedAt',0);
						return true;
					}
					window.setTimeout(function(){
						turnstile.reset('#` + locationID + `');
					},1000);
					$('#` + locationID + `').data('lastGeneratedAt',(new Date()).getTime());
				});
			}
		</script>`
		case HCaptcha:
			if !ctx.Internal().Bool(`CaptchaJSInited.HCaptcha`) {
				ctx.Internal().Set(`CaptchaJSInited.HCaptcha`, true)
				htmlContent += `<script src="https://js.hcaptcha.com/1/api.js?onload=" async defer></script>`
			}
			locationID := `hcaptcha-` + com.RandomAlphanumeric(16)
			htmlContent += `<div class="h-captcha" id="` + locationID + `" data-sitekey="` + c.siteKey + `"></div><input type="hidden" id="` + locationID + `-extend" disabled />`
			htmlContent += `<script>window.onload=function(){$('#` + locationID + `').closest('.input-group-addon').addClass('xxs-padding-top').prev('input').remove();}</script>`
		default:
			locationID := `recaptcha-` + com.RandomAlphanumeric(16)
			htmlContent = `<input type="hidden" id="` + locationID + `" name="g-recaptcha-response" value="" /><input type="hidden" id="` + locationID + `-extend" disabled />`
			if !ctx.Internal().Bool(`CaptchaJSInited.Recaptcha`) {
				ctx.Internal().Set(`CaptchaJSInited.Recaptcha`, true)
				htmlContent += `<script src="https://www.recaptcha.net/recaptcha/api.js?render=` + c.siteKey + `"></script>`
			}
			htmlContent += `<script>
			window.onload=function(){
				grecaptcha.ready(function() {
				  grecaptcha.execute('` + c.siteKey + `', {action: 'submit'}).then(function(token) {
					$('#` + locationID + `').val(token);
					$('#` + locationID + `').data('lastGeneratedAt',(new Date()).getTime());
				  });
				});
				$('#` + locationID + `').closest('.form-group').hide();
				$('#` + locationID + `').closest('.input-group').hide();
				$('#` + locationID + `').closest('.input-group-addon').prev('input').prop('disabled',true);
				var $submit=$('#` + locationID + `').closest('form').find(':submit');
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
					  $this.click();
					});
				});
			}
		</script>`
		}
		return template.HTML(htmlContent)
	}

	b, err := ctx.Fetch(templatePath, options)
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(com.Bytes2str(b))
}

func (c *reCaptcha) Verify(ctx echo.Context, hostAlias string, name string, captchaIdent ...string) echo.Data {
	var formKey string
	if len(captchaIdent) > 0 {
		formKey = captchaIdent[0]
	}
	if len(formKey) == 0 {
		switch c.endpoint {
		case captcha.CloudflareTurnstile:
			formKey = `cf-turnstile-response`
		case HCaptcha:
			formKey = `h-captcha-response`
		default:
			formKey = `g-recaptcha-response`
		}
	}
	token := ctx.Form(formKey)
	if len(token) == 0 { // 为空说明没有验证码
		return ctx.Data().SetError(ErrCaptchaCodeRequired.SetZone(name))
	}
	c.verifier.ExpectedHostname = ctx.Domain()
	if c.endpoint == captcha.GoogleRecaptcha {
		c.verifier.ExpectedAction = ctx.Form(`captchaAction`, `submit`)
		c.verifier.MinScore = 0.5
	}
	resp, ok, err := c.verifier.VerifyActionWithResponse(token, ``, c.verifier.ExpectedAction)
	if err != nil {
		return GenCaptchaErrorWithData(ctx, err, name, nil)
	}
	if !ok {
		log.Warnf(`failed to reCaptcha.Verify: %s`, formatter.AsStringer(resp))
		return GenCaptchaErrorWithData(ctx, ErrCaptcha, name, nil)
	}
	return ctx.Data().SetCode(code.Success.Int())
}

func (c *reCaptcha) MakeData(ctx echo.Context, hostAlias string, name string) echo.H {
	data := echo.H{}
	data.Set("siteKey", c.siteKey)
	data.Set("endpoint", c.endpoint)
	return data
}
