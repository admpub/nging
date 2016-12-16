package handler

import (
	"github.com/admpub/nging/application/library/dbmanager"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	_ "github.com/admpub/nging/application/library/dbmanager/driver/mysql"
	"github.com/webx-top/echo"
)

func DbManager(ctx echo.Context) error {
	var err error
	mgr := dbmanager.New(ctx)
	driverName := ctx.Form(`driver`)
	operation := ctx.Form(`operation`)
	if len(driverName) > 0 {
		err = mgr.Run(driverName, operation)
	}
	ret := Err(ctx, err)
	driverList := []string{}
	for driverName := range driver.GetAll() {
		driverList = append(driverList, driverName)
	}
	ctx.Set(`driverList`, driverList)
	return ctx.Render(`db/index`, ret)
}
