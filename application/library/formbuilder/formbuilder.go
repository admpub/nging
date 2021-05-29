package formbuilder

import (
	"os"

	ncommon "github.com/admpub/nging/application/library/common"

	"github.com/coscms/forms"
	"github.com/coscms/forms/common"
	"github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/validation"
)

// New 表单
//@param m: dbschema
func New(c echo.Context, m interface{}, jsonFile string, options ...Option) *forms.Forms {
	form := forms.New()
	form.Style = common.BOOTSTRAP
	for _, option := range options {
		if option == nil {
			continue
		}
		option(c, form)
	}
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
		opts := []formfilter.Options{formfilter.Include(cfg.GetNames()...)}
		if customs, ok := c.Internal().Get(`formfilter.Options`).([]formfilter.Options); ok {
			opts = append(opts, customs...)
		}
		err = c.MustBind(m, formfilter.Build(opts...))
		if err == nil {
			form.ValidFromConfig()
			valid := form.Validate()
			if valid.HasError() {
				err = valid.Errors[0]
			}
		}
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
	wrap := forms.NewForms(form)
	c.Set(`forms`, wrap)
	// 手动调用:
	// wrap.ParseFromConfig()
	return wrap
}

// NewModel 表单
//@param m: dbschema
func NewModel(c echo.Context, m interface{}, cfg *config.Config, options ...Option) *forms.Forms {
	form := forms.New()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(c, form)
	}
	form.SetLabelFunc(func(txt string) string {
		return c.T(txt)
	})
	if cfg == nil {
		cfg = form.NewConfig()
	}
	form.Init(cfg, m)
	if c.IsPost() {
		opts := []formfilter.Options{formfilter.Include(cfg.GetNames()...)}
		if customs, ok := c.Internal().Get(`formfilter.Options`).([]formfilter.Options); ok {
			opts = append(opts, customs...)
		}
		err := c.MustBind(m, formfilter.Build(opts...))
		if err == nil {
			validFields, _ := c.Internal().Get(`formbuilder.validFields`).([]string)
			err = form.Valid(validFields...)
		}
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

func AddChoiceByKV(field fields.FieldInterface, kvData *echo.KVData, checkedKeys ...string) fields.FieldInterface {
	for _, kv := range kvData.Slice() {
		var checked bool
		if kv.H != nil {
			checked = kv.H.Bool(`checked`)
		}
		if len(checkedKeys) > 0 {
			checked = com.InSlice(kv.K, checkedKeys)
		}
		field.AddChoice(kv.K, kv.V, checked)
	}
	return field
}
