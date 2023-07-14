package backend

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/model"
)

func WebAuthn(ctx echo.Context) error {
	user := handler.User(ctx)
	if user == nil {
		return ctx.NewError(code.Unauthenticated, `请先登录`)
	}
	m := model.NewUserU2F(ctx)
	err := m.ListPageByType(user.Id, `webauthn`, 1)
	if err != nil {
		return err
	}
	ctx.Set(`listData`, m.Objects())
	ctx.Set(`activeSafeItem`, `webauthn`)
	ctx.Set(`safeItems`, model.SafeItems.Slice())
	return ctx.Render(`webauthn/user`, handler.Err(ctx, err))
}
