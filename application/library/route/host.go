package route

import (
	"github.com/webx-top/echo"
)

type Host struct {
	Name        string
	Alias       string
	Middlewares []interface{}
	Group       *Group
}

func (h *Host) Use(middlewares ...interface{}) {
	h.Middlewares = append(h.Middlewares, middlewares...)
}

func (h *Host) SetAlias(alias string) *Host {
	h.Alias = alias
	return h
}

func (h *Host) Register(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
	h.Group.Register(groupName, fn, middlewares...)
}
