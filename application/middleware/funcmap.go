package middleware

import (
	"time"

	"strings"

	"github.com/webx-top/echo"
)

func FuncMap() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetFunc(`Now`, time.Now)
			c.SetFunc(`Add`, func(x int, y int) int {
				return x + y
			})
			c.SetFunc(`Sub`, func(x int, y int) int {
				return x - y
			})
			c.SetFunc(`HasPrefix`, strings.HasPrefix)
			c.SetFunc(`HasSuffix`, strings.HasSuffix)
			return h.Handle(c)
		})
	}
}
