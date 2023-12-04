package user

import (
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/backend/oauth2client"
	"github.com/admpub/nging/v5/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

type oAuthProvider struct {
	On        bool
	Name      string
	LoginURL  string
	Scopes    []string
	IconClass string `json:",omitempty"`
	IconImage string `json:",omitempty"`
	WrapClass string `json:",omitempty"`
	Binded    *dbschema.NgingUserOauth
}

func oAuth(ctx echo.Context) error {
	user := handler.User(ctx)
	if user == nil {
		return ctx.NewError(code.Unauthenticated, `请先登录`)
	}
	m := model.NewUserOAuth(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`uid`, user.Id)
	_, err := m.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, 0, -1, cond.And())
	if err != nil {
		return err
	}
	accIdx := map[string]int{}
	accounts := oauth2client.GetOAuthAccounts()
	oAuthProviders := make([]oAuthProvider, len(accounts))
	for index, account := range accounts {
		accIdx[account.Name] = index
		provider := oAuthProvider{
			On:       account.On,
			Name:     account.Name,
			LoginURL: account.LoginURL,
			Scopes:   make([]string, len(account.Scopes)),
		}
		copy(provider.Scopes, account.Scopes)
		provider.IconImage = account.Extra.String(`iconImage`)
		provider.IconClass = account.Extra.String(`iconClass`)
		provider.WrapClass = account.Extra.String(`wrapClass`)
		title := account.Extra.String(`title`)
		if len(title) > 0 {
			provider.Name = title
		} else {
			provider.Name = com.Title(provider.Name)
		}
		oAuthProviders[index] = provider
	}
	for _, row := range m.Objects() {
		idx, ok := accIdx[row.Type]
		if !ok {
			continue
		}
		rowCopy := *row
		rowCopy.AccessToken = `[HIDE]`
		rowCopy.RefreshToken = `[HIDE]`
		oAuthProviders[idx].Binded = &rowCopy
	}
	ctx.Internal().Set(`enabledOAuthAccounts`, accounts)
	ctx.Set(`activeSafeItem`, `oauth`)
	ctx.Set(`safeItems`, model.SafeItems.Slice())
	ctx.Set(`oAuthProviders`, oAuthProviders)
	return ctx.Render(`user/oauth`, handler.Err(ctx, err))
}

func oAuthDelete(ctx echo.Context) error {
	user := handler.User(ctx)
	if user == nil {
		return ctx.NewError(code.Unauthenticated, `请先登录`)
	}
	id := ctx.Paramx(`id`).Uint()
	if id < 1 {
		return ctx.NewError(code.InvalidParameter, `参数无效`)
	}
	m := model.NewUserOAuth(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`uid`, user.Id)
	cond.AddKV(`id`, id)
	affected, err := m.Deletex(nil, cond.And())
	if err == nil {
		if affected == 0 {
			handler.SendFail(ctx, ctx.T(`没有找到可以删除的数据`))
		} else {
			handler.SendOk(ctx, ctx.T(`操作成功`))
		}
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/user/oauth`))
}
