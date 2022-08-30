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
	"regexp"

	"github.com/admpub/regexp2"
	"github.com/webx-top/echo"
)

type RegexpForm struct {
	Src    string
	Regexp string
	Type   string
	Error  string
	Result interface{}
}

func RegexpTest(c echo.Context) error {
	data := &RegexpForm{}
	c.Bind(data)
	if c.IsPost() {
		if data.Type == `regexp2` {
			if cr, err := regexp2.Compile(data.Regexp, 0); err != nil {
				data.Error = err.Error()
			} else if len(data.Src) > 0 {
				matches, err := cr.FindStringMatch(data.Src)
				if err != nil {
					data.Error = err.Error()
				} else if matches != nil {
					data.Result = matches.Slice2()
				}
			}
		} else {
			if cr, err := regexp.Compile(data.Regexp); err != nil {
				data.Error = err.Error()
			} else if len(data.Src) > 0 {
				data.Result = cr.FindAllStringSubmatch(data.Src, -1)
			}
		}
		return c.JSON(data)
	}
	data.Type = c.Query(`type`, `regexp`)
	return c.Render(`collector/regexp_test`, data)
}
