package multiuser

import (
	"net/http"

	plugin "github.com/admpub/frp/pkg/plugin/server"
	"github.com/admpub/nging/v3/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
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
