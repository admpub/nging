/*
Nging is a toolbox for webmasters
Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package config

import (
	"reflect"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session/engine"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
	"github.com/webx-top/echo/param"
)

type SessionStoreInit func(c *Config, cookieOptions *cookie.CookieOptions, sessionConfig param.Store) (changed bool, err error)

var (
	SessionOptions *echo.SessionOptions
	CookieOptions  *echo.CookieOptions
	SessionEngine  = `file`
	SessionName    = `SID`
	SessionStores  = echo.NewKVxData[SessionStoreInit, any]()

	sessionStoreCookieOptions *cookie.CookieOptions
)

func RegisterSessionStore(name string, title string, initFn SessionStoreInit) {
	SessionStores.Add(name, title, echo.KVxOptX[SessionStoreInit, any](initFn))
}

func InitSessionOptions(c *Config) {

	//==================================
	// session基础设置
	//==================================

	if len(c.Cookie.Path) == 0 {
		c.Cookie.Path = `/`
	}
	if len(c.Cookie.Prefix) == 0 {
		c.Cookie.Prefix = `Nging`
	}
	sessionName := c.Sys.SessionName
	sessionEngine := c.Sys.SessionEngine
	sessionConfig := c.Sys.SessionConfig
	if len(sessionName) == 0 {
		sessionName = SessionName
	}
	if len(sessionEngine) == 0 {
		sessionEngine = SessionEngine
	}
	if sessionConfig == nil {
		sessionConfig = echo.H{}
	}
	_cookieOptions := &echo.CookieOptions{
		Prefix:   c.Cookie.Prefix,
		Domain:   c.Cookie.Domain,
		Path:     c.Cookie.Path,
		MaxAge:   c.Cookie.MaxAge,
		HttpOnly: c.Cookie.HttpOnly,
		SameSite: c.Cookie.SameSite,
	}
	if CookieOptions == nil || SessionOptions == nil ||
		!reflect.DeepEqual(_cookieOptions, CookieOptions) ||
		(SessionOptions.Engine != sessionEngine || SessionOptions.Name != sessionName) {
		if SessionOptions != nil {
			*SessionOptions = *echo.NewSessionOptions(sessionEngine, sessionName, _cookieOptions)
		} else {
			SessionOptions = echo.NewSessionOptions(sessionEngine, sessionName, _cookieOptions)
		}
		if CookieOptions != nil {
			*CookieOptions = *_cookieOptions
		} else {
			CookieOptions = _cookieOptions
		}
	}

	//==================================
	// 注册session存储引擎
	//==================================

	//1. 注册默认引擎：cookie
	_sessionStoreCookieOptions := cookie.NewCookieOptions(c.Cookie.HashKey, c.Cookie.BlockKey)
	if sessionStoreCookieOptions == nil || !reflect.DeepEqual(_sessionStoreCookieOptions, sessionStoreCookieOptions) {
		cookie.RegWithOptions(_sessionStoreCookieOptions)
		sessionStoreCookieOptions = _sessionStoreCookieOptions
	}

	ss := SessionStores.GetItem(sessionEngine)
	if ss == nil {
		log.Errorf(`unsupported session engine: %s`, sessionEngine)
		return
	}
	changed, err := ss.X(c, sessionStoreCookieOptions, sessionConfig)
	if err != nil {
		log.Error(err)
		return
	}
	if changed {
		for _, v := range SessionStores.Slice() {
			if v.K != sessionEngine {
				engine.Del(v.K)
			}
		}
	}
}

func AutoSecure(ctx echo.Context, ses *echo.SessionOptions) {
	if !ses.Secure && ctx.IsSecure() {
		ses.Secure = true
	}
}
