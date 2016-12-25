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
package config

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
)

var (
	SessionOptions *echo.SessionOptions
	CookieOptions  *cookie.CookieOptions
)

func InitSessionOptions() {
	if DefaultConfig.Cookie.Path == `` {
		DefaultConfig.Cookie.Path = `/`
	}
	SessionOptions = &echo.SessionOptions{
		Name:   "SID",
		Engine: "cookie",
		CookieOptions: &echo.CookieOptions{
			Domain:   DefaultConfig.Cookie.Domain,
			Path:     DefaultConfig.Cookie.Path,
			MaxAge:   DefaultConfig.Cookie.MaxAge,
			HttpOnly: DefaultConfig.Cookie.HttpOnly,
		},
	}
	CookieOptions = &cookie.CookieOptions{
		KeyPairs:       [][]byte{},
		SessionOptions: SessionOptions,
	}
	if len(DefaultConfig.Cookie.HashKey) > 0 {
		CookieOptions.KeyPairs = append(CookieOptions.KeyPairs, []byte(DefaultConfig.Cookie.HashKey))

		if len(DefaultConfig.Cookie.BlockKey) > 0 && DefaultConfig.Cookie.BlockKey != DefaultConfig.Cookie.HashKey {
			CookieOptions.KeyPairs = append(CookieOptions.KeyPairs, []byte(DefaultConfig.Cookie.BlockKey))
		}
	}
	cookie.RegWithOptions(CookieOptions)
}
