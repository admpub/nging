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
	"time"

	"github.com/webx-top/echo"
)

var (
	DefaultOptions = &Options{
		Skipper:              echo.DefaultSkipper,
		ErrorPages:           make(map[int]string),
		DefaultHTTPErrorCode: http.StatusInternalServerError,
		SetFuncMap: []echo.HandlerFunc{
			func(c echo.Context) error {
				c.SetFunc(`Lang`, c.Lang)
				c.SetFunc(`Now`, time.Now)
				c.SetFunc(`T`, c.T)
				return nil
			},
		},
	}
)

type Options struct {
	Skipper              echo.Skipper
	ErrorPages           map[int]string
	DefaultHTTPErrorCode int
	SetFuncMap           []echo.HandlerFunc
}

func (opt *Options) AddFuncSetter(set ...echo.HandlerFunc) *Options {
	if opt.SetFuncMap == nil {
		opt.SetFuncMap = make([]echo.HandlerFunc, len(DefaultOptions.SetFuncMap))
		for index, setter := range DefaultOptions.SetFuncMap {
			opt.SetFuncMap[index] = setter
		}
	}
	opt.SetFuncMap = append(opt.SetFuncMap, set...)
	return opt
}

func (opt *Options) SetFuncSetter(set ...echo.HandlerFunc) *Options {
	opt.SetFuncMap = set
	return opt
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

func Auto() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetAuto(true)
			return h.Handle(c)
		})
	}
}

func HTTPErrorHandler(opt *Options) echo.HTTPErrorHandler {
	if opt == nil {
		opt = DefaultOptions
	}
	if opt.ErrorPages == nil {
		opt.ErrorPages = DefaultOptions.ErrorPages
	}
	if opt.DefaultHTTPErrorCode < 1 {
		opt.DefaultHTTPErrorCode = DefaultOptions.DefaultHTTPErrorCode
	}
	if opt.SetFuncMap == nil {
		opt.SetFuncMap = DefaultOptions.SetFuncMap
	}
	tmplNum := len(opt.ErrorPages)
	return func(err error, c echo.Context) {
		code := DefaultOptions.DefaultHTTPErrorCode
		var msg string
		var panicErr *echo.PanicError
		switch e := err.(type) {
		case *echo.HTTPError:
			if e.Code > 0 {
				code = e.Code
			}
			msg = e.Message
		case *echo.PanicError:
			panicErr = e

		}
		title := http.StatusText(code)
		if c.Echo().Debug() {
			msg = err.Error()
		} else if len(msg) == 0 {
			msg = title
		}
		if !c.Response().Committed() {
			switch {
			case c.Request().Method() == echo.HEAD:
				c.NoContent(code)
			case tmplNum > 0:
				t, y := opt.ErrorPages[code]
				if !y && code != 0 {
					t, y = opt.ErrorPages[0]
				}
				if y {
					data := c.Data().Reset().SetInfo(msg, 0)
					if c.Format() == `html` {
						c.SetCode(code)
						c.SetFunc(`Lang`, c.Lang)
						if len(opt.SetFuncMap) > 0 {
							for _, setFunc := range opt.SetFuncMap {
								err = setFunc(c)
								if err != nil {
									c.String(err.Error())
									return
								}
							}
						}
						data.SetData(echo.H{
							"title":   title,
							"content": msg,
							"debug":   c.Echo().Debug(),
							"code":    code,
							"panic":   panicErr,
						}, 0)
					} else {
						c.SetCode(opt.DefaultHTTPErrorCode)
					}
					if err := c.SetAuto(true).Render(t, nil); err != nil {
						msg += "\n" + err.Error()
						y = false
						c.Logger().Error(err)
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
		c.Logger().Debug(err)
	}
}
