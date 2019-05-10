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
package middleware

import (
	"strings"
	"time"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/license"
	"github.com/admpub/nging/application/model"
	"github.com/admpub/nging/application/registry/perm"
	"github.com/webx-top/echo"
)

func AuthCheck(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		//检查是否已安装
		if !config.IsInstalled() {
			return c.Redirect(handler.URLFor(`/setup`))
		}

		//验证授权文件
		if !license.Ok(c.Host()) {
			return c.Redirect(handler.URLFor(`/license`))
		}

		if user := handler.User(c); user != nil {
			if jump, ok := c.Session().Get(`auth2ndURL`).(string); ok && len(jump) > 0 {
				return c.Redirect(jump)
			}
			var (
				rpath = c.Path()
				ppath string
			)
			//println(`--------------------->>>`, rpath)
			if len(handler.BackendPrefix) > 0 {
				rpath = strings.TrimPrefix(rpath, handler.BackendPrefix)
			}
			if user.Id == 1 || strings.HasPrefix(rpath, `/user/`) {
				c.SetFunc(`CheckPerm`, func(route string) error {
					return nil
				})
				return h.Handle(c)
			}
			roleList := handler.RoleList(c)
			roleM := model.NewUserRole(c)
			if checker, ok := perm.SpecialAuths[rpath]; ok {
				var err error
				var ret bool
				err, ppath, ret = checker(h, c, rpath, user, roleM, roleList)
				if ret {
					return err
				}
				if err != nil {
					return err
				}
			} else {
				ppath = rpath
				if len(ppath) >= 13 {
					switch ppath[0:13] {
					case `/term/client/`:
						ppath = `term/client`
					default:
						if strings.HasPrefix(rpath, `/frp/dashboard/`) {
							ppath = `/frp/dashboard`
						} else {
							ppath = strings.TrimPrefix(rpath, `/`)
						}
					}
				} else {
					ppath = strings.TrimPrefix(rpath, `/`)
				}
			}
			if !roleM.CheckPerm2(roleList, ppath) {
				return echo.ErrForbidden
			}
			c.SetFunc(`CheckPerm`, func(route string) error {
				if user.Id == 1 {
					return nil
				}
				if !roleM.CheckPerm2(roleList, route) {
					return echo.ErrForbidden
				}
				return nil
			})
			return h.Handle(c)
		}
		return c.Redirect(handler.URLFor(`/login`))
	}
}

func Auth(c echo.Context, saveSession bool) error {
	user := c.Form(`user`)
	pass := c.Form(`pass`)

	m := model.NewUser(c)
	exists, err := m.CheckPasswd(user, pass)
	if !exists {
		return c.E(`用户不存在`)
	}
	if err == nil {
		if saveSession {
			m.SetSession()
		}
		if m.NeedCheckU2F(m.User.Id) {
			c.Session().Set(`auth2ndURL`, handler.URLFor(`/gauth_check`))
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
