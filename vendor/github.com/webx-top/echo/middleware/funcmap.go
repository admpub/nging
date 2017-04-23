package middleware

import (
	"html/template"

	"github.com/webx-top/echo"
)

func FuncMap(funcMap map[string]interface{}, skipper ...echo.Skipper) echo.MiddlewareFunc {
	var skip echo.Skipper
	if len(skipper) > 0 {
		skip = skipper[0]
	} else {
		skip = echo.DefaultSkipper
	}
	getFuncMap := FuncMapGetter(funcMap)
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if skip(c) {
				return h.Handle(c)
			}

			req := c.Request()
			c.SetFunc(`T`, c.T)
			c.SetFunc(`Lang`, c.Lang)
			c.SetFunc(`Stored`, c.Stored)
			c.SetFunc(`Cookie`, c.Cookie)
			c.SetFunc(`Session`, c.Session)
			c.SetFunc(`Form`, c.Form)
			c.SetFunc(`Query`, c.Query)
			c.SetFunc(`FormValues`, c.FormValues)
			c.SetFunc(`QueryValues`, c.QueryValues)
			c.SetFunc(`Param`, c.Param)
			c.SetFunc(`Atop`, c.Atop)
			c.SetFunc(`URL`, req.URL)
			c.SetFunc(`URI`, req.URI)
			c.SetFunc(`Referer`, c.Referer)
			c.SetFunc(`Header`, req.Header)
			c.SetFunc(`Flash`, c.Flash)
			c.SetFunc(`HasAnyRequest`, c.HasAnyRequest)
			for name, function := range c.Echo().FuncMap {
				c.SetFunc(name, function)
			}
			if getFuncMap != nil {
				for name, function := range getFuncMap(c) {
					c.SetFunc(name, function)
				}
			}
			return h.Handle(c)
		})
	}
}

func FuncMapGetter(funcMap interface{}) func(c echo.Context) map[string]interface{} {
	var getFuncMap func(c echo.Context) map[string]interface{}
	switch v := funcMap.(type) {
	case template.FuncMap:
		funcs := make(map[string]interface{})
		//copy value
		for k, f := range v {
			funcs[k] = f
		}
		getFuncMap = func(c echo.Context) map[string]interface{} {
			return funcs
		}
	case map[string]interface{}:
		funcs := make(map[string]interface{})
		//copy value
		for k, f := range v {
			funcs[k] = f
		}
		getFuncMap = func(c echo.Context) map[string]interface{} {
			return funcs
		}
	case func(echo.Context) map[string]interface{}:
		getFuncMap = v
	case func(echo.Context) template.FuncMap:
		getFuncMap = func(c echo.Context) map[string]interface{} {
			return v(c)
		}
	}
	return getFuncMap
}

func SimpleFuncMap(funcMap map[string]interface{}, skipper ...echo.Skipper) echo.MiddlewareFunc {
	var skip echo.Skipper
	if len(skipper) > 0 {
		skip = skipper[0]
	} else {
		skip = echo.DefaultSkipper
	}
	getFuncMap := FuncMapGetter(funcMap)
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if skip(c) {
				return h.Handle(c)
			}
			if getFuncMap != nil {
				for name, function := range getFuncMap(c) {
					c.SetFunc(name, function)
				}
			}
			return h.Handle(c)
		})
	}
}
