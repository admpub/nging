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
	"time"

	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func LoginLog(ctx echo.Context) error {
	cond := db.NewCompounds()
	ownerType := ctx.Form(`ownerType`)
	if len(ownerType) > 0 {
		cond.AddKV(`owner_type`, ownerType)
	}
	q := ctx.Formx(`username`, ctx.Form(`q`)).String()
	if len(q) > 0 {
		cond.AddKV(`username`, q)
	}
	success := ctx.Formx(`success`).String()
	if len(success) > 0 {
		cond.AddKV(`success`, success)
	}
	m := model.NewLoginLog(ctx)
	list, err := m.ListPage(cond, `-created`)
	if err != nil {
		return err
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`/manager/login_log`, handler.Err(ctx, err))
}

func LoginLogDelete(ctx echo.Context) error {
	m := model.NewLoginLog(ctx)
	oneMonthAgo := time.Now().Local().Unix() - 30*86400 // 删除30天之前的数据
	err := m.Delete(nil, db.Cond{`created`: db.Lt(oneMonthAgo)})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/manager/login_log`))
}
