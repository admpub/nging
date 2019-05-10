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
	"github.com/admpub/ip2region/binding/golang/ip2region"
	"github.com/webx-top/echo"
)

var (
	region   *ip2region.Ip2Region
	dictFile string
)

func init() {
	dictFile = echo.Wd() + echo.FilePathSeparator + `data` + echo.FilePathSeparator + `ip2region` + echo.FilePathSeparator + `ip2region.db`
}

func IP2Region(c echo.Context) (err error) {
	ip := c.Form(`ip`)
	if len(ip) > 0 {
		if region == nil {
			region, err = ip2region.New(dictFile)
			if err != nil {
				return
			}
		}
		var info ip2region.IpInfo
		info, err = region.MemorySearch(ip)
		if err != nil {
			return
		}
		c.Data().SetData(info)
	}
	return c.Render(`/tool/ip`, nil)
}
