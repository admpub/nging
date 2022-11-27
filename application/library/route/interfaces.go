package route

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
)

type IRegister interface {
	Echo() *echo.Echo
	Routes() []*echo.Route
	Logger() logger.Logger
	Prefix() string
	SetPrefix(prefix string)
	MetaHandler(m echo.H, handler interface{}, requests ...interface{}) echo.Handler
	MetaHandlerWithRequest(m echo.H, handler interface{}, request interface{}, methods ...string) echo.Handler
	HandlerWithRequest(handler interface{}, requests interface{}, methods ...string) echo.Handler
	AddGroupNamer(namers ...func(string) string)
	SetGroupNamer(namers ...func(string) string)
	SetRootGroup(groupName string)
	RootGroup() string
	Apply()
	Pre(middlewares ...interface{})
	Use(middlewares ...interface{})
	PreToGroup(groupName string, middlewares ...interface{})
	UseToGroup(groupName string, middlewares ...interface{})
	Register(fn func(echo.RouteRegister))
	RegisterToGroup(groupName string, fn func(echo.RouteRegister), middlewares ...interface{})
	Host(hostName string, middlewares ...interface{}) *Host
}
