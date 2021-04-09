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

package caddy

import (
	"net/url"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func addonAttr(ctx echo.Context, v url.Values) {
	ctx.SetFunc(`Get`, func(addon string, item string, defaults ...string) string {
		if len(addon) > 0 {
			addon += `_`
		}
		k := addon + item
		v := v.Get(k)
		if len(v) == 0 {
			if len(defaults) > 0 {
				return defaults[0]
			}
			return ``
		}
		return v
	})
	ctx.SetFunc(`AddonAttr`, func(addon string, item string, defaults ...string) string {
		if len(addon) > 0 {
			addon += `_`
		}
		k := addon + item
		v := v.Get(k)
		if len(v) == 0 {
			if len(defaults) > 0 && len(defaults[0]) > 0 {
				return item + `   ` + defaults[0]
			}
			return ``
		}
		return item + `   ` + v
	})
	ctx.SetFunc(`Iterator`, func(addon string, item string, prefix string, withQuotes ...bool) string {
		if len(addon) > 0 {
			addon += `_`
		}
		k := addon + item
		values, _ := v[k]
		var r, t string
		var withQuote bool
		if len(withQuotes) > 0 {
			withQuote = withQuotes[0]
		}
		for _, v := range values {
			if withQuote {
				v = `"` + com.AddSlashes(v, '"') + `"`
			}
			r += t + prefix + v
			t = "\n"
		}
		return r
	})
}

func iteratorKV(ctx echo.Context, v url.Values) {
	ctx.SetFunc(`IteratorKV`, func(addon string, item string, prefix string, withQuotes ...bool) string {
		if len(addon) > 0 {
			addon += `_`
		}
		k := addon + item + `_k`
		keys, _ := v[k]

		k = addon + item + `_v`
		values, _ := v[k]

		var r, t string
		var withQuote bool
		if len(withQuotes) > 0 {
			withQuote = withQuotes[0]
		}
		l := len(values)
		for i, k := range keys {
			if i < l {
				v := values[i]
				if withQuote {
					v = `"` + com.AddSlashes(v, '"') + `"`
				}
				r += t + prefix + k + `   ` + v
				t = "\n"
			}
		}
		return r
	})
}

func SetCaddyfileFunc(ctx echo.Context, v url.Values) {
	addonAttr(ctx, v)
	iteratorKV(ctx, v)
	ctx.SetFunc(`AddSlashes`, func(v string) string {
		return com.AddSlashes(v, '"')
	})
	ctx.SetFunc(`GetSlice`, func(key string) param.StringSlice {
		values, _ := v[key]
		return param.StringSlice(values)
	})
	ctx.SetFunc(`GetWebdavUser`, func() []*WebdavUser {
		return parseWebdavUserForm(v)
	})
	ctx.SetFunc(`GetWebdavGlobal`, func() []*WebdavPerm {
		return parseWebdavGlobalForm(v)
	})
}

type WebdavPerm struct {
	Perm     string
	Resource string
}

func (w *WebdavPerm) String() string {
	if strings.Contains(w.Resource, `*`) {
		return w.Perm + `_r    ` + strings.Replace(w.Resource, `*`, `(.*)`, -1)
	}
	return w.Perm + `      ` + w.Resource
}

type WebdavUser struct {
	User     string
	Password string
	Perms    []*WebdavPerm
}

func parseWebdavUserForm(v url.Values) []*WebdavUser {
	indexes, _ := v[`webdav_user_index`]
	users, _ := v[`webdav_user`]
	passwords, _ := v[`webdav_pass`]
	var list []*WebdavUser
	for key, index := range indexes {
		if key >= len(users) || key >= len(passwords) {
			continue
		}
		u := &WebdavUser{
			User:     users[key],
			Password: passwords[key],
		}
		if len(u.User) == 0 {
			continue
		}
		perms, _ := v[`webdav_perms[user][`+index+`]`]
		resources, _ := v[`webdav_resources[user][`+index+`]`]
		for pkey, resource := range resources {
			if pkey >= len(perms) || len(resource) == 0 {
				continue
			}
			p := &WebdavPerm{Perm: perms[pkey], Resource: resource}
			if p.Perm == `1` {
				p.Perm = `allow`
			} else {
				p.Perm = `block`
			}
			u.Perms = append(u.Perms, p)
		}
		list = append(list, u)
	}
	return list
}

func parseWebdavGlobalForm(v url.Values) []*WebdavPerm {
	var list []*WebdavPerm
	perms, _ := v[`webdav_perms[global]`]
	resources, _ := v[`webdav_resources[global]`]
	for pkey, resource := range resources {
		if pkey >= len(perms) || len(resource) == 0 {
			continue
		}
		p := &WebdavPerm{Perm: perms[pkey], Resource: resource}
		if p.Perm == `1` {
			p.Perm = `allow`
		} else {
			p.Perm = `block`
		}
		list = append(list, p)
	}
	return list
}
