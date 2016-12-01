package main

import (
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/middleware/ratelimit"
)

func main() {
	e := echo.New()

	// Create a limiter struct.
	limiter := ratelimit.New(1, time.Second)

	e.Get("/", echo.HandlerFunc(func(c echo.Context) error {
		return c.String("Hello, World!")
	}), ratelimit.LimitHandler(limiter))

	e.Run(standard.New(":4444"))
}
