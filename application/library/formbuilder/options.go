package formbuilder

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
)

type Option func(echo.Context, *FormBuilder)

func IgnoreFields(ignoreFields ...string) Option {
	return func(_ echo.Context, form *FormBuilder) {
		form.CloseValid(ignoreFields...)
	}
}

func ValidFields(validFields ...string) Option {
	return func(c echo.Context, form *FormBuilder) {
		c.Internal().Set(`formbuilder.validFields`, validFields)
	}
}

func Style(style string) Option {
	return func(_ echo.Context, form *FormBuilder) {
		form.Style = style
	}
}

func FormFilter(filters ...formfilter.Options) Option {
	return func(c echo.Context, _ *FormBuilder) {
		c.Internal().Set(`formfilter.Options`, filters)
	}
}

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

func RenderBefore(fn func()) Option {
	return func(_ echo.Context, f *FormBuilder) {
		f.AddBeforeRender(fn)
	}
}
