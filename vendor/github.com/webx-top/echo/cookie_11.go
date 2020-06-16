// +build go1.11

package echo

import (
	"net/http"
	"strings"
)

//NewCookie create a cookie instance
func NewCookie(name string, value string, opts ...*CookieOptions) *Cookie {
	opt := &CookieOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	cookie := newCookie(name, value, opt)
	if len(opt.SameSite) > 0 {
		cookie.SameSite(opt.SameSite)
	}
	return cookie
}

//SameSite 设置SameSite
func (c *Cookie) SameSite(p string) *Cookie {
	switch strings.ToLower(p) {
	case `lax`:
		c.cookie.SameSite = http.SameSiteLaxMode
	case `strict`:
		c.cookie.SameSite = http.SameSiteStrictMode
	default:
		c.cookie.SameSite = http.SameSiteDefaultMode
	}
	return c
}

//CopyCookieOptions copy cookie options
func CopyCookieOptions(from *http.Cookie, to *Cookie) {
	to.MaxAge(from.MaxAge)
	to.Expires(from.Expires)
	if len(from.Path) == 0 {
		from.Path = `/`
	}
	to.Path(from.Path)
	to.Domain(from.Domain)
	to.Secure(from.Secure)
	to.HttpOnly(from.HttpOnly)
	to.cookie.SameSite = from.SameSite
}
