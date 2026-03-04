package file

import (
	"github.com/webx-top/echo"
)

func Finder(ctx echo.Context) error {
	return FinderWithOwner(ctx, ``, 0)
}
