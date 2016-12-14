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
	ctx.Set(`driverList`, echo.H{
		"mysql": "MySQL",
	})
	return ctx.Render(`db/index`, ret)
}
