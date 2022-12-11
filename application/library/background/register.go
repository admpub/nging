package background

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func Register(c echo.Context, op string, cacheKey string, bg *Background) (*Group, error) {
	actual, _ := Backgrounds.LoadOrStore(op, &Group{})
	group := actual.(*Group)
	if group.Exists(cacheKey) {
		return group, c.NewError(code.OperationProcessing, `任务正在后台处理中，请稍候...`)
	}
	group.Add(op, cacheKey, bg)
	return group, nil
}
