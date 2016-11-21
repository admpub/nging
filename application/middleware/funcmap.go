package middleware

import (
	"time"

	"github.com/webx-top/echo"
)

func FuncMap() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetFunc(`Now`, time.Now)
			return h.Handle(c)
		})
	}
}
