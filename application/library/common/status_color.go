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

package common

import (
	"github.com/admpub/color"
)

var (
	bootstrapColors = map[StatusColor]string{
		`green`:  `success`,
		`red`:    `danger`,
		`yellow`: `warning`,
		`cyan`:   `info`,
	}
	terminalColors = map[StatusColor]func(string, ...interface{}){
		`green`:  color.Green,
		`red`:    color.Red,
		`yellow`: color.Yellow,
		`cyan`:   color.Cyan,
	}
)

// StatusColor 状态色
type StatusColor string

func (s StatusColor) String() string {
	return string(s)
}

// Bootstrap 前端框架 bootstrap css 状态样式
func (s StatusColor) Bootstrap() string {
	return bootstrapColors[s]
}

// Terminal 控制台样式
func (s StatusColor) Terminal() func(string, ...interface{}) {
	return terminalColors[s]
}

// HTTPStatusColor HTTP状态码相应颜色
func HTTPStatusColor(httpCode int) StatusColor {
	s := `green`
	switch {
	case httpCode >= 500:
		s = `red`
	case httpCode >= 400:
		s = `yellow`
	case httpCode >= 300:
		s = `cyan`
	}
	return StatusColor(s)
}
