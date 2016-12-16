package handler

import (
	"github.com/admpub/nging/application/library/dbmanager"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	_ "github.com/admpub/nging/application/library/dbmanager/driver/mysql"
	"github.com/webx-top/echo"
)

func DbManager(ctx echo.Context) error {
	var err error
	auth := &driver.DbAuth{}
	mgr := dbmanager.New(ctx, auth)
	driverName := ctx.Form(`driver`)
	operation := ctx.Form(`operation`)
	dbName := ctx.Form(`db`)
	if data, ok := ctx.Session().Get(`dbAuth`).(*driver.DbAuth); ok {
		if len(dbName) > 0 && len(data.Db) == 0 {
			data.Db = dbName
		}
		auth.CopyFrom(data)
		err = mgr.Run(auth.Driver, `login`)
		if err == nil {
			driverName = auth.Driver
			operation = `listTable`
		}
	}
	if len(driverName) > 0 {
		switch operation {
		case `login`:
			data := &driver.DbAuth{}
			ctx.Bind(data)
			if len(data.Username) == 0 {
				data.Username = `root`
			}
			if len(data.Host) == 0 {
				data.Host = `127.0.0.1:3306`
			}
			auth.CopyFrom(data)
			ctx.Session().Set(`dbAuth`, data)
		}
		ctx.Set(`driver`, driverName)
		if err == nil {
			err = mgr.Run(driverName, operation)
		}
		if err == nil {
			switch operation {
			case `login`:
			case `logout`:
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
	ret := Err(ctx, err)
	driverList := []string{}
	for driverName := range driver.GetAll() {
		driverList = append(driverList, driverName)
	}
	ctx.Set(`driverList`, driverList)
	return ctx.Render(`db/index`, ret)
}
