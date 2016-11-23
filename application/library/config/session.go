package config

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
)

var (
	SessionOptions *echo.SessionOptions
	CookieOptions  *cookie.CookieOptions
)

func InitSessionOptions() {
	if DefaultConfig.Cookie.Path == `` {
		DefaultConfig.Cookie.Path = `/`
	}
	SessionOptions = &echo.SessionOptions{
		Name:   "SID",
		Engine: "cookie",
		CookieOptions: &echo.CookieOptions{
			Domain:   DefaultConfig.Cookie.Domain,
			Path:     DefaultConfig.Cookie.Path,
			MaxAge:   DefaultConfig.Cookie.MaxAge,
			HttpOnly: DefaultConfig.Cookie.HttpOnly,
		},
	}
	CookieOptions = &cookie.CookieOptions{
		KeyPairs:       [][]byte{},
		SessionOptions: SessionOptions,
	}
	if len(DefaultConfig.Cookie.HashKey) > 0 {
		CookieOptions.KeyPairs = append(CookieOptions.KeyPairs, []byte(DefaultConfig.Cookie.HashKey))

		if len(DefaultConfig.Cookie.BlockKey) > 0 && DefaultConfig.Cookie.BlockKey != DefaultConfig.Cookie.HashKey {
			CookieOptions.KeyPairs = append(CookieOptions.KeyPairs, []byte(DefaultConfig.Cookie.BlockKey))
		}
	}
	cookie.RegWithOptions(CookieOptions)
}
