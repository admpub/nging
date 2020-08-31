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

package tool

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/library/ip2region"
)

func IP2Region(c echo.Context) error {
	ip := c.Form(`ip`)
	if len(ip) > 0 {
		info, err := ip2region.IPInfo(ip)
		if err != nil {
			return err
		}
		c.Data().SetData(info)
	} else {
		c.Request().Form().Set(`ip`, c.RealIP())
	}
	return c.Render(`/tool/ip`, nil)
}
