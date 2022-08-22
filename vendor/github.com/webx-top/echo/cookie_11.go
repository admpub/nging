//go:build go1.11
// +build go1.11

package echo

import (
	"net/http"
	"strings"
)

// CookieSameSite 设置SameSite
func CookieSameSite(stdCookie *http.Cookie, p string) {
	switch strings.ToLower(p) {
	case `lax`:
		stdCookie.SameSite = http.SameSiteLaxMode
	case `strict`:
		stdCookie.SameSite = http.SameSiteStrictMode
	default:
		stdCookie.SameSite = http.SameSiteDefaultMode
	}
}

// CopyCookieOptions copy cookie options
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
	to.SameSite = from.SameSite
}
