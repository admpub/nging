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

package handler

import (
	"html/template"
	"net/url"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"

	"github.com/nging-plugins/caddymanager/application/library/webdav"
)

func NewFormValues(values url.Values) *FormValues {
	return &FormValues{Values: values}
}

type FormValues struct {
	url.Values
}

func (v FormValues) GetWebdavGlobal() []*webdav.WebdavPerm {
	return webdav.ParseGlobalForm(v.Values)
}

func (v FormValues) GetWebdavUser() []*webdav.WebdavUser {
	return webdav.ParseUserForm(v.Values)
}

func (v FormValues) GetSlice(key string) param.StringSlice {
	values, _ := v.Values[key]
	return param.StringSlice(values)
}

func (v FormValues) AddSlashes(val string) string {
	return com.AddCSlashes(val, '"')
}

func (v FormValues) IteratorKV(addon string, item string, prefix string, withQuotes ...bool) interface{} {
	if len(addon) > 0 && len(item) > 0 {
		addon += `_`
	}
	k := addon + item + `_k`
	keys, _ := v.Values[k]

	k = addon + item + `_v`
	values, _ := v.Values[k]

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
				v = `"` + com.AddCSlashes(v, '"') + `"`
			}
			r += t + prefix + k + `   ` + v
			t = "\n"
		}
	}
	if withQuote {
		return template.HTML(r)
	}
	return r
}

func (v FormValues) GetAttrVal(addon string, item string, defaults ...string) string {
	if len(addon) > 0 {
		addon += `_`
	}
	k := addon + item
	val := v.Values.Get(k)
	if len(val) == 0 {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ``
	}
	return val
}

func (v FormValues) AddonAttr(addon string, item string, defaults ...string) string {
	if len(addon) > 0 {
		addon += `_`
	}
	k := addon + item
	val := v.Values.Get(k)
	if len(val) == 0 {
		if len(defaults) > 0 && len(defaults[0]) > 0 {
			return item + `   ` + defaults[0]
		}
		return ``
	}
	return item + `   ` + val
}

func (v FormValues) Iterator(addon string, item string, prefix string, withQuotes ...bool) interface{} {
	if len(addon) > 0 {
		addon += `_`
	}
	k := addon + item
	values, _ := v.Values[k]
	var r, t string
	var withQuote bool
	if len(withQuotes) > 0 {
		withQuote = withQuotes[0]
	}
	for _, v := range values {
		if withQuote {
			v = `"` + com.AddCSlashes(v, '"') + `"`
		}
		r += t + prefix + v
		t = "\n"
	}
	if withQuote {
		return template.HTML(r)
	}
	return r
}
