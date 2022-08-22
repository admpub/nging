//go:build !go1.11
// +build !go1.11

package echo

import (
	"net/http"
)

// CookieSameSite 设置SameSite
func CookieSameSite(_ *http.Cookie, _ string) {
}

func CopyCookieOptions(from *http.Cookie, to *http.Cookie) {
	to.MaxAge = from.MaxAge
	to.Expires = from.Expires
	if len(from.Path) == 0 {
		from.Path = `/`
	}
	to.Path = from.Path
	to.Domain = from.Domain
	to.Secure = from.Secure
	to.HttpOnly = from.HttpOnly
}
