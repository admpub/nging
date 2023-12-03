package oauth2server

import (
	"encoding/gob"
	"net/http"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/formfilter"

	"github.com/admpub/nging/v5/application/dbschema"
	handlerIndex "github.com/admpub/nging/v5/application/handler/index"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/model"
)

func init() {
	gob.Register(&dbschema.NgingOauthApp{})
}

var (
	// Debug 调试模式
	Debug bool
)

var authCodeFormFilter = formfilter.Build(formfilter.SplitValues(`scope`, ` `))

// 获取授权码(response_type=code)或token(response_type=token) (首先进入执行)
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := defaults.MustGetContext(r.Context())
	applyCachedFormData(ctx)

	if Debug {
		println(`[authorizeHandler.Forms]:`, echo.Dump(ctx.Forms(), false))
	}

	// 调用 UserAuthorizeHandler()
	if err := Default.Server().HandleAuthorizeRequest(w, ctx.Request().StdRequest().WithContext(ctx)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := defaults.MustGetContext(r.Context())

	applyCachedFormData(ctx)

	if Debug {
		println(`[loginHandler.Forms]:`, echo.Dump(ctx.Forms(), false))
	}

	ctx.Request().Form().Set(`next`, RoutePrefix+`/auth`)
	err := handlerIndex.Login(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// 取回被缓存的数据(在访问authorizeHandler调用UserAuthorizeHandler时缓存下来的)
func applyCachedFormData(ctx echo.Context) {
	if values, ok := ctx.Session().Get(requestFormDataCacheKey).(map[string][]string); ok && values != nil {
		common.CopyFormDataFrom(ctx, values)
	}

	//ctx.Session().Delete(requestFormDataCacheKey)
}

// 用户授权页面
func authHandler(w http.ResponseWriter, r *http.Request) {
	ctx := defaults.MustGetContext(r.Context())
	user, ok := ctx.Session().Get(`user`).(*dbschema.NgingUser)
	var err error
	if !ok || user == nil {
		err = ctx.Redirect(RoutePrefix + `/login`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	applyCachedFormData(ctx)
	if Debug {
		println(`[authHandler.Forms]:`, echo.Dump(ctx.Forms(), false))
	}

	// 获取OpenApp数据
	clientID := ctx.Form("client_id")
	openApp := model.NewOAuthApp(ctx)
	err = openApp.GetByAppID(clientID)
	if err != nil {
		ctx.Session().Delete(requestFormDataCacheKey)
		http.Error(w, ctx.T(`获取应用信息失败`)+`: `+err.Error(), http.StatusBadRequest)
		return
	}

	scope := ctx.Form(`scope`)
	var scopes []string
	if len(scope) > 0 {
		scopes = strings.Split(scope, " ")
	}
	err = openApp.VerifyApp(``, scopes)
	if err != nil {
		ctx.Session().Delete(requestFormDataCacheKey)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := model.NewOAuthAgree(ctx)
	var agreed bool
	agreed, err = m.IsAgreed(user.Id, clientID, scopes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if agreed {
		err = ctx.Redirect(RoutePrefix + `/authorize`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	if ctx.IsPost() {
		accept := ctx.Formx(`accept`).Bool()
		if accept {
			err = m.Save(user.Id, clientID, scopes)
			if err == nil {
				err = ctx.Redirect(RoutePrefix + `/authorize`)
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	ctx.SetFunc(`getOpenApp`, func() dbschema.NgingOauthApp {
		v := *openApp.NgingOauthApp
		v.AppSecret = ``
		return v
	})
	ctx.Set(`scopes`, []echo.H{})
	err = ctx.Render(`oauth2server/auth`, common.Err(ctx, err))
	if err != nil {
		ctx.Session().Delete(requestFormDataCacheKey)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := defaults.MustGetContext(r.Context())
	err := handlerIndex.Logout(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ctx.Session().Delete(requestFormDataCacheKey)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := defaults.MustGetContext(r.Context())
	if Debug {
		println(`[tokenHandler.Forms]:`, echo.Dump(ctx.Forms(), false))
	}
	err := Default.Server().HandleTokenRequest(w, ctx.Request().StdRequest().WithContext(ctx))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
