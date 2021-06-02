package formbuilder

import (
	"github.com/coscms/forms"
	"github.com/coscms/forms/common"
	"github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

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
			checked = kv.H.Bool(`checked`) || kv.H.Bool(`selected`)
		}
		if len(checkedKeys) > 0 {
			checked = com.InSlice(kv.K, checkedKeys)
		}
		field.AddChoice(kv.K, kv.V, checked)
	}
	return field
}

func SetChoiceByKV(field fields.FieldInterface, kvData *echo.KVData, checkedKeys ...string) fields.FieldInterface {
	choices := []fields.InputChoice{}
	for _, kv := range kvData.Slice() {
		var checked bool
		if kv.H != nil {
			checked = kv.H.Bool(`checked`) || kv.H.Bool(`selected`)
		}
		if len(checkedKeys) > 0 {
			checked = com.InSlice(kv.K, checkedKeys)
		}
		choices = append(choices, fields.InputChoice{
			ID:      kv.K,
			Val:     kv.V,
			Checked: checked,
		})
	}

	field.SetChoices(choices)
	return field
}
