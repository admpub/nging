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

	"github.com/webx-top/echo"
)

func addonAttr(ctx echo.Context, v url.Values) {
	ctx.SetFunc(`AddonAttr`, func(addon string, item string) string {
		if len(addon) > 0 {
			addon += `_`
		}
		k := addon + item
		v := v.Get(k)
		if len(v) == 0 {
			return ``
		}
		return item + `   ` + v
	})
	ctx.SetFunc(`Iterator`, func(addon string, item string, prefix string) string {
		if len(addon) > 0 {
			addon += `_`
		}
		k := addon + item
		values, _ := v[k]
		r := ``
		t := ``
		for i, v := range values {
			r += t + prefix + v + `   ` + values[i]
			t = "\n"
		}
		return r
	})
}

func iteratorKV(ctx echo.Context, v url.Values) {
	ctx.SetFunc(`IteratorKV`, func(addon string, item string, prefix string) string {
		if len(addon) > 0 {
			addon += `_`
		}
		k := addon + item + `_k`
		keys, _ := v[k]

		k = addon + item + `_v`
		values, _ := v[k]

		r := ``
		l := len(values)
		t := ``
		for i, v := range keys {
			if i < l {
				r += t + prefix + v + `   ` + values[i]
				t = "\n"
			}
		}
		return r
	})
}

func SetCaddyfileFunc(ctx echo.Context, v url.Values) {
	addonAttr(ctx, v)
	iteratorKV(ctx, v)
}
