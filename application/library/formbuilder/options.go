package formbuilder

import (
	"github.com/coscms/forms"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
)

type Option func(echo.Context, *forms.Form)

func IgnoreFields(ignoreFields ...string) Option {
	return func(_ echo.Context, form *forms.Form) {
		form.CloseValid(ignoreFields...)
	}
}

func ValidFields(validFields ...string) Option {
	return func(c echo.Context, form *forms.Form) {
		c.Internal().Set(`formbuilder.validFields`, validFields)
	}
}

func Style(style string) Option {
	return func(_ echo.Context, form *forms.Form) {
		form.Style = style
	}
}

func FormFilter(filters ...formfilter.Options) Option {
	return func(c echo.Context, _ *forms.Form) {
		c.Internal().Set(`formfilter.Options`, filters)
	}
}
