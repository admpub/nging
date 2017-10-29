/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package database

import (
	"strings"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/admpub/nging/application/middleware"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func init() {
	handler.Register(func(e *echo.Echo) {
		e.Route(`GET,POST`, `/db/account`, AccountIndex, middleware.AuthCheck)
		e.Route(`GET,POST`, `/db/account_add`, AccountAdd, middleware.AuthCheck)
		e.Route(`GET,POST`, `/db/account_edit`, AccountEdit, middleware.AuthCheck)
		e.Route(`GET,POST`, `/db/account_delete`, AccountDelete, middleware.AuthCheck)
	})
}

func AccountIndex(ctx echo.Context) error {
	user := handler.User(ctx)
	m := model.NewDbAccount(ctx)
	page, size, totalRows, p := handler.PagingWithPagination(ctx)
	cond := db.Cond{
		`uid`: user.Id,
	}
	cnt, err := m.List(nil, func(r db.Result) db.Result {
		return r
	}, page, size, cond)
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
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(`/db/account`)
			}
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
		err = ctx.MustBind(m.DbAccount, func(k string, v []string) (string, []string) {
			switch strings.ToLower(k) {
			case `created`, `uid`: //禁止修改创建时间
				return ``, v
			}
			return k, v
		})

		if err == nil {
			m.Id = id
			err = m.Edit(nil, cond)
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(`/db/account`)
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.DbAccount, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
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

	return ctx.Redirect(`/db/account`)
}
