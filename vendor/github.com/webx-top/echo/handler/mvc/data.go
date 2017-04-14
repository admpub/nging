package mvc

import (
	"github.com/webx-top/echo"
)

type Data interface {
	Assign(key string, val interface{})
	Assignx(values *map[string]interface{})
	SetTmplFuncs()
	Render(tmpl string, code ...int) error
	String() string
	Set(code int, args ...interface{})
	Gets() (code echo.State, info interface{}, zone interface{}, data interface{})
	GetData() interface{}
}
