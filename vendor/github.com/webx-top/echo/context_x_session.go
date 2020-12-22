package echo

import "net/http"

func (c *xContext) Session() Sessioner {
	return c.sessioner
}

func (c *xContext) Flash(names ...string) (r interface{}) {
	if v := c.sessioner.Flashes(names...); len(v) > 0 {
		r = v[len(v)-1]
	}
	return r
}

func (c *xContext) SetCookieOptions(opts *CookieOptions) {
	c.SessionOptions().CookieOptions = opts
}

func (c *xContext) CookieOptions() *CookieOptions {
	return c.SessionOptions().CookieOptions
}

func (c *xContext) SetSessionOptions(opts *SessionOptions) {
	c.sessionOptions = opts
}

func (c *xContext) SessionOptions() *SessionOptions {
	if c.sessionOptions == nil {
		c.sessionOptions = DefaultSessionOptions
	}
	return c.sessionOptions
}

func (c *xContext) NewCookie(key string, value string) *http.Cookie {
	return NewCookie(key, value, c.CookieOptions())
}

func (c *xContext) Cookie() Cookier {
	return c.cookier
}

func (c *xContext) GetCookie(key string) string {
	return c.cookier.Get(key)
}

func (c *xContext) SetCookie(key string, val string, args ...interface{}) {
	c.cookier.Set(key, val, args...)
}
