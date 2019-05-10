// +build !go1.11

package echo

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
