package route

import (
	"github.com/webx-top/echo"
)

type Hoster interface {
	Use(middlewares ...interface{})
	SetAlias(alias string) Hoster
	Register(groupName string, fn func(echo.RouteRegister), middlewares ...interface{})
}

type Host struct {
	Name        string
	Alias       string
	Middlewares []interface{}
	Group       *Group
}

func (h *Host) Use(middlewares ...interface{}) {
	h.Middlewares = append(h.Middlewares, middlewares...)
}

func (h *Host) SetAlias(alias string) Hoster {
	h.Alias = alias
	return h
}

func (h *Host) Register(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
	h.Group.Register(groupName, fn, middlewares...)
}

type noopHost struct {
}

func (h *noopHost) Use(middlewares ...interface{}) {
}

func (h *noopHost) SetAlias(alias string) Hoster {
	return h
}

func (h *noopHost) Register(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
}
