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

package tool

import (
	"regexp"

	"github.com/admpub/nging/v5/application/library/fileupdater/listener"
	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/nging/v5/application/model/file/storer"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"
)

var replaceURLRegex = regexp.MustCompile(`^(?s)http[s]?://(.+)/$`)

func ReplaceURL(c echo.Context) (err error) {
	storerInfo := storer.Get()
	toURL := c.Formx(`to`).String()
	m := model.NewCloudStorage(c)
	if len(storerInfo.ID) > 0 {
		err = m.NgingCloudStorage.Get(nil, `id`, storerInfo.ID)
		if len(toURL) == 0 && err == nil {
			toURL = m.BaseURL() + `/`
			c.Request().Form().Set(`to`, toURL)
		}
	}
	if c.IsPost() {
		data := c.Data()
		fromURL := c.Formx("from").String()
		var total int64
		if len(fromURL) == 0 {
			err = c.E(`旧网址不能为空`)
		} else if !replaceURLRegex.MatchString(fromURL) {
			err = c.E(`要查找的旧网址无效`)
		} else if !replaceURLRegex.MatchString(toURL) {
			err = c.E(`新网址无效`)
		} else if fromURL != toURL {
			for proj, tables := range listener.UpdaterInfos {
				_ = proj
				for table, fields := range tables {
					for field, info := range fields {
						var affected int64
						if table == `nging_config` && field == `value` {
							affected, err = mysql.ReplacePrefix(0, table, field, fromURL, toURL, true, "`type` IN('html','list','json')")
							if err != nil {
								goto END
							}
							total += affected
							affected, err = mysql.ReplacePrefix(0, table, field, fromURL, toURL, false, "`type` NOT IN('html','list','json')")
							if err != nil {
								goto END
							}
							total += affected
							continue
						}
						if info.Embedded || len(info.Seperator) > 0 {
							affected, err = mysql.ReplacePrefix(0, table, field, fromURL, toURL, true)
							if err != nil {
								goto END
							}
							total += affected
							continue
						}
						affected, err = mysql.ReplacePrefix(0, table, field, fromURL, toURL, false)
						if err != nil {
							goto END
						}
						total += affected
					}
				}
			}
		}

	END:
		data.SetData(echo.H{`total`: total})
		if err != nil {
			data.SetError(err)
		}
		return c.JSON(data)
	}
	c.Set(`updaterInfos`, listener.UpdaterInfos)
	return c.Render(`tool/replaceurl`, nil)
}
