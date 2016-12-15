/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package render

import (
	"net/http"

	"github.com/webx-top/echo"
)

var (
	DefaultDataKey       = `data`
	DefaultTmplKey       = `tmpl`
	DefaultTmplName      = `index`
	DefaultErrorTmpl     = `error`
	JSONPCallbackName    = `callback`
	DefaultErrorFunc     = OutputError
	DefaultOutputFunc    = Output
	DefaultErrorHTTPCode = http.StatusInternalServerError
)

// Middleware set renderer
func Middleware(d echo.Renderer) echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetRenderer(d)
			return h.Handle(c)
		})
	}
}

// AutoOutput Outputs the specified format
func AutoOutput(outputFunc func(string, echo.Context) error, skipper ...echo.Skipper) echo.MiddlewareFunc {
	isSkiped := echo.DefaultSkipper
	if len(skipper) > 0 {
		isSkiped = skipper[0]
	}
	if outputFunc == nil {
		outputFunc = DefaultOutputFunc
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if isSkiped(c) {
				return h.Handle(c)
			}
			format := c.Format()
			if err := h.Handle(c); err != nil {
				return DefaultErrorFunc(err, format, c)
			}
			return outputFunc(format, c)
		})
	}
}

// Output Outputs the specified format
func Output(format string, c echo.Context) error {
	switch format {
	case `json`:
		return c.JSON(c.Get(DefaultDataKey))
	case `jsonp`:
		return c.JSONP(c.Query(JSONPCallbackName), c.Get(DefaultDataKey))
	case `xml`:
		return c.XML(c.Get(DefaultDataKey))
	default:
		tmpl, ok := c.Get(DefaultTmplKey).(string)
		if !ok {
			tmpl = DefaultTmplName
		}
		data := c.Get(DefaultDataKey)
		if v, y := data.(*echo.Data); y {
			SetFuncs(c, v)
			return c.Render(tmpl, v.Data)
		}
		if h, y := data.(echo.H); y {
			v := h.ToData()
			SetFuncs(c, v)
			return c.Render(tmpl, v.Data)
		}
		return c.Render(tmpl, data)
	}
}

// SetFuncs register template function
func SetFuncs(c echo.Context, v *echo.Data) {
	c.SetFunc(`Info`, func() interface{} {
		return v.Info
	})
	c.SetFunc(`Code`, func() interface{} {
		return v.Code
	})
	c.SetFunc(`Zone`, func() interface{} {
		return v.Zone
	})
}

// OutputError Outputs the specified format
func OutputError(err error, format string, c echo.Context) error {
	if apiData, ok := err.(*echo.Data); ok {
		c.Set(DefaultDataKey, apiData)
	} else {
		c.Set(DefaultDataKey, echo.NewData(c.Code(), err.Error()))
	}
	c.Set(DefaultTmplKey, DefaultErrorTmpl)
	c.SetCode(DefaultErrorHTTPCode)
	return Output(format, c)
}

func HTTPErrorHandler(templates map[int]string, formatRender ...func(string, echo.Context) error) echo.HTTPErrorHandler {
	if templates == nil {
		templates = make(map[int]string)
	}
	tmplNum := len(templates)
	var output func(string, echo.Context) error
	if len(formatRender) > 0 {
		output = formatRender[0]
	}
	if output == nil {
		output = DefaultOutputFunc
	}
	return func(err error, c echo.Context) {
		code := DefaultErrorHTTPCode
		msg := http.StatusText(code)
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code > 0 {
				code = he.Code
			}
			msg = he.Message
		}
		if c.Echo().Debug() {
			msg = err.Error()
		}
		if !c.Response().Committed() {
			switch {
			case c.Request().Method() == echo.HEAD:
				c.NoContent(code)
			case tmplNum > 0:
				t, y := templates[code]
				if !y && tmplNum > 1 {
					t, y = templates[0]
				}
				if y {
					c.Set(DefaultDataKey, msg)
					c.Set(DefaultTmplKey, t)
					c.SetCode(code)
					if err := output(c.Format(), c); err != nil {
						msg += "\n" + err.Error()
						y = false
						c.Echo().Logger().Error(err)
					}
				}
				if y {
					break
				}
				fallthrough
			default:
				c.String(msg, code)
			}
		}
		c.Echo().Logger().Debug(err)
	}
}
