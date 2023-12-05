package oauth2server

import (
	"context"
	"net/http"
	"strings"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/middleware"
	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/oauth2/v4/errors"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
)

// PasswordAuthorizationHandler 密码认证处理
func PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	err = errors.ErrAccessDenied
	return
	/*/ 后台的账号登录安全要求非常高，不开放此功能
	c := defaults.MustGetContext(ctx)
	m := model.NewUser(c)
	c.Request().Form().Set(`user`, username)
	c.Request().Form().Set(`pass`, password)
	err = middleware.Auth(c)
	if err != nil {
		log.Debug(`oauth2server.PasswordAuthorizationHandler: `, err.Error())
		return
	}
	userID = param.AsString(m.Id)
	return
	//*/
}

// UserAuthorizeHandler 查询用户授权信息
func UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	ctx := defaults.MustGetContext(r.Context())
	if Debug {
		println(`[UserAuthorizeHandler.Forms]:`, echo.Dump(ctx.Forms(), false))
	}

	user, ok := ctx.Session().Get(`user`).(*dbschema.NgingUser)
	if !ok || user == nil {
		// 没有登录时记录当前的提交数据
		ctx.Session().Set(requestFormDataCacheKey, ctx.Forms())
		err = ctx.Redirect(handler.URLFor(RoutePrefix + `/login`))
		return
	}
	var need bool
	need, err = middleware.TwoFactorAuth(ctx, func() error {
		ctx.Session().Set(requestFormDataCacheKey, ctx.Forms())
		return nil
	})
	if need {
		return
	}
	clientID := ctx.Form("client_id")
	var scopes []string
	scope := ctx.Form("scope")
	if len(scope) > 0 {
		scopes = strings.Split(scope, " ")
	}
	m := model.NewOAuthAgree(ctx)
	var agreed bool
	agreed, err = m.IsAgreed(user.Id, clientID, scopes)
	if err != nil {
		return
	}
	if !agreed {
		// 没有授权时记录当前的提交数据
		ctx.Session().Set(requestFormDataCacheKey, ctx.Forms())
		err = ctx.Redirect(handler.URLFor(RoutePrefix + `/auth`))
		return
	}
	userID = param.AsString(user.Id)
	ctx.Session().Delete(requestFormDataCacheKey)
	ctx.Session().Save()
	return
}

type AgreedScopeStorer interface {
	IsAgreed(userID uint, appID string, scopes []string) (bool, error)
	Save(userID uint, appID string, scopes []string) error
}
