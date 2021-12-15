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

package echo

import (
	"fmt"
	"net/http"
)

// WrapHandler wrap `interface{}` into `echo.Handler`.
func WrapHandler(h interface{}) Handler {
	switch v := h.(type) {
	case HandlerFunc:
		return v
	case Handler:
		return v
	case func(Context) error:
		return HandlerFunc(v)
	case http.Handler:
		return HandlerFunc(func(ctx Context) error {
			v.ServeHTTP(
				ctx.Response().StdResponseWriter(),
				ctx.Request().StdRequest().WithContext(ctx),
			)
			return nil
		})
	case func(http.ResponseWriter, *http.Request):
		return HandlerFunc(func(ctx Context) error {
			v(
				ctx.Response().StdResponseWriter(),
				ctx.Request().StdRequest().WithContext(ctx),
			)
			return nil
		})
	case func(http.ResponseWriter, *http.Request) error:
		return HandlerFunc(func(ctx Context) error {
			return v(
				ctx.Response().StdResponseWriter(),
				ctx.Request().StdRequest().WithContext(ctx),
			)
		})

	// lazyload
	case func() HandlerFunc:
		return v()
	case func() func(Context) error:
		return HandlerFunc(v())

	default:
		panic(fmt.Sprintf(`unknown handler: %T`, h))
	}
}

// WrapMiddleware wrap `interface{}` into `echo.Middleware`.
func WrapMiddleware(m interface{}) Middleware {
	switch h := m.(type) {
	case MiddlewareFunc:
		return h
	case MiddlewareFuncd:
		return h
	case Middleware:
		return h
	case HandlerFunc:
		return WrapMiddlewareFromHandler(h)
	case func(Context) error:
		return WrapMiddlewareFromHandler(HandlerFunc(h))
	case func(Handler) func(Context) error:
		return MiddlewareFunc(func(next Handler) Handler {
			return HandlerFunc(h(next))
		})
	case func(Handler) HandlerFunc:
		return MiddlewareFunc(func(next Handler) Handler {
			return h(next)
		})
	case func(HandlerFunc) HandlerFunc:
		return MiddlewareFunc(func(next Handler) Handler {
			return h(next.Handle)
		})
	case func(Handler) Handler:
		return MiddlewareFunc(h)
	case func(func(Context) error) func(Context) error:
		return MiddlewareFunc(func(next Handler) Handler {
			return HandlerFunc(h(next.Handle))
		})
	case http.Handler:
		return WrapMiddlewareFromStdHandler(h)
	case func(http.ResponseWriter, *http.Request):
		return WrapMiddlewareFromStdHandleFunc(h)
	case func(http.ResponseWriter, *http.Request) error:
		return WrapMiddlewareFromStdHandleFuncd(h)

	// lazyload
	case func() MiddlewareFunc:
		return h()
	case func() MiddlewareFuncd:
		return h()
	case func() HandlerFunc:
		return WrapMiddlewareFromHandler(h())
	case func() func(Context) error:
		return WrapMiddlewareFromHandler(HandlerFunc(h()))

	default:
		panic(fmt.Sprintf(`unknown middleware: %T`, m))
	}
}

// WrapMiddlewareFromHandler wrap `echo.HandlerFunc` into `echo.Middleware`.
func WrapMiddlewareFromHandler(h HandlerFunc) Middleware {
	return MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(c Context) error {
			if err := h.Handle(c); err != nil {
				return err
			}
			return next.Handle(c)
		})
	})
}

// WrapMiddlewareFromStdHandler wrap `http.HandlerFunc` into `echo.Middleware`.
func WrapMiddlewareFromStdHandler(h http.Handler) Middleware {
	return MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(c Context) error {
			h.ServeHTTP(
				c.Response().StdResponseWriter(),
				c.Request().StdRequest().WithContext(c),
			)
			if c.Response().Committed() {
				return nil
			}
			return next.Handle(c)
		})
	})
}

// WrapMiddlewareFromStdHandleFunc wrap `func(http.ResponseWriter, *http.Request)` into `echo.Middleware`.
func WrapMiddlewareFromStdHandleFunc(h func(http.ResponseWriter, *http.Request)) Middleware {
	return MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(c Context) error {
			h(
				c.Response().StdResponseWriter(),
				c.Request().StdRequest().WithContext(c),
			)
			if c.Response().Committed() {
				return nil
			}
			return next.Handle(c)
		})
	})
}

// WrapMiddlewareFromStdHandleFuncd wrap `func(http.ResponseWriter, *http.Request)` into `echo.Middleware`.
func WrapMiddlewareFromStdHandleFuncd(h func(http.ResponseWriter, *http.Request) error) Middleware {
	return MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(c Context) error {
			if err := h(
				c.Response().StdResponseWriter(),
				c.Request().StdRequest().WithContext(c),
			); err != nil {
				return err
			}
			if c.Response().Committed() {
				return nil
			}
			return next.Handle(c)
		})
	})
}
