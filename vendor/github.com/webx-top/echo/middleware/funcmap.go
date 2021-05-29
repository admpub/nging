package middleware

import (
	"html/template"
	"strings"
	"time"

	"github.com/admpub/humanize"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"
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
			c.SetFunc(`Formx`, c.Formx)
			c.SetFunc(`Query`, c.Query)
			c.SetFunc(`Queryx`, c.Queryx)
			c.SetFunc(`FormValues`, c.FormValues)
			c.SetFunc(`QueryValues`, c.QueryValues)
			c.SetFunc(`FormxValues`, c.FormxValues)
			c.SetFunc(`GetByIndex`, param.GetByIndex)
			c.SetFunc(`QueryxValues`, c.QueryxValues)
			c.SetFunc(`Param`, c.Param)
			c.SetFunc(`Paramx`, c.Paramx)
			c.SetFunc(`Atop`, c.Atop)
			c.SetFunc(`URL`, req.URL)
			c.SetFunc(`URI`, req.URI)
			c.SetFunc(`Site`, c.Site)

			var pageURL string
			c.SetFunc(`SiteURI`, func() string {
				if len(pageURL) > 0 {
					return pageURL
				}
				pageURL = c.Site() + strings.TrimPrefix(req.URI(), `/`)
				return pageURL
			})
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
			c.SetFunc(`DurationFormat`, func(t interface{}, args ...string) *com.Durafmt {
				return tplfunc.DurationFormat(c.Lang(), t, args...)
			})
			c.SetFunc(`TsHumanize`, func(startTime interface{}, endTime ...interface{}) string {
				humanizer, err := humanize.New(c.Lang())
				if err != nil {
					return err.Error()
				}
				var (
					startDate = tplfunc.ToTime(startTime)
					endDate   time.Time
				)
				if len(endTime) > 0 {
					endDate = tplfunc.ToTime(endTime[0])
				}
				if endDate.IsZero() {
					endDate = time.Now().Local()
				}
				return humanizer.TimeDiff(endDate, startDate, 0)
			})
			c.SetFunc(`CaptchaForm`, func(args ...interface{}) template.HTML {
				return tplfunc.CaptchaFormWithURLPrefix(c.Echo().Prefix(), args...)
			})
			c.SetFunc(`MakeURL`, c.Echo().URL)
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
