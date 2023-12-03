package index

import (
	"sync"

	"github.com/markbates/goth"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/handler/oauth2"
	"github.com/webx-top/echo/middleware/session"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/handler/setup"
	"github.com/admpub/nging/v5/application/initialize/backend"
	"github.com/admpub/nging/v5/application/library/backend/oauth2nging"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/extend"
	"github.com/admpub/nging/v5/application/model"
	"github.com/coscms/oauth2s/client/goth/providers"
)

func init() {
	handler.Register(func(e echo.RouteRegister) {
		InitOauth(handler.IRegister().Echo())
	})

	setup.OnInstalled(onInstalled)
	extend.Register(`oauth2backend`, func() interface{} {
		return &OAuth2Config{}
	})
}

var (
	oauthLock         sync.RWMutex
	defaultOAuth      *oauth2.OAuth
	SuccessHandler    interface{}  = successHandler
	BeginAuthHandler  echo.Handler = echo.HandlerFunc(oauth2.BeginAuthHandler)
	AfterLoginSuccess []func(ctx echo.Context, ouser *goth.User) (end bool, err error)
)

func OnAfterLoginSuccess(hooks ...func(ctx echo.Context, ouser *goth.User) (end bool, err error)) {
	AfterLoginSuccess = append(AfterLoginSuccess, hooks...)
}

func FireAfterLoginSuccess(ctx echo.Context, ouser *goth.User) (end bool, err error) {
	for _, hook := range AfterLoginSuccess {
		end, err = hook(ctx, ouser)
		if err != nil || end {
			return
		}
	}
	return
}

type OAuth2Config struct {
	On       bool // on / off
	Accounts []OAuthAccount
}

func (c *OAuth2Config) ToAccounts() []*oauth2.Account {
	accounts := make([]*oauth2.Account, 0, len(c.Accounts))
	isProduction := config.FromFile().Sys.IsEnv(`prod`)
	for _, account := range c.Accounts {
		if !account.On {
			continue
		}
		acc := account.ToAccount()
		var provider func(account *oauth2.Account) goth.Provider
		if !isProduction {
			provider = providers.Get(acc.Name + `_dev`)
		}
		if provider == nil {
			provider = providers.Get(acc.Name)
		}
		if provider != nil {
			log.Infof(`backend oauth2 account: %s`, acc.Name)
			acc.SetConstructor(provider)
		} else {
			log.Errorf(`no provider for %q exists`, acc.Name)
		}
		accounts = append(accounts, acc)
	}
	return accounts
}

func (c *OAuth2Config) Reload() error {
	if c == nil {
		return nil
	}
	accounts := c.ToAccounts()
	oauthLock.Lock()
	defaultOAuth.Config.Accounts = accounts
	defaultOAuth.Config.GenerateProviders()
	oauthLock.Unlock()
	log.Debug(`Update backend oauth configuration information`)
	return nil
}

type OAuthAccount struct {
	On     bool // on / off
	Name   string
	Key    string
	Secret string
	Extra  echo.H
}

func (c *OAuthAccount) ToAccount() *oauth2.Account {
	a := &oauth2.Account{
		On:     c.On,
		Name:   c.Name,
		Secret: c.Secret,
		Key:    c.Key,
		Extra:  c.Extra,
	}
	if c.Extra == nil {
		a.Extra = map[string]interface{}{}
	}
	return a
}

// InitOauth 第三方登录
func InitOauth(e *echo.Echo) {
	oCfg, _ := common.ExtendConfig().Get(`oauth2backend`).(*OAuth2Config)
	if oCfg == nil || !oCfg.On || len(oCfg.Accounts) == 0 {
		return
	}
	host := subdomains.Default.URL(``, `backend`)
	if len(host) == 0 {
		backendDomain := config.FromCLI().BackendDomain
		if len(backendDomain) > 0 {
			host = backend.MakeSubdomains(backendDomain, backend.DefaultLocalHostNames)[0]
		}
		if len(host) == 0 {
			host = `127.0.0.1:28081`
		}
		host = `http://` + host
	}
	log.Warnf(`oauth host: %s`, host)
	oauth2Config := &oauth2.Config{}
	RegisterProvider(oauth2Config)
	oauth2Config.AddAccount(oCfg.ToAccounts()...)
	oauthLock.Lock()
	defaultOAuth = oauth2.New(host, oauth2Config)
	defaultOAuth.SetSuccessHandler(SuccessHandler)
	defaultOAuth.SetBeginAuthHandler(BeginAuthHandler)
	defaultOAuth.Wrapper(e, session.Middleware(config.SessionOptions))
	oauthLock.Unlock()
}

// RegisterProvider 注册Provider
func RegisterProvider(c *oauth2.Config) {
	providers.Register(`nging`, func(account *oauth2.Account) goth.Provider {
		hostURL := account.Extra.String(`HostURL`)
		if len(account.CallbackURL) == 0 {
			account.CallbackURL = hostURL + "/oauth/callback/" + account.Name //c.CallbackURL(account.Name)
		}
		m := oauth2nging.New(account.Key, account.Secret, account.CallbackURL, hostURL, `profile`)
		m.SetName(`nging`)
		return m
	})
}

func DefaultOAuth() *oauth2.OAuth {
	oauthLock.RLock()
	doa := defaultOAuth
	oauthLock.RUnlock()
	return doa
}

// 通过oauth登录第三方网站成功之后的处理
func successHandler(ctx echo.Context) error {
	doa := DefaultOAuth()
	var ouser *goth.User
	if user := doa.User(ctx); len(user.Provider) > 0 {
		ouser = &user
	}
	if len(ouser.UserID) == 0 {
		return ctx.NewError(code.InvalidParameter, `oauth2登录后获取UserID无效`)
	}
	end, err := FireAfterLoginSuccess(ctx, ouser)
	if err != nil || end {
		return err
	}
	var next string
	oauthM := model.NewUserOAuth(ctx)
	user := handler.User(ctx)
	err = oauthM.GetByOutUser(ouser)
	if err != nil {
		if err == db.ErrNoMoreRows { // 没有绑定过
			if user != nil { // 用户已经登录时，自动绑定
				oauthM.CopyFrom(ouser)
				oauthM.Uid = user.Id
				_, err = oauthM.Add()
			} else {
				err = ctx.NewError(code.DataNotFound, `请先绑定账号`)
			}
		}
		if err != nil {
			return err
		}
	} else {
		if user != nil && oauthM.Uid != user.Id {
			err = ctx.NewError(code.DataUnavailable, `此外部账号已经被其他用户绑定`)
			return err
		}
	}

	oauthSet := echo.H{}
	if ouser.AccessToken != oauthM.AccessToken {
		oauthSet[`access_token`] = ouser.AccessToken
	}
	if ouser.RefreshToken != oauthM.RefreshToken {
		oauthSet[`refresh_token`] = ouser.RefreshToken
	}
	if !ouser.ExpiresAt.IsZero() {
		oauthSet[`expired`] = ouser.ExpiresAt.Unix()
	}

	// 直接登录
	userM := model.NewUser(ctx)
	err = userM.Get(nil, `id`, oauthM.Uid)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		oauthM.Delete(nil, `id`, oauthM.Id) // 删除垃圾数据
		return ctx.NewError(code.UserNotFound, `用户不存在`)
	}
	// 更新用户的旧资料
	if len(ouser.AvatarURL) > 0 && ouser.AvatarURL != oauthM.Avatar {
		oauthSet[`avatar`] = ouser.AvatarURL
		if oauthM.Avatar == userM.Avatar {
			err = userM.NgingUser.UpdateField(nil, `avatar`, ouser.AvatarURL, `id`, userM.Id)
			if err != nil {
				log.Error(ctx.T(`更新本地用户头像为oauth2用户头像时失败`), `: `, err.Error())
			}
		}
	}

	if len(oauthSet) > 0 {
		err = oauthM.UpdateFields(nil, oauthSet, `id`, oauthM.Id)
		if err != nil {
			log.Error(ctx.T(`更新用户oauth2的数据(%s)失败`, echo.Dump(oauthSet, false)), `: `, err.Error())
		}
	}
	authType := model.AuthTypeOauth2
	// 未登录时设置登录状态
	if user == nil {
		err = userM.FireLoginSuccess(authType)
	}
	if err != nil {
		userM.FireLoginFailure(authType, ``, err)
		return err
	}
	if len(next) == 0 {
		next, _ = ctx.Session().Get(`next`).(string)
		if len(next) == 0 {
			next = `/`
		}
	}
	return ctx.Redirect(next)
}

// UpdateOAuthAccount 第三方登录平台账号
func UpdateOAuthAccount() error {
	oCfg, _ := common.ExtendConfig().Get(`oauth2backend`).(*OAuth2Config)
	if oCfg == nil || !oCfg.On || len(oCfg.Accounts) == 0 {
		return nil
	}
	accounts := oCfg.ToAccounts()
	oauthLock.Lock()
	defaultOAuth.Config.Accounts = accounts
	defaultOAuth.Config.GenerateProviders()
	oauthLock.Unlock()
	log.Debug(`Update backend oauth configuration information`)
	return nil
}

func getOAuthAccounts() []*oauth2.Account {
	oauthLock.RLock()
	account := defaultOAuth.Config.Accounts
	oauthLock.RUnlock()
	return account
}

func onInstalled(ctx echo.Context) error {
	return UpdateOAuthAccount()
}
