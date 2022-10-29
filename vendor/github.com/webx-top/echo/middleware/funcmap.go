package middleware

import (
	"github.com/webx-top/echo"
)

func FuncMap(funcMap map[string]interface{}, skipper ...echo.Skipper) echo.MiddlewareFunc {
	var skip echo.Skipper
	if len(skipper) > 0 {
		skip = skipper[0]
	} else {
		skip = echo.DefaultSkipper
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if skip(c) {
				return h.Handle(c)
			}

			for name, function := range funcMap {
				c.SetFunc(name, function)
			}
			return h.Handle(c)
		})
	}
}
