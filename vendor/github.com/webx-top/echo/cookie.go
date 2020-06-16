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

	"github.com/webx-top/echo/param"
)

var (
	DefaultCookieOptions = &CookieOptions{
		Path: `/`,
	}
)

// CookieOptions cookie options
type CookieOptions struct {
	Prefix string

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge int

	// Expires
	Expires time.Time

	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	SameSite string // strict / lax
}

func (c *CookieOptions) Clone() *CookieOptions {
	clone := *c
	return &clone
}

func (c *CookieOptions) SetMaxAge(maxAge int) *CookieOptions {
	c.MaxAge = maxAge
	c.Expires = param.EmptyTime
	return c
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
func newCookie(name string, value string, opt *CookieOptions) *Cookie {
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
			Expires:  opt.Expires,
			Secure:   opt.Secure,
			HttpOnly: opt.HttpOnly,
		},
	}
	return cookie
}

// Cookie 操作封装
type Cookie struct {
	cookie *http.Cookie
}

// Path 设置路径
func (c *Cookie) Path(p string) *Cookie {
	c.cookie.Path = p
	return c
}

// Domain 设置域名
func (c *Cookie) Domain(p string) *Cookie {
	c.cookie.Domain = p
	return c
}

// MaxAge 设置有效时长（秒）
// IE6/7/8不支持
// 如果同时设置了MaxAge和Expires，则优先使用MaxAge
// 设置MaxAge则代表每次保存Cookie都会续期，因为MaxAge是基于保存时间来设置的
func (c *Cookie) MaxAge(p int) *Cookie {
	c.cookie.MaxAge = p
	if p > 0 {
		c.cookie.Expires = time.Unix(time.Now().Unix()+int64(p), 0)
	} else if p < 0 {
		c.cookie.Expires = time.Unix(1, 0)
	} else {
		c.cookie.Expires = param.EmptyTime
	}
	return c
}

// Expires 设置过期时间
// 所有浏览器都支持
// 如果仅仅设置Expires，因为过期时间是固定的，所以不会导致保存Cookie时被续期
func (c *Cookie) Expires(expires time.Time) *Cookie {
	if expires.IsZero() {
		return c
	}
	c.cookie.MaxAge = 0
	c.cookie.Expires = expires
	return c
}

// Secure 设置是否启用HTTPS
func (c *Cookie) Secure(p bool) *Cookie {
	c.cookie.Secure = p
	return c
}

// HttpOnly 设置是否启用HttpOnly
func (c *Cookie) HttpOnly(p bool) *Cookie {
	c.cookie.HttpOnly = p
	return c
}

// Send 发送cookie数据到响应头
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

// Set Set cookie value
// @param string key
// @param string value
// @param int|int64|time.Duration args[0]:maxAge (seconds)
// @param string args[1]:path (/)
// @param string args[2]:domain
// @param bool args[3]:secure
// @param bool args[4]:httpOnly
// @param string args[5]:sameSite (lax/strict/default)
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
	case 6:
		sameSite, _ := args[5].(string)
		cookie.SameSite(sameSite)
		fallthrough
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
		switch v := args[0].(type) {
		case *http.Cookie:
			CopyCookieOptions(v, cookie)
		case *CookieOptions:
			cookie.MaxAge(v.MaxAge)
			cookie.Expires(v.Expires)
			if len(v.Path) == 0 {
				v.Path = `/`
			}
			cookie.Path(v.Path)
			cookie.Domain(v.Domain)
			cookie.Secure(v.Secure)
			cookie.HttpOnly(v.HttpOnly)
			cookie.SameSite(v.SameSite)
		case int:
			cookie.MaxAge(v)
		case int64:
			cookie.MaxAge(int(v))
		case time.Duration:
			cookie.MaxAge(int(v.Seconds()))
		case time.Time:
			cookie.Expires(v)
		}
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
