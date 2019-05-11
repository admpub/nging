package index

import (
	"io/ioutil"

	"github.com/webx-top/echo"
)

func Ping(ctx echo.Context) error {
	header := ctx.Request().Header()
	body := ctx.Request().Body()
	b, _ := ioutil.ReadAll(body)
	body.Close()
	r := echo.H{
		`header`: header.Object(),
		`form`:   echo.NewMapx(ctx.Request().Form().All()).AsStore(),
		`body`:   string(b),
	}
	data := ctx.Data()
	data.SetData(r)
	callback := ctx.Form(`callback`)
	if len(callback) > 0 {
		return ctx.JSONP(callback, data)
	}
	return ctx.JSON(data)
}
