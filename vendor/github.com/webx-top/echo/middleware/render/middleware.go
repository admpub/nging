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
	DefaultOptions = &Options{
		Skipper:              echo.DefaultSkipper,
		DataKey:              `data`,
		TmplKey:              `tmpl`,
		DefaultTmpl:          `index`,
		DefaultErrorTmpl:     `error`,
		JSONPCallbackName:    `callback`,
		ErrorFunc:            OutputError,
		OutputFunc:           Output,
		DefaultErrorHTTPCode: http.StatusInternalServerError,
	}
)

type Options struct {
	Skipper              echo.Skipper
	DataKey              string
	TmplKey              string
	DefaultTmpl          string
	DefaultErrorTmpl     string
	JSONPCallbackName    string
	ErrorFunc            func(err error, format string, c echo.Context, opt *Options) error
	OutputFunc           func(format string, c echo.Context, opt *Options) error
	DefaultErrorHTTPCode int
}

// Middleware set renderer
func Middleware(d echo.Renderer) echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetRenderer(d)
			return h.Handle(c)
		})
	}
}

func SetDefaultOptions(opt *Options) *Options {
	if opt.Skipper == nil {
		opt.Skipper = DefaultOptions.Skipper
	}
	if opt.ErrorFunc == nil {
		opt.ErrorFunc = DefaultOptions.ErrorFunc
	}
	if opt.OutputFunc == nil {
		opt.OutputFunc = DefaultOptions.OutputFunc
	}
	if len(opt.DataKey) == 0 {
		opt.DataKey = DefaultOptions.DataKey
	}
	if len(opt.TmplKey) == 0 {
		opt.TmplKey = DefaultOptions.TmplKey
	}
	if len(opt.DefaultTmpl) == 0 {
		opt.DefaultTmpl = DefaultOptions.DefaultTmpl
	}
	if len(opt.DefaultErrorTmpl) == 0 {
		opt.DefaultErrorTmpl = DefaultOptions.DefaultErrorTmpl
	}
	if len(opt.JSONPCallbackName) == 0 {
		opt.JSONPCallbackName = DefaultOptions.JSONPCallbackName
	}
	return opt
}

func checkOptions(options ...*Options) *Options {
	var opt *Options
	if len(options) > 0 {
		opt = options[0]
	}
	if opt == nil {
		opt = DefaultOptions
	}
	return opt
}

// AutoOutput Outputs the specified format
func AutoOutput(options ...*Options) echo.MiddlewareFunc {
	opt := checkOptions(options...)
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if opt.Skipper(c) {
				return h.Handle(c)
			}
			format := c.Format()
			if err := h.Handle(c); err != nil {
				return opt.ErrorFunc(err, format, c, opt)
			}
			return opt.OutputFunc(format, c, opt)
		})
	}
}

// Output Outputs the specified format
func Output(format string, c echo.Context, opt *Options) error {
	switch format {
	case `json`:
		return c.JSON(c.Get(opt.DataKey))
	case `jsonp`:
		return c.JSONP(c.Query(opt.JSONPCallbackName), c.Get(opt.DataKey))
	case `xml`:
		return c.XML(c.Get(opt.DataKey))
	default:
		tmpl, ok := c.Get(opt.TmplKey).(string)
		if !ok {
			tmpl = opt.DefaultTmpl
		}
		data := c.Get(opt.DataKey)
		if v, y := data.(*echo.Data); y {
			SetFuncs(c, v)
			return c.Render(tmpl, v.Data)
		}
		if h, y := data.(echo.H); y {
			v := h.ToData().SetContext(c)
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
func OutputError(err error, format string, c echo.Context, opt *Options) error {
	if apiData, ok := err.(*echo.Data); ok {
		c.Set(opt.DataKey, apiData)
	} else {
		c.Set(opt.DataKey, echo.NewData(c, c.Code(), err.Error()))
	}
	c.Set(opt.TmplKey, opt.DefaultErrorTmpl)
	c.SetCode(opt.DefaultErrorHTTPCode)
	return Output(format, c, opt)
}

func HTTPErrorHandler(templates map[int]string, options ...*Options) echo.HTTPErrorHandler {
	if templates == nil {
		templates = make(map[int]string)
	}
	tmplNum := len(templates)
	opt := checkOptions(options...)
	return func(err error, c echo.Context) {
		code := opt.DefaultErrorHTTPCode
		msg := http.StatusText(code)
		title := msg
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
					c.Set(opt.DataKey, c.NewData().SetInfo(echo.H{
						"title":   title,
						"content": msg,
						"debug":   c.Echo().Debug(),
						"code":    code,
					}))
					c.Set(opt.TmplKey, t)
					c.SetCode(code)
					if err := opt.OutputFunc(c.Format(), c, opt); err != nil {
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
