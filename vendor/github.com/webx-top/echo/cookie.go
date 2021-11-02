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
	"sync"
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
	Add(cookies ...*http.Cookie) Cookier
	Set(key string, val string, args ...interface{}) Cookier
	Send()
}

//NewCookier create a cookie instance
func NewCookier(ctx Context) Cookier {
	return &cookie{
		context: ctx,
		cookies: []*http.Cookie{},
		indexes: map[string]int{},
	}
}

type cookie struct {
	context Context
	cookies []*http.Cookie
	indexes map[string]int
	lock    sync.RWMutex
}

func (c *cookie) Send() {
	c.lock.RLock()
	for _, cookie := range c.cookies {
		c.context.Response().SetCookie(cookie)
	}
	c.lock.RUnlock()
}

func (c *cookie) record(stdCookie *http.Cookie) {
	c.lock.Lock()
	if idx, ok := c.indexes[stdCookie.Name]; ok {
		c.cookies[idx] = stdCookie
		c.lock.Unlock()
		return
	}
	c.indexes[stdCookie.Name] = len(c.cookies)
	c.cookies = append(c.cookies, stdCookie)
	c.lock.Unlock()
}

func (c *cookie) Get(key string) string {
	var val string
	if v := c.context.Request().Cookie(c.context.CookieOptions().Prefix + key); len(v) > 0 {
		val, _ = url.QueryUnescape(v)
	}
	return val
}

func (c *cookie) Add(cookies ...*http.Cookie) Cookier {
	c.lock.Lock()
	for _, cookie := range c.cookies {
		if idx, ok := c.indexes[cookie.Name]; ok {
			c.cookies[idx] = cookie
			continue
		}
		c.indexes[cookie.Name] = len(c.cookies)
		c.cookies = append(c.cookies, cookie)
	}
	c.lock.Unlock()
	return c
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
	cookie := NewCookie(key, val, c.context.CookieOptions())
	switch len(args) {
	case 6:
		sameSite, _ := args[5].(string)
		CookieSameSite(cookie, sameSite)
		fallthrough
	case 5:
		httpOnly, _ := args[4].(bool)
		cookie.HttpOnly = httpOnly
		fallthrough
	case 4:
		secure, _ := args[3].(bool)
		cookie.Secure = secure
		fallthrough
	case 3:
		domain, _ := args[2].(string)
		cookie.Domain = domain
		fallthrough
	case 2:
		ppath, _ := args[1].(string)
		if len(ppath) == 0 {
			ppath = `/`
		}
		cookie.Path = ppath
		fallthrough
	case 1:
		switch v := args[0].(type) {
		case *http.Cookie:
			CopyCookieOptions(v, cookie)
		case *CookieOptions:
			cookie.MaxAge = v.MaxAge
			cookie.Expires = v.Expires
			if len(v.Path) == 0 {
				v.Path = `/`
			}
			cookie.Path = v.Path
			cookie.Domain = v.Domain
			cookie.Secure = v.Secure
			cookie.HttpOnly = v.HttpOnly
			CookieSameSite(cookie, v.SameSite)
		case int:
			CookieMaxAge(cookie, v)
		case int64:
			CookieMaxAge(cookie, int(v))
		case time.Duration:
			CookieMaxAge(cookie, int(v.Seconds()))
		case time.Time:
			CookieExpires(cookie, v)
		}
	}
	c.record(cookie)
	return c
}

// CookieMaxAge 设置有效时长（秒）
// IE6/7/8不支持
// 如果同时设置了MaxAge和Expires，则优先使用MaxAge
// 设置MaxAge则代表每次保存Cookie都会续期，因为MaxAge是基于保存时间来设置的
func CookieMaxAge(stdCookie *http.Cookie, p int) {
	stdCookie.MaxAge = p
	if p > 0 {
		stdCookie.Expires = time.Unix(time.Now().Unix()+int64(p), 0)
	} else if p < 0 {
		stdCookie.Expires = time.Unix(1, 0)
	} else {
		stdCookie.Expires = param.EmptyTime
	}
}

// CookieExpires 设置过期时间
// 所有浏览器都支持
// 如果仅仅设置Expires，因为过期时间是固定的，所以不会导致保存Cookie时被续期
func CookieExpires(stdCookie *http.Cookie, expires time.Time) {
	if expires.IsZero() {
		return
	}
	stdCookie.MaxAge = 0
	stdCookie.Expires = expires
}

// NewCookie 新建cookie对象
func NewCookie(key, value string, opt *CookieOptions) *http.Cookie {
	c := &http.Cookie{
		Name:     opt.Prefix + key,
		Value:    value,
		Path:     `/`,
		Domain:   opt.Domain,
		MaxAge:   opt.MaxAge,
		Expires:  opt.Expires,
		Secure:   opt.Secure,
		HttpOnly: opt.HttpOnly,
	}
	if len(opt.Path) > 0 {
		c.Path = opt.Path
	}
	if len(opt.SameSite) > 0 {
		CookieSameSite(c, opt.SameSite)
	}
	return c
}
