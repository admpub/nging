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
	"github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/dbmanager"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	_ "github.com/admpub/nging/application/library/dbmanager/driver/mysql" //mysql
	_ "github.com/admpub/nging/application/library/dbmanager/driver/redis" //redis
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func Manager(ctx echo.Context) error {
	user := handler.User(ctx)
	m := model.NewDbAccount(ctx)
	var err error
	auth := &driver.DbAuth{}
	mgr := dbmanager.New(ctx, auth)
	driverName := ctx.Form(`driver`)
	operation := ctx.Form(`operation`)
	var accountID uint
	if user != nil {
		accountID = ctx.Formx(`accountId`).Uint()
		if accountID > 0 {
			err = m.Get(nil, db.And(
				db.Cond{`id`: accountID},
				db.Cond{`uid`: user.Id},
			))
			if err != nil && db.ErrNoMoreRows != err {
				log.Error(err)
			}
			accountID = m.Id
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
			if ctx.Form(`remember`) == `1` {
				m.Title = auth.Driver + `://` + auth.Username + `@` + auth.Host + `/` + auth.Db
				m.Engine = auth.Driver
				m.Host = auth.Host
				m.User = auth.Username
				m.Password = auth.Password
				m.Name = auth.Db
				if accountID < 1 || err == db.ErrNoMoreRows {
					m.Uid = user.Id
					_, err = m.Add()
				} else {
					err = m.Edit(accountID, nil, db.Cond{`id`: accountID})
				}
			}
			ctx.Session().Set(`dbAuth`, auth)
		}

	case `logout`:
		//pass

	default:
		if accountID > 0 {
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
			ctx.Set(`signedIn`, true)
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
				return ctx.Redirect(handler.URLFor(`/db`))
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

	if accountID > 0 {
		ctx.Request().Form().Set(`db`, m.Name)
		ctx.Request().Form().Set(`driver`, m.Engine)
		ctx.Request().Form().Set(`username`, m.User)
		ctx.Request().Form().Set(`password`, m.Password)
		ctx.Request().Form().Set(`host`, m.Host)
	}

	ctx.Set(`driverList`, driverList)
	ctx.Set(`dbType`, ctx.T(`数据库`))
	return ctx.Render(`db/index`, ret)
}
