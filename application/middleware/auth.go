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

	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/license"
	"github.com/admpub/nging/application/model"
	"github.com/admpub/nging/application/registry/perm"
)

func AuthCheck(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		//检查是否已安装
		if !config.IsInstalled() {
			c.Data().SetError(c.E(`请先安装`))
			return c.Redirect(handler.URLFor(`/setup`))
		}

		//验证授权文件
		if !license.Ok(c.Host()) {
			c.Data().SetError(c.E(`请先获取本系统授权`))
			return c.Redirect(handler.URLFor(`/license`))
		}

		if user := handler.User(c); user != nil {
			if jump, ok := c.Session().Get(`auth2ndURL`).(string); ok && len(jump) > 0 {
				c.Data().SetError(c.E(`请先进行第二步验证`))
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
			permission, ok := c.Internal().Get(`permission`).(*model.RolePermission)
			if !ok {
				return echo.ErrForbidden
			}
			if checker, ok := perm.SpecialAuths[rpath]; ok {
				var err error
				var ret bool
				err, ppath, ret = checker(h, c, rpath, user, permission)
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
						ppath = `/term/client`
					default:
						if strings.HasPrefix(rpath, `/frp/dashboard/`) {
							ppath = `/frp/dashboard`
						}
					}
				}
			}
			c.SetFunc(`CheckPerm`, func(route string) error {
				if user.Id == 1 {
					return nil
				}
				if !permission.Check(route) {
					return echo.ErrForbidden
				}
				return nil
			})
			return h.Handle(c)
		}
		c.Data().SetError(c.E(`请先登录`))
		return c.Redirect(handler.URLFor(`/login`))
	}
}

// CheckAnyPerm 检查是否匹配任意给定路径权限
func CheckAnyPerm(c echo.Context, ppaths ...string) (err error) {
	check := c.GetFunc(`CheckPerm`).(func(string) error)
	for _, ppath := range ppaths {
		if err = check(ppath); err == nil {
			return nil
		}
	}
	return err
}

// CheckAllPerm 检查是否匹配所有给定路径权限
func CheckAllPerm(c echo.Context, ppaths ...string) (err error) {
	check := c.GetFunc(`CheckPerm`).(func(string) error)
	for _, ppath := range ppaths {
		if err = check(ppath); err != nil {
			return err
		}
	}
	return nil
}

func Auth(c echo.Context, saveSession bool) error {
	user := c.Form(`user`)
	pass := c.Form(`pass`)
	common.DecryptedByRandomSecret(c, `loginPassword`, &pass)
	m := model.NewUser(c)
	exists, err := m.CheckPasswd(user, pass)
	if !exists {
		return c.E(`用户不存在`)
	}
	if err == nil {
		if saveSession {
			m.SetSession()
		}
		if m.NeedCheckU2F(m.NgingUser.Id) {
			c.Session().Set(`auth2ndURL`, handler.URLFor(`/gauth_check`))
		}
		m.NgingUser.LastLogin = uint(time.Now().Unix())
		m.NgingUser.LastIp = c.RealIP()
		m.NgingUser.SetFields(nil, map[string]interface{}{
			`last_login`: m.NgingUser.LastLogin,
			`last_ip`:    m.NgingUser.LastIp,
		}, `id`, m.NgingUser.Id)
		common.DeleteRandomSecret(c, `loginPassword`)
	}
	return err
}
