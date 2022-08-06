package handler

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/nging-plugins/frpmanager/pkg/dbschema"
)

func ServerDashboard(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint()
	m := dbschema.NewNgingFrpServer(ctx)
	err := m.Get(nil, db.And(
		db.Cond{`id`: id},
		db.Cond{`disabled`: `N`},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return ctx.NewError(code.DataNotFound, `没有找到启用的配置信息`)
	}
	if m.DashboardPort > 0 {
		dashboardHost := m.DashboardAddr
		if m.DashboardAddr == `0.0.0.0` || len(m.DashboardAddr) == 0 {
			dashboardHost = ctx.Domain()
		}
		return ctx.Redirect(fmt.Sprintf(`http://%s:%d/`, dashboardHost, m.DashboardPort))
	}
	return ctx.NewError(code.Unsupported, `配置信息中未启用管理面板访问功能。如要启用，请将面板端口设为一个大于0的数值`)
}

func ClientDashboard(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint()
	m := dbschema.NewNgingFrpClient(ctx)
	err := m.Get(nil, db.And(
		db.Cond{`id`: id},
		db.Cond{`disabled`: `N`},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return ctx.NewError(code.DataNotFound, `没有找到启用的配置信息`)
	}
	if m.AdminPort > 0 {
		dashboardHost := m.AdminAddr
		if m.AdminAddr == `0.0.0.0` || len(m.AdminAddr) == 0 {
			dashboardHost = ctx.Domain()
		}
		return ctx.Redirect(fmt.Sprintf(`http://%s:%d/`, dashboardHost, m.AdminPort))
	}
	return ctx.NewError(code.Unsupported, `配置信息中未启用管理面板访问功能。如要启用，请将面板端口设为一个大于0的数值`)
}
