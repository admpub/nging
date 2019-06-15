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
package database

import (
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/db`, func(g echo.RouteRegister) {
		e := handler.Echo()
		g.Route(`GET,POST`, ``, Manager)
		g.Route(`GET,POST`, `/account`, e.MetaHandler(echo.H{`name`: `账号列表`}, AccountIndex))
		g.Route(`GET,POST`, `/account_add`, e.MetaHandler(echo.H{`name`: `添加账号`}, AccountAdd))
		g.Route(`GET,POST`, `/account_edit`, e.MetaHandler(echo.H{`name`: `修改账号`}, AccountEdit))
		g.Route(`GET,POST`, `/account_delete`, e.MetaHandler(echo.H{`name`: `删除账号`}, AccountDelete))
	})
}

func AccountIndex(ctx echo.Context) error {
	user := handler.User(ctx)
	m := model.NewDbAccount(ctx)
	page, size, totalRows, p := handler.PagingWithPagination(ctx)
	cond := db.Compounds{
		db.Cond{`uid`: user.Id},
	}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`name`, db.Like(`%`+q+`%`))
	}
	cnt, err := m.List(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, page, size, cond.And())
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	ret := handler.Err(ctx, err)
	driverList := map[string]string{}
	for driverName, driver := range driver.GetAll() {
		driverList[driverName] = driver.Name()
	}
	ctx.Set(`driverList`, driverList)
	ctx.Set(`pagination`, p)
	ctx.Set(`listData`, m.Objects())
	ctx.Set(`activeURL`, `/db/account`)
	return ctx.Render(`db/account`, ret)
}

func AccountAdd(ctx echo.Context) error {
	user := handler.User(ctx)
	var err error
	if ctx.IsPost() {
		m := model.NewDbAccount(ctx)
		err = ctx.MustBind(m.DbAccount)
		if err == nil {
			m.Uid = user.Id
			_, err = m.Add()
			if err == nil {
				if ctx.IsAjax() {
					data := ctx.Data().SetInfo(ctx.T(`数据库账号成功`)).SetData(m.DbAccount)
					return ctx.JSON(data)
				}
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/db/account`))
			}
		}
		if err != nil && ctx.IsAjax() {
			return ctx.JSON(ctx.Data().SetError(err))
		}
	}
	ret := handler.Err(ctx, err)
	driverList := map[string]string{}
	for driverName, driver := range driver.GetAll() {
		driverList[driverName] = driver.Name()
	}
	ctx.Set(`driverList`, driverList)
	ctx.Set(`activeURL`, `/db/account_add`)
	return ctx.Render(`db/account_edit`, ret)
}

func AccountEdit(ctx echo.Context) error {
	user := handler.User(ctx)
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewDbAccount(ctx)
	cond := db.And(db.Cond{`id`: id}, db.Cond{`uid`: user.Id})
	err = m.Get(nil, cond)
	if ctx.IsPost() {
		err = ctx.MustBind(m.DbAccount, echo.ExcludeFieldName(`created`, `uid`))

		if err == nil {
			m.Id = id
			err = m.Edit(id, nil, cond)
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/db/account`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.DbAccount, ``, echo.LowerCaseFirstLetter)
	}

	ret := handler.Err(ctx, err)
	driverList := map[string]string{}
	for driverName, driver := range driver.GetAll() {
		driverList[driverName] = driver.Name()
	}
	ctx.Set(`driverList`, driverList)
	ctx.Set(`activeURL`, `/db/account`)
	return ctx.Render(`db/account_edit`, ret)
}

func AccountDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewDbAccount(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/db/account`))
}
