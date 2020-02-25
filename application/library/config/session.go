/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

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
	"path/filepath"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
	"github.com/webx-top/echo/middleware/session/engine/file"
)

var (
	SessionOptions *echo.SessionOptions
	CookieOptions  *cookie.CookieOptions
	SessionEngine  = `file`
	SessionName    = `SID`
)

func InitSessionOptions(c *Config) {

	//==================================
	// session基础设置
	//==================================

	if len(c.Cookie.Path) == 0 {
		c.Cookie.Path = `/`
	}
	sessionName := c.Sys.SessionName
	if len(sessionName) == 0 {
		sessionName = SessionName
	}
	SessionOptions = echo.NewSessionOptions(SessionEngine, sessionName, &echo.CookieOptions{
		Domain:   c.Cookie.Domain,
		Path:     c.Cookie.Path,
		MaxAge:   c.Cookie.MaxAge,
		HttpOnly: c.Cookie.HttpOnly,
	})

	//==================================
	// 注册session存储引擎
	//==================================

	//1. 注册默认引擎：cookie
	CookieOptions = cookie.NewCookieOptions(c.Cookie.HashKey, c.Cookie.BlockKey)
	cookie.RegWithOptions(CookieOptions)

	//2. 注册文件引擎：file
	file.RegWithOptions(&file.FileOptions{
		SavePath: filepath.Join(echo.Wd(), `data`, `cache`, `sessions`),
		KeyPairs: CookieOptions.KeyPairs,
	})
}
