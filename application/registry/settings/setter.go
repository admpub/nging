package settings

import "github.com/webx-top/echo"

type FormSetter func(*SettingForm)

func OptTmpl(tmpl ...string) FormSetter {
	return func(form *SettingForm) {
		form.Tmpl = tmpl
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
func OptHookGet(hookGet ...func(echo.Context) error) FormSetter {
	return func(form *SettingForm) {
		form.hookGet = hookGet
	}
}
