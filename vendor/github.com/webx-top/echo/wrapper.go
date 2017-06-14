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

import "net/http"

// WrapHandler wrap `interface{}` into `echo.Handler`.
func WrapHandler(h interface{}) Handler {
	if v, ok := h.(HandlerFunc); ok {
		return v
	}
	if v, ok := h.(Handler); ok {
		return v
	}
	if v, ok := h.(func(Context) error); ok {
		return HandlerFunc(v)
	}
	if v, ok := h.(http.Handler); ok {
		return HandlerFunc(func(ctx Context) error {
			v.ServeHTTP(ctx.Response().StdResponseWriter(), ctx.Request().StdRequest())
			return nil
		})
	}
	if v, ok := h.(func(http.ResponseWriter, *http.Request)); ok {
		return HandlerFunc(func(ctx Context) error {
			v(ctx.Response().StdResponseWriter(), ctx.Request().StdRequest())
			return nil
		})
	}
	if v, ok := h.(func(http.ResponseWriter, *http.Request) error); ok {
		return HandlerFunc(func(ctx Context) error {
			return v(ctx.Response().StdResponseWriter(), ctx.Request().StdRequest())
		})
	}
	panic(`unknown handler`)
}

// WrapMiddleware wrap `interface{}` into `echo.Middleware`.
func WrapMiddleware(m interface{}) Middleware {
	if h, ok := m.(MiddlewareFunc); ok {
		return h
	}
	if h, ok := m.(MiddlewareFuncd); ok {
		return h
	}
	if h, ok := m.(Middleware); ok {
		return h
	}
	if h, ok := m.(HandlerFunc); ok {
		return WrapMiddlewareFromHandler(h)
	}
	if h, ok := m.(func(Context) error); ok {
		return WrapMiddlewareFromHandler(HandlerFunc(h))
	}
	if h, ok := m.(func(Handler) func(Context) error); ok {
		return MiddlewareFunc(func(next Handler) Handler {
			return HandlerFunc(h(next))
		})
	}
	if h, ok := m.(func(Handler) HandlerFunc); ok {
		return MiddlewareFunc(func(next Handler) Handler {
			return h(next)
		})
	}
	if h, ok := m.(func(HandlerFunc) HandlerFunc); ok {
		return MiddlewareFunc(func(next Handler) Handler {
			return h(next.Handle)
		})
	}
	if h, ok := m.(func(Handler) Handler); ok {
		return MiddlewareFunc(h)
	}
	if h, ok := m.(func(func(Context) error) func(Context) error); ok {
		return MiddlewareFunc(func(next Handler) Handler {
			return HandlerFunc(h(next.Handle))
		})
	}
	if v, ok := m.(http.Handler); ok {
		return WrapMiddlewareFromStdHandler(v)
	}
	if v, ok := m.(func(http.ResponseWriter, *http.Request)); ok {
		return WrapMiddlewareFromStdHandleFunc(v)
	}
	if v, ok := m.(func(http.ResponseWriter, *http.Request) error); ok {
		return WrapMiddlewareFromStdHandleFuncd(v)
	}
	panic(`unknown middleware`)
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
			h.ServeHTTP(c.Response().StdResponseWriter(), c.Request().StdRequest())
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
			h(c.Response().StdResponseWriter(), c.Request().StdRequest())
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
			if err := h(c.Response().StdResponseWriter(), c.Request().StdRequest()); err != nil {
				return err
			}
			if c.Response().Committed() {
				return nil
			}
			return next.Handle(c)
		})
	})
}
