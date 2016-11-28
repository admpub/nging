package middleware

import (
	"github.com/webx-top/echo"
)

// MaxAllowed limits simultaneous requests; can help with high traffic load
func MaxAllowed(n int) echo.MiddlewareFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			acquire() // before request
			err := h.Handle(c)
			release() // after request
			return err
		})
	}
}
