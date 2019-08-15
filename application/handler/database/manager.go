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
	"github.com/admpub/nging/application/library/dbmanager/driver/mysql"   //mysql
	_ "github.com/admpub/nging/application/library/dbmanager/driver/redis" //redis
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

var defaultGenURL = func(_ string, _ ...string) string {
	return ``
}

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
	genURL := defaultGenURL
	switch operation {
	case `login`:
		if err = login(mgr, accountID, m, user); err != nil {
			return err
		}

	case `logout`:
		//pass

	default:
		var signedIn bool
		err, signedIn = authentication(mgr, accountID, m)
		ctx.Set(`signedIn`, signedIn)
		ctx.Set(`dbUsername`, auth.Username)
		ctx.Set(`dbHost`, auth.Host)
		if signedIn {
			driverName = auth.Driver
			if len(operation) == 0 {
				operation = `listDb`
			}
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
	mgr.SetURLGenerator(genURL)
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
	ctx.Set(`charsetList`, mysql.Charsets)
	return ctx.Render(`db/index`, ret)
}
