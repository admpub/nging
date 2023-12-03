package oauth2server

import (
	"strings"

	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/oauth2/v4"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
)

// RefreshingScopeHandler check the scope of the refreshing token
func RefreshingScopeHandler(tgr *oauth2.TokenGenerateRequest, oldScope string) (allowed bool, err error) {
	ctx := defaults.MustGetContext(tgr.Request.Context())
	if Debug {
		println(`[refreshingScopeHandler.Forms]:`, echo.Dump(ctx.Forms(), false))
	}
	var token oauth2.TokenInfo
	token, err = Default.Server().ValidationBearerToken(tgr.Request)
	if err != nil {
		return
	}
	clientID := token.GetClientID()
	userID := param.AsUint(token.GetUserID())
	scope := token.GetScope()
	var scopes []string
	if len(scope) > 0 {
		scopes = strings.Split(scope, " ")
	}
	m := model.NewOAuthAgree(ctx)
	allowed, err = m.IsAgreed(userID, clientID, scopes)
	return
}

// RefreshingValidationHandler check if refresh_token is still valid. eg no revocation or other
func RefreshingValidationHandler(ti oauth2.TokenInfo) (allowed bool, err error) {
	allowed = true
	return
}
