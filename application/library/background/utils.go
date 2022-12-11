package background

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func Register(c echo.Context, op string, cacheKey string, bgExec *Background) (*Exec, error) {
	actual, _ := Backgrounds.LoadOrStore(op, &Exec{})
	exports := actual.(*Exec)
	if exports.Exists(cacheKey) {
		return exports, c.NewError(code.OperationProcessing, `任务正在后台处理中，请稍候...`)
	}
	exports.Add(op, cacheKey, bgExec)
	return exports, nil
}
