package main

import (
	"os"
	"path/filepath"

	"github.com/webx-top/echo"
	// "github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/markdown"
)

func main() {
	e := echo.New()
	e.Use(mw.Log(), mw.Recover())
	e.Use(markdown.Markdown(&markdown.Options{
		Path:   "/book/",
		Root:   filepath.Join(os.Getenv(`GOPATH`), `src`, `github.com/admpub/gopl-zh`),
		Browse: true,
	}))

	e.Get("/", echo.HandlerFunc(func(c echo.Context) error {
		return c.String("Hello, World!")
	}))

	// FastHTTP
	// e.Run(fasthttp.New(":4444"))

	// Standard
	e.Run(standard.New(":4444"))
}
