package captcha

import (
	"html/template"
	"path"

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
	provider string
	endpoint captcha.Endpoint
	siteKey  string
	verifier *captcha.SimpleCaptchaVerifier
	jsURL    string
}

func (c *captchaAPI) Init(opt echo.H) error {
	c.provider = opt.String(`provider`)
	cfg := opt.GetStore(c.provider)
	c.siteKey = cfg.String(`siteKey`)
	switch c.provider {
	case `turnstile`:
		c.endpoint = captcha.CloudflareTurnstile
		c.jsURL = `https://challenges.cloudflare.com/turnstile/v0/api.js`
	default:
		c.endpoint = captcha.GoogleRecaptcha
		c.jsURL = `https://www.recaptcha.net/recaptcha/api.js?render=` + c.siteKey
	}
	captchaSecret := cfg.String(`secret`)
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
	uniqid := com.RandomAlphanumeric(16)
	options.Set("uniqid", uniqid)
	if len(templatePath) == 0 {
		var htmlContent string
		switch c.endpoint {
		case captcha.CloudflareTurnstile:
			if len(jsURL) > 0 {
				htmlContent += `<script src="` + jsURL + `"></script>`
			}
			locationID := `turnstile-` + uniqid
			htmlContent += `<div class="cf-turnstile" id="` + locationID + `" data-sitekey="` + c.siteKey + `"></div><input type="hidden" id="` + locationID + `-extend" disabled />`
			htmlContent += `<script>
			var windowOriginalOnload` + uniqid + `=window.onload;
			window.onload=function(){
				$('#` + locationID + `').closest('.input-group-addon').addClass('xxs-padding-top').prev('input').remove();
				turnstile.ready(function(){$('#` + locationID + `').data('lastGeneratedAt',(new Date()).getTime());});
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
				windowOriginalOnload` + uniqid + ` && windowOriginalOnload` + uniqid + `.apply(this,arguments);
			}
		</script>`
		default:
			locationID := `recaptcha-` + uniqid
			htmlContent = `<input type="hidden" id="` + locationID + `" name="g-recaptcha-response" value="" /><input type="hidden" id="` + locationID + `-extend" disabled />`
			if len(jsURL) > 0 {
				htmlContent += `<script src="` + jsURL + `"></script>`
			}
			htmlContent += `<script>
			var windowOriginalOnload` + uniqid + `=window.onload;
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
					  $this.trigger('click');
					});
				});
				windowOriginalOnload` + uniqid + ` && windowOriginalOnload` + uniqid + `.apply(this,arguments);
			}
		</script>`
		}
		return template.HTML(htmlContent)
	}

	b, err := ctx.Fetch(path.Join(`captcha/api`, templatePath), options)
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(com.Bytes2str(b))
}

func (c *captchaAPI) Verify(ctx echo.Context, hostAlias string, name string, captchaIdent ...string) echo.Data {
	var formKey string
	if len(captchaIdent) > 0 {
		formKey = captchaIdent[0]
	}
	if len(formKey) == 0 {
		switch c.endpoint {
		case captcha.CloudflareTurnstile:
			formKey = `cf-turnstile-response`
		default:
			formKey = `g-recaptcha-response`
		}
	}
	token := ctx.Form(formKey)
	if len(token) == 0 { // 为空说明没有验证码
		return ctx.Data().SetError(ErrCaptchaCodeRequired.SetMessage(ctx.T(`请先进行人机验证`)).SetZone(name))
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
		log.Warnf(`failed to captchaAPI.Verify: %s`, formatter.AsStringer(resp))
		return GenCaptchaErrorWithData(ctx, ErrCaptcha.SetMessage(ctx.T(`抱歉，未能通过人机验证`)), name, nil)
	}
	return ctx.Data().SetCode(code.Success.Int())
}

func (c *captchaAPI) MakeData(ctx echo.Context, hostAlias string, name string) echo.H {
	data := echo.H{}
	data.Set("siteKey", c.siteKey)
	data.Set("endpoint", c.endpoint)
	return data
}
