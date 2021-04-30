package common

import (
	"errors"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/render"
)

var ErrorProcessors = []render.ErrorProcessor{
	func(ctx echo.Context, err error) (processed bool, newErr error) {
		if errors.Is(err, db.ErrNoMoreRows) {
			return true, echo.NewError(ctx.T(`数据不存在`), stdCode.DataNotFound)
		}
		return false, err
	},
}

func ProcessError(ctx echo.Context, err error) error {
	for _, processor := range ErrorProcessors {
		if processor == nil {
			continue
		}
		var processed bool
		processed, err = processor(ctx, err)
		if processed {
			break
		}
	}
	return err
}
