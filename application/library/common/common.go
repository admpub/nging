package common

import (
	stdErr "errors"

	"github.com/admpub/nging/application/library/errors"
	"github.com/webx-top/echo"
)

var PageMaxSize = 1000

func Paging(ctx echo.Context) (page int, size int) {
	page = ctx.Formx(`page`).Int()
	size = ctx.Formx(`size`).Int()
	if page < 1 {
		page = 1
	}
	if size < 1 || size > PageMaxSize {
		size = 50
	}
	return
}

func Ok(v string) errors.Successor {
	return errors.NewOk(v)
}

func Err(ctx echo.Context, err error) (ret interface{}) {
	if err == nil {
		flash := ctx.Flash()
		if flash != nil {
			if errMsg, ok := flash.(string); ok {
				ret = stdErr.New(errMsg)
			} else {
				ret = flash
			}
		}
	} else {
		ret = err
	}
	return
}
