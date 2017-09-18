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
	"net/http"
	"net/url"
	"time"
)

// CookieOptions cookie options
type CookieOptions struct {
	Prefix string

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge int

	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}

func (c *CookieOptions) Clone() *CookieOptions {
	clone := *c
	return &clone
}

//Cookier interface
type Cookier interface {
	Get(key string) string
	Set(key string, val string, args ...interface{}) Cookier
}

//NewCookier create a cookie instance
func NewCookier(ctx Context) Cookier {
	return &cookie{
		context: ctx,
		cookies: []*Cookie{},
	}
}

//NewCookie create a cookie instance
func NewCookie(name string, value string, opts ...*CookieOptions) *Cookie {
	opt := &CookieOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	if len(opt.Path) == 0 {
		opt.Path = `/`
	}
	cookie := &Cookie{
		cookie: &http.Cookie{
			Name:     opt.Prefix + name,
			Value:    value,
			Path:     opt.Path,
			Domain:   opt.Domain,
			MaxAge:   opt.MaxAge,
			Secure:   opt.Secure,
			HttpOnly: opt.HttpOnly,
		},
	}
	return cookie
}

//Cookie 操作封装
type Cookie struct {
	cookie *http.Cookie
}

//Path 设置路径
func (c *Cookie) Path(p string) *Cookie {
	c.cookie.Path = p
	return c
}

//Domain 设置域名
func (c *Cookie) Domain(p string) *Cookie {
	c.cookie.Domain = p
	return c
}

//MaxAge 设置有效时长（秒）
func (c *Cookie) MaxAge(p int) *Cookie {
	c.cookie.MaxAge = p
	return c
}

//Expires 设置过期时间戳
func (c *Cookie) Expires(p int64) *Cookie {
	if p > 0 {
		c.cookie.Expires = time.Unix(time.Now().Unix()+p, 0)
	} else if p < 0 {
		c.cookie.Expires = time.Unix(1, 0)
	}
	return c
}

//Secure 设置是否启用HTTPS
func (c *Cookie) Secure(p bool) *Cookie {
	c.cookie.Secure = p
	return c
}

//HttpOnly 设置是否启用HttpOnly
func (c *Cookie) HttpOnly(p bool) *Cookie {
	c.cookie.HttpOnly = p
	return c
}

//Send 发送cookie数据到响应头
func (c *Cookie) Send(ctx Context) {
	ctx.Response().SetCookie(c.cookie)
}

type cookie struct {
	context Context
	cookies []*Cookie
}

func (c *cookie) Get(key string) string {
	var val string
	if v := c.context.Request().Cookie(c.context.CookieOptions().Prefix + key); len(v) > 0 {
		val, _ = url.QueryUnescape(v)
	}
	return val
}

func (c *cookie) Set(key string, val string, args ...interface{}) Cookier {
	val = url.QueryEscape(val)
	var cookie *Cookie
	var found bool
	for _, v := range c.cookies {
		if key == v.cookie.Name {
			cookie = v
			found = true
			break
		}
	}
	if cookie == nil {
		cookie = NewCookie(key, val, c.context.CookieOptions())
	}
	switch len(args) {
	case 5:
		httpOnly, _ := args[4].(bool)
		cookie.HttpOnly(httpOnly)
		fallthrough
	case 4:
		secure, _ := args[3].(bool)
		cookie.Secure(secure)
		fallthrough
	case 3:
		domain, _ := args[2].(string)
		cookie.Domain(domain)
		fallthrough
	case 2:
		ppath, _ := args[1].(string)
		cookie.Path(ppath)
		fallthrough
	case 1:
		var liftTime int64
		switch args[0].(type) {
		case int:
			liftTime = int64(args[0].(int))
		case int64:
			liftTime = args[0].(int64)
		case time.Duration:
			liftTime = int64(args[0].(time.Duration).Seconds())
		}
		cookie.Expires(liftTime)
	}
	if !found {
		c.cookies = append(c.cookies, cookie)
		cookie.Send(c.context)
	} else {
		c.context.Response().Header().Del(HeaderSetCookie)
		for _, cookie := range c.cookies {
			cookie.Send(c.context)
		}
	}
	return c
}
