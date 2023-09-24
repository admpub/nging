package formbuilder

import (
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo/formfilter"
)

type Option func(*FormBuilder)

// IgnoreFields 疏略某些字段的验证
func IgnoreFields(ignoreFields ...string) Option {
	return func(form *FormBuilder) {
		form.CloseValid(ignoreFields...)
	}
}

// Theme 设置forms模板风格
func Theme(theme string) Option {
	return func(form *FormBuilder) {
		form.Theme = theme
	}
}

// FormFilter 设置表单过滤
func FormFilter(filters ...formfilter.Options) Option {
	return func(form *FormBuilder) {
		form.filters = filters
	}
}

// ConfigFile 指定要解析的配置文件。如果silent=true则仅仅设置配置文件而不解析
func ConfigFile(jsonFile string, silent ...bool) Option {
	return func(f *FormBuilder) {
		f.configFile = jsonFile
		if len(silent) > 0 && silent[0] {
			return
		}
		if err := f.ParseConfigFile(); err != nil {
			panic(err)
		}
	}
}

// RenderBefore 设置渲染表单前的钩子函数
func RenderBefore(fn func()) Option {
	return func(f *FormBuilder) {
		f.AddBeforeRender(fn)
	}
}

// DBI 指定模型数据表所属DBI(数据库信息)
func DBI(dbi *factory.DBI) Option {
	return func(f *FormBuilder) {
		f.dbi = dbi
	}
}
