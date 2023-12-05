package model

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/backend/oauth2server/oauth2serverutils"
	"github.com/admpub/nging/v5/application/library/common"
)

func NewOAuthApp(ctx echo.Context) *OAuthApp {
	m := &OAuthApp{
		NgingOauthApp: dbschema.NewNgingOauthApp(ctx),
	}
	return m
}

type OAuthApp struct {
	*dbschema.NgingOauthApp
}

func (f *OAuthApp) check() error {
	f.SiteDomains = strings.TrimSpace(f.SiteDomains)
	if len(f.SiteDomains) > 0 {
		f.SiteDomains = strings.Join(param.StringSlice(com.TrimSpaceForRows(f.SiteDomains)).Unique().String(), `,`)
	}
	return nil
}

func (f *OAuthApp) Add() (pk interface{}, err error) {
	err = f.Context().Begin()
	if err != nil {
		return
	}
	defer func() {
		f.Context().End(err == nil)
	}()
	err = f.check()
	if err != nil {
		return
	}
	f.NgingOauthApp.AppId, err = f.GenerateAppID()
	if err != nil {
		return
	}
	f.NgingOauthApp.AppSecret = f.GenerateAppSecret()
	pk, err = f.NgingOauthApp.Insert()
	return
}

func (f *OAuthApp) GenerateAppID() (string, error) {
	return common.UniqueID()
}

func (f *OAuthApp) GenerateAppSecret() string {
	return com.RandomAlphanumeric(32)
}

func (f *OAuthApp) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	err = f.Context().Begin()
	if err != nil {
		return
	}
	defer func() {
		f.Context().End(err == nil)
	}()
	err = f.check()
	if err != nil {
		return
	}
	if len(f.NgingOauthApp.AppId) == 0 {
		f.NgingOauthApp.AppId, err = f.GenerateAppID()
		if err != nil {
			return
		}
	}
	if len(f.NgingOauthApp.AppSecret) == 0 {
		f.NgingOauthApp.AppSecret = f.GenerateAppSecret()
	}
	err = f.NgingOauthApp.Update(mw, args...)
	return
}

func (f *OAuthApp) Exists(appID string) error {
	exists, err := f.NgingOauthApp.Exists(nil, db.Cond{`app_id`: appID})
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().E(`AppID“%s”已经存在`, appID)
	}
	return err
}

func (f *OAuthApp) MatchDomain(domain string) bool {
	return MatchDomain(domain, f.NgingOauthApp)
}

func (f *OAuthApp) MatchScope(scopes ...string) bool {
	return len(InvalidScope(scopes, f.NgingOauthApp)) == 0
}

func (f *OAuthApp) GetAndVerify(appID string, domain string, scopes []string) error {
	err := f.GetByAppID(appID)
	if err != nil {
		return err
	}
	return f.VerifyApp(domain, scopes)
}

func (f *OAuthApp) VerifyApp(domain string, scopes []string) error {
	return VerifyApp(f.NgingOauthApp, domain, scopes)
}

func (f *OAuthApp) GetByAppID(appID string, ignoreStatus ...bool) error {
	err := f.NgingOauthApp.Get(nil, db.Cond{`app_id`: appID})
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = f.Context().NewError(stdCode.DataNotFound, `appID无效`)
		}
		return err
	}
	if len(ignoreStatus) == 0 || !ignoreStatus[0] {
		if f.NgingOauthApp.Disabled != common.BoolN {
			return f.Context().NewError(stdCode.DataUnavailable, `应用已停用`)
		}
	}
	return err
}

func MatchDomain(domain string, f *dbschema.NgingOauthApp) bool {
	if f.Id < 1 {
		return false
	}
	return oauth2serverutils.MatchDomain(domain, f.SiteDomains)
}

func InvalidScope(scopes []string, f *dbschema.NgingOauthApp) (invalidScopes []string) {
	if f.Id < 1 {
		return []string{`invalid-data`}
	}
	if len(f.Scopes) == 0 || len(scopes) == 0 {
		return
	}
	allowScopes := strings.Split(f.Scopes, `,`)
	for _, scope := range scopes {
		if !com.InSlice(scope, allowScopes) {
			invalidScopes = append(invalidScopes, scope)
		}
	}
	return
}

func VerifyApp(f *dbschema.NgingOauthApp, domain string, scopes []string) error {
	ctx := f.Context()
	if !MatchDomain(domain, f) {
		return ctx.NewError(stdCode.InvalidType, ctx.T(`应用与域名不匹配`))
	}
	invalidScopes := InvalidScope(scopes, f)
	if len(invalidScopes) > 0 {
		return ctx.NewError(stdCode.InvalidType, ctx.T(`应用未获取授权: %s`, strings.Join(invalidScopes, `, `)))
	}
	return nil
}
