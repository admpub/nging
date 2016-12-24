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
	var genURL func(string, ...string) string
	switch operation {
	case `login`:
		data := &driver.DbAuth{}
		ctx.Bind(data)
		if len(data.Username) == 0 {
			data.Username = `root`
		}
		if len(data.Host) == 0 {
			data.Host = `127.0.0.1`
		}
		auth.CopyFrom(data)
		ctx.Session().Set(`dbAuth`, data)
	case `logout`:
	default:
		if data, ok := ctx.Session().Get(`dbAuth`).(*driver.DbAuth); ok {
			if len(dbName) > 0 && len(data.Db) == 0 {
				data.Db = dbName
			}
			auth.CopyFrom(data)
			err = mgr.Run(auth.Driver, `login`)
			if err == nil {
				driverName = auth.Driver
				if operation == `` {
					operation = `listDb`
				}
				ctx.Set(`dbUsername`, auth.Username)
				ctx.Set(`dbHost`, auth.Host)
				genURL = func(op string, args ...string) string {
					if op == `` {
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
			}
		}
	}
	if genURL == nil {
		genURL = func(_ string, _ ...string) string {
			return ``
		}
	}
	mgr.GenURL = genURL
	ctx.SetFunc(`dbMgrURL`, genURL)
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
