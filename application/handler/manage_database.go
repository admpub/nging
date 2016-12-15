package handler

import (
	"github.com/admpub/caddyui/application/library/dbmanager"
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
	for driverName := range dbmanager.GetAll() {
		driverList = append(driverList, driverName)
	}
	ctx.Set(`driverList`, driverList)
	return ctx.Render(`db/index`, ret)
}
