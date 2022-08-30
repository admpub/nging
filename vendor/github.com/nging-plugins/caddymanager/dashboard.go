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

package caddymanager

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/registry/dashboard"
	"github.com/nging-plugins/caddymanager/application/model"
)

func RegisterDashboard(dd *dashboard.Dashboards) {
	dd.Backend.Cards.Add(-1,
		(&dashboard.Card{
			IconName:  `fa-sitemap`,
			IconColor: `primary`,
			Short:     `SITES`,
			Name:      `网站数量`,
			Summary:   ``,
		}).SetContentGenerator(func(ctx echo.Context) interface{} {
			//网站统计
			vhostMdl := model.NewVhost(ctx)
			vhostCount, _ := vhostMdl.Count(nil)
			return vhostCount
		}),
	)

}
