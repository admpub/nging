package settings

import (
	"html/template"

	"github.com/admpub/nging/v3/application/dbschema"
	formsconfig "github.com/coscms/forms/config"
	"github.com/webx-top/echo"
)

type FormSetter func(*SettingForm)

func OptTmpl(tmpl ...string) FormSetter {
	return func(form *SettingForm) {
		form.Tmpl = tmpl
	}
}
func OptAddTmpl(tmpl ...string) FormSetter {
	return func(form *SettingForm) {
		form.Tmpl = append(form.Tmpl, tmpl...)
	}
}
func OptHeadTmpl(tmpl ...string) FormSetter {
	return func(form *SettingForm) {
		form.HeadTmpl = tmpl
	}
}
func OptAddHeadTmpl(tmpl ...string) FormSetter {
	return func(form *SettingForm) {
		form.HeadTmpl = append(form.HeadTmpl, tmpl...)
	}
}
func OptFootTmpl(tmpl ...string) FormSetter {
	return func(form *SettingForm) {
		form.FootTmpl = tmpl
	}
}
func OptAddFootTmpl(tmpl ...string) FormSetter {
	return func(form *SettingForm) {
		form.FootTmpl = append(form.FootTmpl, tmpl...)
	}
}
func OptShort(short string) FormSetter {
	return func(form *SettingForm) {
		form.Short = short
	}
}
func OptLabel(label string) FormSetter {
	return func(form *SettingForm) {
		form.Label = label
	}
}
func OptGroup(group string) FormSetter {
	return func(form *SettingForm) {
		form.Group = group
	}
}
func OptHookPost(hookPost ...func(echo.Context) error) FormSetter {
	return func(form *SettingForm) {
		form.hookPost = hookPost
	}
}
func OptAddHookPost(hookPost ...func(echo.Context) error) FormSetter {
	return func(form *SettingForm) {
		form.hookPost = append(form.hookPost, hookPost...)
	}
}
func OptHookGet(hookGet ...func(echo.Context) error) FormSetter {
	return func(form *SettingForm) {
		form.hookGet = hookGet
	}
}
func OptAddHookGet(hookGet ...func(echo.Context) error) FormSetter {
	return func(form *SettingForm) {
		form.hookGet = append(form.hookGet, hookGet...)
	}
}
func OptFormConfig(formcfg *formsconfig.Config) FormSetter {
	return func(form *SettingForm) {
		form.SetFormConfig(formcfg)
	}
}
func OptDataTransfer(name string, dataInitor DataInitor, dataFrom DataFrom) FormSetter {
	return func(form *SettingForm) {
		form.SetDataTransfer(name, dataInitor, dataFrom)
	}
}
func OptAddConfig(configs ...*dbschema.NgingConfig) FormSetter {
	return func(form *SettingForm) {
		form.AddConfig(configs...)
	}
}
func OptRenderer(renderer func(echo.Context) template.HTML) FormSetter {
	return func(form *SettingForm) {
		form.SetRenderer(renderer)
	}
}
