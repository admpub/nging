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
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/config/cmder"
	"github.com/admpub/nging/v4/application/registry/dashboard"

	"github.com/nging-plugins/servermanager/application/registry"
)

func addLogCategory(logCategories *echo.KVList, k, v string) {
	logFilename, _ := config.FromFile().Log.LogFilename(k)
	if len(logFilename) > 0 {
		logFilename = filepath.Base(logFilename)
	}
	logCategories.Add(k, v, echo.KVOptHKV(`logFilename`, logFilename))
}

func Service(ctx echo.Context) error {
	logCategories := &echo.KVList{}
	addLogCategory(logCategories, log.DefaultLog.Category, ctx.T(`Nging日志`))
	if strings.Contains(config.FromFile().Log.LogFile(), `{category}`) {
		ctx.Set(`logWithCategory`, true)
		categories := config.FromFile().Log.LogCategories()
		for _, k := range categories {
			k = strings.SplitN(k, `,`, 2)[0]
			v := k
			switch k {
			case `db`:
				v = ctx.T(`SQL日志`)
			case `echo`:
				v = ctx.T(`Web框架日志`)
			default:
				v = ctx.T(`%s日志`, strings.Title(k))
			}
			addLogCategory(logCategories, k, v)
		}
	} else {
		ctx.Set(`logWithCategory`, false)
	}
	ctx.Set(`logCategories`, *logCategories)
	ctx.SetFunc(`HasService`, cmder.Has)
	ctx.SetFunc(`ServiceControls`, func() dashboard.Buttons {
		buttons := registry.ServiceControls
		buttons.Ready(ctx)
		return buttons
	})
	return ctx.Render(`server/service`, nil)
}
