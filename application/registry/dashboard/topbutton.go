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

func NewButton(tmpl string, handle ...func(echo.Context) error) *Button {
	var content func(echo.Context) error
	if len(handle) > 0 {
		content = handle[0]
	}
	return &Button{Tmpl: tmpl, content: content}
}

type Button = Tmplx
type Buttons = Tmplxs

var topButtons = Buttons{
	{
		Tmpl: `manager/topbutton/donation`,
	},
	{
		Tmpl: `manager/topbutton/language`,
	},
	{
		Tmpl: `manager/topbutton/source`,
	},
	{
		Tmpl: `manager/topbutton/bug-report`,
	},
}

func TopButtonRegister(topButton ...*Button) {
	topButtons.Add(-1, topButton...)
}

func TopButtonAdd(index int, topButton ...*Button) {
	topButtons.Add(index, topButton...)
}

//TopButtonRemove 删除元素
func TopButtonRemove(index int) {
	topButtons.Remove(index)
}

func TopButtonSearch(cb func(*Button) bool) int {
	return topButtons.Search(cb)
}

func TopButtonFindTmpl(tmpl string) int {
	return topButtons.FindTmpl(tmpl)
}

func TopButtonRemoveByTmpl(tmpl string) {
	topButtons.RemoveByTmpl(tmpl)
}

//TopButtonSet 设置元素
func TopButtonSet(index int, list ...*Button) {
	topButtons.Set(index, list...)
}

func TopButtonAll(_ echo.Context) Buttons {
	return topButtons
}
