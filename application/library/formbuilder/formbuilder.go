package formbuilder

import (
	"os"

	ncommon "github.com/admpub/nging/application/library/common"

	"github.com/coscms/forms"
	"github.com/coscms/forms/common"
	"github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/validation"
)

// New 表单
//@param m: dbschema
func New(c echo.Context, m interface{}, jsonFile string, ingoreFields ...string) *forms.Forms {
	form := forms.New()
	form.Style = common.BOOTSTRAP
	form.SetLabelFunc(func(txt string) string {
		return c.T(txt)
	})
	var cfg *config.Config
	renderer := c.Renderer().(driver.Driver)
	jsonFile += `.form.json`
	jsonFile = renderer.TmplPath(c, jsonFile)
	if len(jsonFile) == 0 {
		return nil
	}
	b, err := renderer.RawContent(jsonFile)
	if err != nil {
		if os.IsNotExist(err) && renderer.Manager() != nil {
			form.SetModel(m)
			cfg = form.ToConfig()
			var jsonb []byte
			jsonb, err = form.ToJSONBlob(cfg)
			if err == nil {
				err = renderer.Manager().SetTemplate(jsonFile, jsonb)
				if err == nil {
					c.Logger().Infof(c.T(`生成表单配置文件“%v”成功。`), jsonFile)
				}
			}
		}
	} else {
		cfg, err = forms.Unmarshal(b, jsonFile)
	}
	if err != nil {
		c.Logger().Error(err)
	}
	if cfg == nil {
		cfg = form.NewConfig()
	}
	form.Init(cfg, m)
	if c.IsPost() {
		form.CloseValid(ingoreFields...)
		form.ValidFromConfig()
		valid := form.Validate()
		if valid.HasError() {
			c.Data().SetInfo(valid.Errors[0].Message, 0).SetZone(valid.Errors[0].Field)
		}
	}
	setNextURLField := func() {
		if len(cfg.Action) == 0 {
			form.SetParam(`action`, c.Request().URI())
		}
		nextURL := c.Form(ncommon.DefaultReturnToURLVarName)
		if len(nextURL) == 0 {
			nextURL = c.Referer()
		}
		form.Elements(fields.HiddenField(ncommon.DefaultReturnToURLVarName).SetValue(nextURL))
	}
	csrfToken, ok := c.Get(`csrf`).(string)
	if ok {
		form.AddBeforeRender(func() {
			form.Elements(fields.HiddenField(`csrf`).SetValue(csrfToken))
			setNextURLField()
		})
	} else {
		form.AddBeforeRender(setNextURLField)
	}
	wrap := forms.NewForms(form)
	c.Set(`forms`, wrap)
	// 手动调用:
	// wrap.ParseFromConfig()
	return wrap
}

// NewModel 表单
//@param m: dbschema
func NewModel(c echo.Context, m interface{}, cfg *config.Config, validFields ...string) *forms.Forms {
	form := forms.New()
	form.SetLabelFunc(func(txt string) string {
		return c.T(txt)
	})
	if cfg == nil {
		cfg = form.NewConfig()
	}
	form.Init(cfg, m)
	if c.IsPost() {
		err := form.Valid(validFields...)
		if err != nil {
			if vErr, ok := err.(*validation.ValidationError); ok {
				c.Data().SetInfo(vErr.Message, 0).SetZone(vErr.Field)
			} else {
				c.Data().SetError(err)
			}
		}
	}
	setNextURLField := func() {
		if len(cfg.Action) == 0 {
			form.SetParam(`action`, c.Request().URI())
		}
		nextURL := c.Form(ncommon.DefaultReturnToURLVarName)
		if len(nextURL) == 0 {
			nextURL = c.Referer()
		}
		form.Elements(fields.HiddenField(ncommon.DefaultReturnToURLVarName).SetValue(nextURL))
	}
	csrfToken, ok := c.Get(`csrf`).(string)
	if ok {
		form.AddBeforeRender(func() {
			form.Elements(fields.HiddenField(`csrf`).SetValue(csrfToken))
			setNextURLField()
		})
	} else {
		form.AddBeforeRender(setNextURLField)
	}
	form.AddClass("form-horizontal").SetParam("role", "form")
	wrap := forms.NewForms(form)
	c.Set(`forms`, wrap)
	// 手动调用:
	// wrap.ParseFromConfig()
	return wrap
}

// NewConfig 表单配置
func NewConfig(theme, tmpl, method, action string) *config.Config {
	cfg := forms.NewConfig()
	cfg.Theme = theme
	cfg.Template = tmpl
	cfg.Method = method
	cfg.Action = action
	return cfg
}

// NewSnippet 表单片段
func NewSnippet(theme ...string) *forms.Form {
	cfg := forms.NewConfig()
	if len(theme) > 0 {
		cfg.Theme = theme[0]
	}
	cfg.Template = common.TmplDir(cfg.Theme) + `/allfields.html`
	form := forms.NewWithConfig(cfg)
	return form
}

func ClearCache() {
	common.ClearCachedConfig()
	common.ClearCachedTemplate()
}

func DelCachedConfig(file string) bool {
	return common.DelCachedConfig(file)
}
