package handler

import (
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/errors"
	"github.com/webx-top/echo"
)

func Paging(ctx echo.Context) (page int, size int) {
	return common.Paging(ctx)
}

func Ok(v string) errors.Successor {
	return common.Ok(v)
}

func Err(ctx echo.Context, err error) (ret interface{}) {
	return common.Err(ctx, err)
}
