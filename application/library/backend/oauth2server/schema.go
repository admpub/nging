package oauth2server

import (
	"strings"

	"github.com/webx-top/echo"
)

type AuthCodeRequestData struct {
	ClientID     string   `json:"client_id" xml:"client_id"`
	RedirectURI  string   `json:"redirect_uri" xml:"redirect_uri"`
	ResponseType string   `json:"response_type" xml:"response_type"` // code|token
	Scope        []string `json:"scope" xml:"scope"`
	State        string   `json:"state" xml:"state"`
}

func (a *AuthCodeRequestData) FromContext(ctx echo.Context) *AuthCodeRequestData {
	a.ClientID = ctx.Form("client_id")
	a.RedirectURI = ctx.Form("redirect_uri")
	a.ResponseType = ctx.Form("response_type")
	scope := ctx.Formx(`scope`).String()
	if len(scope) > 0 {
		a.Scope = strings.Split(scope, " ")
	}
	a.State = ctx.Form("state")
	return a
}

type TokenRequestData struct {
	Code        string `json:"code" xml:"code"`
	GrantType   string `json:"grant_type" xml:"grant_type"`
	RedirectURI string `json:"redirect_uri" xml:"redirect_uri"`
}

func (a *TokenRequestData) FromContext(ctx echo.Context) *TokenRequestData {
	a.Code = ctx.Form("code")
	a.GrantType = ctx.Form("grant_type")
	a.RedirectURI = ctx.Form("redirect_uri")
	return a
}
