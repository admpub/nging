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
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/dbmanager"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	_ "github.com/admpub/nging/application/library/dbmanager/driver/mysql" //mysql
	"github.com/admpub/nging/application/middleware"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func init() {
	handler.Register(func(e *echo.Echo) {
		e.Route(`GET,POST`, `/db`, Manager, middleware.AuthCheck)
	})
}

func Manager(ctx echo.Context) error {
	user := handler.User(ctx)
	m := model.NewDbAccount(ctx)
	var err error
	auth := &driver.DbAuth{}
	mgr := dbmanager.New(ctx, auth)
	driverName := ctx.Form(`driver`)
	operation := ctx.Form(`operation`)
	var accountId uint
	if user != nil {
		accountId = ctx.Formx(`accountId`).Uint()
		if accountId > 0 {
			err = m.Get(nil, db.And(db.Cond{`id`: accountId}, db.Cond{`uid`: user.Id}))
		}
	}
	var genURL func(string, ...string) string
	switch operation {
	case `login`:
		if ctx.IsPost() {
			data := &driver.DbAuth{}
			ctx.Bind(data)
			if len(data.Username) == 0 {
				data.Username = `root`
			}
			if len(data.Host) == 0 {
				data.Host = `127.0.0.1`
			}
			auth.CopyFrom(data)
			if accountId > 0 && ctx.Form(`remember`) == `1` {
				m.Engine = auth.Driver
				m.Host = auth.Host
				m.User = auth.Username
				m.Password = auth.Password
				m.Name = auth.Db
				if err == db.ErrNoMoreRows {
					m.Uid = user.Id
					_, err = m.Add()
				} else {
					err = m.Edit(nil, db.Cond{`id`: accountId})
				}
			}
			ctx.Session().Set(`dbAuth`, auth)
		}

	case `logout`:
		//pass

	default:
		if accountId > 0 {
			auth.Driver = m.Engine
			auth.Username = m.User
			auth.Password = m.Password
			auth.Host = m.Host
			auth.Db = m.Name
			ctx.Session().Set(`dbAuth`, auth)
			err = mgr.Run(auth.Driver, `login`)
		} else {
			if data, exists := ctx.Session().Get(`dbAuth`).(*driver.DbAuth); exists {
				auth.CopyFrom(data)
				err = mgr.Run(auth.Driver, `login`)
			}
		}
		if err == nil {
			driverName = auth.Driver
			if len(operation) == 0 {
				operation = `listDb`
			}
			ctx.Set(`dbUsername`, auth.Username)
			ctx.Set(`dbHost`, auth.Host)
			genURL = func(op string, args ...string) string {
				if len(op) == 0 {
					op = operation
				}
				var p string
				switch len(args) {
				case 2:
					p += `&table=` + args[1]
					fallthrough
				case 1:
					p += `&db=` + args[0]
				}
				return `/db?driver=` + driverName + `&username=` + auth.Username + `&operation=` + op + p
			}
			defer mgr.Run(auth.Driver, `logout`)
		} else {
			if err != nil {
				handler.SendFail(ctx, err.Error())
				err = nil
			}
			driverName = ``
		}
	}
	if genURL == nil {
		genURL = func(_ string, _ ...string) string {
			return ``
		}
	}
	mgr.GenURL = genURL
	ctx.SetFunc(`dbMgrURL`, genURL)
	ctx.Set(`operation`, operation)
	if len(driverName) > 0 {
		ctx.Set(`driver`, driverName)
		if err == nil {
			err = mgr.Run(driverName, operation)
		}
		if err == nil {
			switch operation {
			case `login`:
				mgr.Run(auth.Driver, `logout`)
				return ctx.Redirect(`/db`)
			case `logout`:
				mgr.Run(auth.Driver, `logout`)
				ctx.Session().Delete(`dbAuth`)
			default:
				return err
			}
		} else {
			if operation != `login` {
				return err
			}
		}

	}
	ret := handler.Err(ctx, err)
	driverList := map[string]string{}
	for driverName, driver := range driver.GetAll() {
		driverList[driverName] = driver.Name()
	}

	if accountId > 0 {
		ctx.Request().Form().Set(`db`, m.Name)
		ctx.Request().Form().Set(`driver`, m.Engine)
		ctx.Request().Form().Set(`username`, m.User)
		ctx.Request().Form().Set(`password`, m.Password)
		ctx.Request().Form().Set(`host`, m.Host)
	}

	ctx.Set(`driverList`, driverList)
	return ctx.Render(`db/index`, ret)
}
