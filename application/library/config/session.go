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
	"path/filepath"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session/engine"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
	"github.com/webx-top/echo/middleware/session/engine/file"
	"github.com/webx-top/echo/middleware/session/engine/redis"
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
	SessionOptions = echo.NewSessionOptions(sessionEngine, sessionName, &echo.CookieOptions{
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

	switch sessionEngine {
	case `file`: //2. 注册文件引擎：file
		fileOptions := &file.FileOptions{
			SavePath: sessionConfig.String(`savePath`),
			KeyPairs: CookieOptions.KeyPairs,
			MaxAge:   sessionConfig.Int(`maxAge`),
		}
		if len(fileOptions.SavePath) == 0 {
			fileOptions.SavePath = filepath.Join(echo.Wd(), `data`, `cache`, `sessions`)
		}
		file.RegWithOptions(fileOptions)
		engine.Del(`redis`)
	case `redis`: //3. 注册redis引擎：redis
		redisOptions := &redis.RedisOptions{
			Size:         sessionConfig.Int(`maxIdle`),
			Network:      sessionConfig.String(`network`),
			Address:      sessionConfig.String(`address`),
			Password:     sessionConfig.String(`password`),
			DB:           sessionConfig.Uint(`db`),
			KeyPairs:     CookieOptions.KeyPairs,
			MaxAge:       sessionConfig.Int(`maxAge`),
			MaxReconnect: sessionConfig.Int(`maxReconnect`),
		}
		if redisOptions.Size <= 0 {
			redisOptions.Size = 10
		}
		if len(redisOptions.Network) == 0 {
			redisOptions.Network = `tcp`
		}
		if len(redisOptions.Address) == 0 {
			redisOptions.Address = `127.0.0.1:6379`
		}
		redis.RegWithOptions(redisOptions)
		engine.Del(`file`)
	}
}
