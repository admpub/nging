package oauth2server

import (
	"context"
	"net/url"
	"strings"

	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/oauth2/v4"
	"github.com/admpub/oauth2/v4/models"
	"github.com/coscms/oauth2s"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/defaults"
)

var (
	DefaultAppClient oauth2.ClientStore = NewClientStore(nil)
	Default                             = oauth2s.NewConfig()
)

var requestFormDataCacheKey = `backend.` + oauth2s.RequestFormDataCacheKey

func NewClientStore(fn func(ctx echo.Context, appID string) (oauth2.ClientInfo, error)) oauth2.ClientStore {
	if fn == nil {
		fn = getByAppID
	}
	return &AppClient{GetByAppID: fn}
}

func getByAppID(ctx echo.Context, appID string) (oauth2.ClientInfo, error) {
	openM := model.NewOAuthApp(ctx)
	err := openM.GetByAppID(appID)
	if err != nil {
		return nil, err
	}
	scope := ctx.Form("scope")
	var scopes []string
	if len(scope) > 0 {
		scopes = strings.Split(scope, " ")
	}
	redirectURI := ctx.Form("redirect_uri")
	u, err := url.Parse(redirectURI)
	if err != nil {
		return nil, ctx.NewError(code.InvalidParameter, `参数 %s 不是有效的网址`, `redirect_uri`)
	}
	domain := com.SplitHost(u.Host)
	err = openM.VerifyApp(domain, scopes)
	if err != nil {
		return nil, err
	}
	clientInfo := &models.Client{
		ID:     openM.AppId,
		Secret: openM.AppSecret,
		// Domain: u.Scheme + `://` + u.Host,
		// UserID: ``,
	}
	return clientInfo, nil
}

type AppClient struct {
	GetByAppID func(ctx echo.Context, appID string) (oauth2.ClientInfo, error)
}

func (a *AppClient) GetByID(stdCtx context.Context, appID string) (oauth2.ClientInfo, error) {
	ctx := defaults.MustGetContext(stdCtx)
	if Debug {
		println(`[AppClient.GetByID.Forms]:`, echo.Dump(ctx.Forms(), false))
	}
	return a.GetByAppID(ctx, appID)
}
