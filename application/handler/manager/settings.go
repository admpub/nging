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

package manager

import (
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/registry/settings"
	"github.com/webx-top/echo"
)

func Settings(ctx echo.Context) error {
	//panic(echo.Dump(settings.ConfigAsStore(), false))
	var err error
	if ctx.IsPost() {
		err = configPost(ctx)
		if err != nil {
			return err
		}
		err = settings.RunHookPost(ctx)
		if err != nil {
			return err
		}
		if err == nil {
			config.DefaultConfig.SetDebug(ctx.Formx(`base[debug][value]`).Bool())
			config.DefaultConfig.ConfigInDB.Init()
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/manager/settings`))
		}
	}
	if _err := configGet(ctx); _err != nil {
		return _err
	}
	if _err := settings.RunHookGet(ctx); _err != nil {
		return _err
	}
	ret := handler.Err(ctx, err)
	ctx.Set(`config`, config.DefaultConfig)
	ctx.Set(`settings`, settings.Settings())
	return ctx.Render(`/manager/settings`, ret)
}
