package formbuilder

import (
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
)

type Option func(echo.Context, *FormBuilder)

// IgnoreFields 疏略某些字段的验证
func IgnoreFields(ignoreFields ...string) Option {
	return func(_ echo.Context, form *FormBuilder) {
		form.CloseValid(ignoreFields...)
	}
}

// ValidFields 指定仅仅验证某些字段
func ValidFields(validFields ...string) Option {
	return func(c echo.Context, form *FormBuilder) {
		c.Internal().Set(`formbuilder.validFields`, validFields)
	}
}

// Theme 设置forms模板风格
func Theme(theme string) Option {
	return func(_ echo.Context, form *FormBuilder) {
		form.Theme = theme
	}
}

// FormFilter 设置表单过滤
func FormFilter(filters ...formfilter.Options) Option {
	return func(c echo.Context, _ *FormBuilder) {
		c.Internal().Set(`formfilter.Options`, filters)
	}
}

// ConfigFile 指定要解析的配置文件。如果silent=true则仅仅设置配置文件而不解析
func ConfigFile(jsonFile string, silent ...bool) Option {
	return func(_ echo.Context, f *FormBuilder) {
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
	return func(_ echo.Context, f *FormBuilder) {
		f.AddBeforeRender(fn)
	}
}

// DBI 指定模型数据表所属DBI(数据库信息)
func DBI(dbi *factory.DBI) Option {
	return func(_ echo.Context, f *FormBuilder) {
		f.dbi = dbi
	}
}
