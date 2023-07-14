package embed

import "github.com/webx-top/echo"

type Config struct {
	Index    string
	Prefix   string
	FilePath func(echo.Context) (string, error)
}

var DefaultConfig = Config{
	Index: "index.html",
	FilePath: func(c echo.Context) (string, error) {
		return c.Param(`*`), nil
	},
}
