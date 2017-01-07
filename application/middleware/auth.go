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
package middleware

import (
	"errors"

	"github.com/admpub/nging/application/library/config"
	"github.com/webx-top/echo"
)

func AuthCheck(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		if user, _ := c.Session().Get(`user`).(string); len(user) > 0 {
			c.Set(`user`, user)
			c.SetFunc(`Username`, func() string { return user })
			return h.Handle(c)
		}

		//临时认证
		if err := Auth(c, false); err != nil {

			//检查是否已安装
			if !config.IsInstalled() {
				return c.Redirect(`/setup`)
			}

			return c.Redirect(`/login`)
		}
		return h.Handle(c)
	}
}

func Auth(c echo.Context, saveSession bool) error {
	user := c.Form(`user`)
	pass := c.Form(`pass`)
	if pwd, ok := config.DefaultConfig.Sys.Accounts[user]; ok && pwd == pass {
		if saveSession {
			c.Session().Set(`user`, user)
		}
		return nil
	}
	return errors.New(c.T(`登录失败，用户名或密码不正确`))
}
