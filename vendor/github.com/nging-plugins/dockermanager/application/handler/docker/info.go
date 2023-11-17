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

package docker

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
)

func Index(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	info, err := c.Info(ctx)
	if err != nil {
		return err
	}
	if ctx.Format() == echo.ContentTypeJSON {
		data := ctx.Data()
		switch ctx.Form(`prop`) {
		case `plugins.volume`:
			data.SetData(echo.H{`listData`: info.Plugins.Volume})
		case `plugins.log`:
			data.SetData(echo.H{`listData`: info.Plugins.Log})
		case `plugins.network`:
			data.SetData(echo.H{`listData`: info.Plugins.Network})
		case `plugins.authorization`:
			data.SetData(echo.H{`listData`: info.Plugins.Authorization})
		case `plugins`:
			data.SetData(info.Plugins)
		}
		return ctx.JSON(data)
	}
	ctx.Set(`info`, info)
	return ctx.Render(`docker/base/info/index`, handler.Err(ctx, err))
}
