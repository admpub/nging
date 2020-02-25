package route

import (
	"github.com/webx-top/echo"
)

type Host struct {
	Name        string
	Middlewares []interface{}
	Group       *Group
}

func (h *Host) Use(middlewares ...interface{}) {
	h.Middlewares = append(h.Middlewares, middlewares...)
}

func (h *Host) Register(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
	h.Group.Register(groupName, fn, middlewares...)
}
