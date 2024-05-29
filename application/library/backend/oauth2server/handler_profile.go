package oauth2server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
)

func profileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := defaults.MustGetContext(r.Context())
	if Debug {
		println(`[profileHandler.header.Authorization]:`, r.Header.Get("Authorization"))
		println(`[profileHandler.Forms]:`, echo.Dump(ctx.Forms(), false))
		println(`[profileHandler.Form]:`, echo.Dump(r.Form, false))
	}
	token, err := Default.Server().ValidationBearerToken(r)
	if err != nil {
		log.Errorf(`failed to oauth2server.ValidationBearerToken: %v`, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	err = m.Get(userID, clientID)
	if err != nil {
		http.Error(w, ctx.T(`此账号尚未授权访问权限`), http.StatusForbidden)
		return
	}
	var accepted []string
	if len(m.Scopes) > 0 {
		accepted = strings.Split(m.Scopes, `,`)
	}
	var noAccepts []string
	for _, scope := range scopes {
		if !com.InSlice(scope, accepted) {
			noAccepts = append(noAccepts, scope)
		}
	}
	if len(noAccepts) > 0 {
		http.Error(w, ctx.T(`没有访问该用户 %s 的权限`, strings.Join(noAccepts, `,`)), http.StatusForbidden)
		return
	}

	userM := model.NewUser(ctx)
	err = userM.Get(nil, `id`, userID)
	if err != nil {
		if err == db.ErrNoMoreRows {
			http.Error(w, ctx.T(`用户不存在`), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !strings.Contains(userM.Avatar, `://`) {
		siteURL := common.BackendURL(ctx)
		if len(userM.Avatar) == 0 {
			userM.Avatar = siteURL + `/public/assets/backend/images/user_128.png`
		} else {
			userM.Avatar = siteURL + `/` + strings.TrimPrefix(userM.Avatar, `/`)
		}
	}
	data := map[string]interface{}{
		"expires_in":  int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":   clientID,
		"id":          userID,
		"name":        userM.Username,
		"avatar":      userM.Avatar,
		"email":       userM.Email,
		"description": ``,
	}
	e := json.NewEncoder(w)
	if Debug {
		e.SetIndent("", "  ")
	}
	e.Encode(data)
}
