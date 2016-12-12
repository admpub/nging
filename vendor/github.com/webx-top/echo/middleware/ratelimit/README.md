## tollbooth_echo

[Echo](https://github.com/webx-top/echo) middleware for rate limiting HTTP requests.


## Five Minutes Tutorial

```
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
	limiter := ratelimit.NewLimiter(1, time.Second)

	e.Get("/", echo.HandlerFunc(func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	}), ratelimit.LimitHandler(limiter))

	e.Run(standard.New(":4444"))
}

```