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

package dashboard

import (
	"github.com/webx-top/echo"
)

func NewGlobalFooter(tmpl string, handle ...func(echo.Context) error) *GlobalFooter {
	var content func(echo.Context) error
	if len(handle) > 0 {
		content = handle[0]
	}
	return &GlobalFooter{Tmpl: tmpl, content: content}
}

type GlobalFooter = Tmplx
type GlobalFooters = Tmplxs

var globalFooters = GlobalFooters{}

func GlobalFooterRegister(footer ...*GlobalFooter) {
	globalFooters.Add(-1, footer...)
}

func GlobalFooterAdd(index int, footer ...*GlobalFooter) {
	globalFooters.Add(index, footer...)
}

//GlobalFooterRemove 删除元素
func GlobalFooterRemove(index int) {
	globalFooters.Remove(index)
}

//GlobalFooterSet 设置元素
func GlobalFooterSet(index int, list ...*GlobalFooter) {
	globalFooters.Set(index, list...)
}

func GlobalFooterAll(_ echo.Context) GlobalFooters {
	return globalFooters
}
