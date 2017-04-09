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
	"time"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/echo"
)

func AuthCheck(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		if user, _ := c.Session().Get(`user`).(*dbschema.User); user != nil {
			if jump, ok := c.Session().Get(`auth2ndURL`).(string); ok && len(jump) > 0 {
				return c.Redirect(jump)
			}
			c.Set(`user`, user)
			c.SetFunc(`Username`, func() string { return user.Username })
			return h.Handle(c)
		}

		//检查是否已安装
		if !config.IsInstalled() {
			return c.Redirect(`/setup`)
		}
		return c.Redirect(`/login`)
	}
}

func Auth(c echo.Context, saveSession bool) error {
	user := c.Form(`user`)
	pass := c.Form(`pass`)

	m := model.NewUser(c)
	exists, err := m.CheckPasswd(user, pass)
	if !exists {
		return errors.New(c.T(`用户不存在`))
	}
	if err == nil {
		if saveSession {
			m.SetSession()
		}
		if m.NeedCheckU2F(m.User.Id) {
			c.Session().Set(`auth2ndURL`, `/gauth_check`)
		}
		m.User.LastLogin = uint(time.Now().Unix())
		m.User.LastIp = c.RealIP()
		m.User.Param().SetSend(map[string]interface{}{
			`last_login`: m.User.LastLogin,
			`last_ip`:    m.User.LastIp,
		}).SetArgs(`id`, m.User.Id).Update()
	}
	return err
}
