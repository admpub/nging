// +build !go1.11

package echo

import (
	"net/http"
)

//NewCookie create a cookie instance
func NewCookie(name string, value string, opts ...*CookieOptions) *Cookie {
	opt := &CookieOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	cookie := newCookie(name, value, opt)
	return cookie
}

//SameSite 设置SameSite
func (c *Cookie) SameSite(_ string) *Cookie {
	return c
}

func CopyCookieOptions(from *http.Cookie, to *Cookie) {
	to.Expires(from.MaxAge)
	if len(from.Path) == 0 {
		from.Path = `/`
	}
	to.Path(from.Path)
	to.Domain(from.Domain)
	to.Secure(from.Secure)
	to.HttpOnly(from.HttpOnly)
}
