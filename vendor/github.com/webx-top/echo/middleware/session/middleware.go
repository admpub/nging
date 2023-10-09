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
package session

import (
	"time"

	"github.com/admpub/sessions"
	"github.com/webx-top/echo"
)

const (
	CookieMaxAgeKey = `CookieMaxAge`
)

func ForgotMaxAge(c echo.Context) {
	if c.Session().Get(CookieMaxAgeKey) != nil {
		c.Session().Delete(CookieMaxAgeKey)
	}
}

func RememberMaxAge(c echo.Context, maxAge int) {
	if maxAge > 0 {
		c.CookieOptions().SetMaxAge(maxAge)
		c.Session().Set(CookieMaxAgeKey, maxAge)
	} else {
		ForgotMaxAge(c)
	}
}

func RememberExpires(c echo.Context, expires time.Time) {
	if !expires.IsZero() {
		c.CookieOptions().Expires = expires
		c.Session().Set(CookieMaxAgeKey, expires)
	} else {
		ForgotMaxAge(c)
	}
}

func saveSession(c echo.Context) {
	if err := c.Session().Save(); err != nil {
		c.Logger().Error(err)
	}
}

func Sessions(options *echo.SessionOptions, store sessions.Store) echo.MiddlewareFuncd {
	var newSession func(ctx echo.Context) echo.Sessioner
	if options == nil {
		newSession = func(ctx echo.Context) echo.Sessioner {
			return NewMySession(store, ctx.SessionOptions().Name, ctx)
		}
	} else {
		newSession = func(ctx echo.Context) echo.Sessioner {
			sessionOptions := options.Clone()
			ctx.SetSessionOptions(sessionOptions)
			return NewMySession(store, sessionOptions.Name, ctx)
		}
	}
	return func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			s := newSession(c)
			c.SetSessioner(s)
			s.SetPreSaveHook(func(c echo.Context) error {
				switch v := s.Get(CookieMaxAgeKey).(type) {
				case int:
					c.CookieOptions().SetMaxAge(v)
				case time.Time:
					c.CookieOptions().Expires = v
				}
				return nil
			})
			c.AddPreResponseHook(func() error {
				saveSession(c)
				return nil
			})
			err := h.Handle(c)
			saveSession(c)
			return err
		}
	}
}

func Middleware(options *echo.SessionOptions) echo.MiddlewareFuncd {
	if options == nil {
		options = DefaultSessionOptions()
	}
	store := StoreEngine(options)
	return Sessions(options, store)
}
