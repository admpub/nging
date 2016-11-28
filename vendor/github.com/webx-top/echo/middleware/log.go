package middleware

import (
	"fmt"
	"time"

	"github.com/webx-top/echo"
)

func Log() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			logger := c.Logger()

			start := time.Now()
			if err := h.Handle(c); err != nil {
				c.Error(err)
			}

			remoteAddr := req.RealIP()
			stop := time.Now()
			method := req.Method()
			uri := req.URI()
			size := res.Size()
			code := res.Status()
			logger.Info(remoteAddr + " " + method + " " + uri + " " + fmt.Sprint(code) + " " + stop.Sub(start).String() + " " + fmt.Sprint(size))
			return nil
		})
	}
}
