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

func NewGlobalHead(tmpl string, handle ...func(echo.Context) error) *GlobalHead {
	var content func(echo.Context) error
	if len(handle) > 0 {
		content = handle[0]
	}
	return &GlobalHead{Tmpl: tmpl, content: content}
}

type GlobalHead = Tmplx
type GlobalHeads = Tmplxs

var globalHeads = GlobalHeads{}

func GlobalHeadRegister(footer ...*GlobalHead) {
	globalHeads.Add(-1, footer...)
}

func GlobalHeadAdd(index int, footer ...*GlobalHead) {
	globalHeads.Add(index, footer...)
}

//GlobalHeadRemove 删除元素
func GlobalHeadRemove(index int) {
	globalHeads.Remove(index)
}

//GlobalHeadSet 设置元素
func GlobalHeadSet(index int, list ...*GlobalHead) {
	globalHeads.Set(index, list...)
}

func GlobalHeadAll(_ echo.Context) GlobalHeads {
	return globalHeads
}
