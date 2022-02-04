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

package middleware

import (
	"strings"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/codec"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/license"
	"github.com/admpub/nging/v4/application/library/role"
	"github.com/admpub/nging/v4/application/model"
)

func AuthCheck(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		//检查是否已安装
		if !config.IsInstalled() {
			c.Data().SetError(c.E(`请先安装`))
			return c.Redirect(handler.URLFor(`/setup`))
		}

		//验证授权文件
		if !license.Ok(c) {
			c.Data().SetError(c.E(`请先获取本系统授权`))
			return c.Redirect(handler.URLFor(`/license`))
		}
		handlerPermission := c.Route().String(`permission`)
		if handlerPermission == `guest` {
			return h.Handle(c)
		}
		user := handler.User(c)
		if user == nil {
			c.Data().SetError(c.E(`请先登录`))
			return c.Redirect(handler.URLFor(`/login`))
		}
		if jump, ok := c.Session().Get(`auth2ndURL`).(string); ok && len(jump) > 0 {
			c.Data().SetError(c.E(`请先进行第二步验证`))
			return c.Redirect(jump)
		}
		var (
			rpath = c.Path()
			upath = c.Request().URL().Path()
			ppath string
		)
		//println(`--------------------->>>`, rpath)
		if len(handler.BackendPrefix) > 0 {
			rpath = strings.TrimPrefix(rpath, handler.BackendPrefix)
		}
		//echo.Dump(c.Route().Meta)
		if user.Id == 1 || strings.HasPrefix(rpath, `/user/`) {
			c.SetFunc(`CheckPerm`, func(route string) error {
				return nil
			})
			return h.Handle(c)
		}
		permission := UserPermission(c)
		c.SetFunc(`CheckPerm`, func(route string) error {
			if user.Id == 1 {
				return nil
			}
			if !permission.Check(c, route) {
				return common.ErrUserNoPerm
			}
			return nil
		})
		if handlerPermission == `public` {
			return h.Handle(c)
		}
		checker, ok := role.SpecialAuths[rpath]
		if !ok {
			checker, ok = role.SpecialAuths[upath]
		}
		if ok {
			var (
				ret bool
				err error
			)
			if ppath, ret, err = checker(h, c, user, permission); ret {
				return err
			} else if err != nil {
				return err
			}
		} else {
			ppath = rpath
		}
		if !permission.Check(c, ppath) {
			return common.ErrUserNoPerm
		}
		return h.Handle(c)
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
	return common.ErrUserNoPerm
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
	var err error
	pass, err = codec.DefaultSM2DecryptHex(pass)
	if err != nil {
		return c.NewError(code.InvalidParameter, c.T(`密码解密失败: %v`, err))
	}
	loginLogM := model.NewLoginLog(c)
	loginLogM.OwnerType = `user`
	loginLogM.Username = user
	m := model.NewUser(c)
	exists, err := m.CheckPasswd(user, pass)
	if !exists {
		loginLogM.Errpwd = pass
		loginLogM.Failmsg = c.T(`用户不存在`)
		loginLogM.Success = `N`
		loginLogM.Add()
		return c.NewError(code.UserNotFound, c.T(`用户不存在`))
	}
	if err == nil {
		if saveSession {
			m.SetSession()
		}
		if m.NeedCheckU2F(m.NgingUser.Id) {
			c.Session().Set(`auth2ndURL`, handler.URLFor(`/gauth_check`))
		}
		m.NgingUser.LastLogin = uint(time.Now().Unix())
		set := echo.H{
			`last_login`: m.NgingUser.LastLogin,
		}
		if !echo.Bool(`backend.Anonymous`) {
			m.NgingUser.LastIp = c.RealIP()
			set[`last_ip`] = m.NgingUser.LastIp
		}
		if len(m.NgingUser.SessionId) > 0 {
			if m.NgingUser.SessionId != c.Session().MustID() {
				c.Session().RemoveID(m.NgingUser.SessionId)
				set.Set(`session_id`, c.Session().MustID())
			}
		} else {
			set.Set(`session_id`, c.Session().MustID())
		}
		m.NgingUser.UpdateFields(nil, set, `id`, m.NgingUser.Id)

		loginLogM.OwnerId = uint64(m.Id)
		loginLogM.Success = `Y`
	} else {
		loginLogM.Errpwd = pass
		loginLogM.Failmsg = err.Error()
		loginLogM.Success = `N`
	}
	loginLogM.Add()
	return err
}
