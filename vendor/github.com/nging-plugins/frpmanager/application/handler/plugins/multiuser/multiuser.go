package multiuser

import (
	"net/http"

	plugin "github.com/admpub/frp/pkg/plugin/server"
	"github.com/admpub/log"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/nging-plugins/frpmanager/application/dbschema"
	"github.com/nging-plugins/frpmanager/application/library/cmder"
	"github.com/nging-plugins/frpmanager/application/library/utils"
	"github.com/nging-plugins/frpmanager/application/model"
)

func Login(ctx echo.Context) error {
	op := ctx.Query(`op`)
	if op != `Login` {
		return echo.ErrForbidden
	}
	//ver := ctx.Query(`version`)
	var r plugin.Request
	var content plugin.LoginContent
	r.Content = &content
	if err := ctx.MustBind(&r); err != nil {
		return ctx.String(err.Error(), http.StatusBadRequest)
	}
	var res plugin.Response
	token := content.Metas["token"]
	serverIDStr := content.Metas["serverid"]
	if len(content.User) == 0 || len(token) == 0 {
		res.Reject = true
		res.RejectReason = "user or meta token can not be empty"
		return ctx.JSON(res)
	}
	//echo.Dump(echo.H{`user`: content.User, `token`: token, `serverid`: serverIDStr})
	serverID := param.AsUint(serverIDStr)
	if serverID < 1 {
		res.Reject = true
		if len(serverIDStr) == 0 {
			res.RejectReason = "meta serverid can not be empty"
		} else {
			res.RejectReason = "meta serverid can not less than 1"
		}
		return ctx.JSON(res)
	}
	m := model.NewFrpUser(ctx)
	err := m.CheckPasswd(serverID, content.User, token)
	if err != nil {
		res.Reject = true
		if err == db.ErrNoMoreRows {
			res.RejectReason = "invalid meta token"
		} else {
			res.RejectReason = err.Error()
		}
	} else {
		res.Unchange = true
	}
	return ctx.JSON(res)
}

func OnChangeBackendURL(ctx echo.Context) error {
	c := dbschema.NewNgingFrpServer(ctx)
	_, err := c.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`disabled`: `N`},
		mysql.FindInSet(`plugins`, `multiuser_login`),
	))
	if err != nil {
		return err
	}
	cm, err := cmder.GetServer()
	if err != nil {
		return nil
	}
	for _, row := range c.Objects() {
		err := utils.SaveConfigFile(row)
		if err != nil {
			log.Error(err)
			continue
		}
		id := param.AsString(row.Id)
		err = cm.RestartBy(id)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}
