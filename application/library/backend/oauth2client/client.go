package oauth2client

import (
	"github.com/admpub/goth"
	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/backend/oauth2nging"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/nging/v5/application/registry/route"
	"github.com/coscms/oauth2s/client/goth/providers"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/handler/oauth2"
	"github.com/webx-top/echo/subdomains"
)

var (
	defaultOAuth      *oauth2.OAuth
	SuccessHandler    interface{}  = successHandler
	BeginAuthHandler  echo.Handler = echo.HandlerFunc(oauth2.BeginAuthHandler)
	AfterLoginSuccess []func(ctx echo.Context, ouser *goth.User) (end bool, err error)
)

func Default() *oauth2.OAuth {
	return defaultOAuth
}

func OnChangeBackendURL(d config.Diff) error {
	if defaultOAuth == nil || !d.IsDiff {
		return nil
	}
	host := d.String()
	if len(host) == 0 {
		host = subdomains.Default.URL(``, `backend`)
	}
	defaultOAuth.HostURL = host
	defaultOAuth.Config.RangeAccounts(func(account *oauth2.Account) bool {
		// 清空生成的网址，以便于在后面的 GenerateProviders() 函数中重新生成新的网址
		account.CallbackURL = ``
		account.LoginURL = ``
		return true
	})
	defaultOAuth.Config.GenerateProviders()
	return nil
}

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
	On       bool           `json:"on" xml:"on"` // on / off
	Accounts []OAuthAccount `json:"accounts" xml:"accounts"`
}

func (c *OAuth2Config) ToAccounts() []*oauth2.Account {
	if !c.On {
		return []*oauth2.Account{}
	}
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
	if c == nil || defaultOAuth == nil {
		return nil
	}
	accounts := c.ToAccounts()
	defaultOAuth.Config.ClearAccounts()
	if len(accounts) > 0 {
		defaultOAuth.Config.AddAccount(accounts...)
		defaultOAuth.Config.GenerateProviders()
	}
	log.Debugf(`reload backend oauth configuration information`)
	return nil
}

type OAuthAccount struct {
	On     bool   `json:"on" xml:"on"` // on / off
	Name   string `json:"name" xml:"name"`
	AppID  string `json:"appID" xml:"appID"`
	Secret string `json:"secret" xml:"secret"`
	Extra  echo.H `json:"extra" xml:"extra"`
}

func (c *OAuthAccount) ToAccount() *oauth2.Account {
	a := &oauth2.Account{
		On:     c.On,
		Name:   c.Name,
		Secret: c.Secret,
		Key:    c.AppID,
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
	host := common.BackendURL(nil)
	if len(host) == 0 {
		host = subdomains.Default.URL(``, `backend`)
	}
	oauth2Config := oauth2.NewConfig()
	RegisterProvider(oauth2Config)
	if oCfg != nil {
		oauth2Config.AddAccount(oCfg.ToAccounts()...)
	}
	defaultOAuth = oauth2.New(host, oauth2Config)
	defaultOAuth.SetSuccessHandler(SuccessHandler)
	defaultOAuth.SetBeginAuthHandler(BeginAuthHandler)
	e.Group(defaultOAuth.Config.Path).SetMetaKV(route.PermGuestKV())
	defaultOAuth.Wrapper(e)
}

// RegisterProvider 注册Provider
func RegisterProvider(c *oauth2.Config) {
	providers.Register(`nging`, func(account *oauth2.Account) goth.Provider {
		hostURL := account.Extra.String(`hostURL`)
		if len(account.CallbackURL) == 0 {
			account.CallbackURL = oauth2.DefaultPath + "/callback/" + account.Name
		}
		m := oauth2nging.New(account.Key, account.Secret, account.CallbackURL, hostURL, `profile`)
		m.SetName(`nging`)
		return m
	})
}

func DefaultOAuth() *oauth2.OAuth {
	doa := defaultOAuth
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
				if err != nil {
					handler.SendErr(ctx, ctx.E(`绑定失败: %s`, err.Error()))
				} else {
					handler.SendOk(ctx, ctx.T(`绑定成功`))
				}
				return ctx.Redirect(handler.URLFor(`/user/oauth`))
			}
			err = ctx.NewError(code.DataNotFound, `请先绑定账号`)
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
	if defaultOAuth == nil {
		return nil
	}
	oCfg, _ := common.ExtendConfig().Get(`oauth2backend`).(*OAuth2Config)
	if oCfg == nil {
		return nil
	}
	accounts := oCfg.ToAccounts()
	defaultOAuth.Config.ClearAccounts()
	defaultOAuth.Config.AddAccount(accounts...)
	defaultOAuth.Config.GenerateProviders()
	log.Debug(`update backend oauth configuration information`)
	return nil
}

func GetOAuthAccounts(skipHided ...bool) []oauth2.Account {
	var accounts []oauth2.Account
	if defaultOAuth == nil {
		return accounts
	}
	var skipHide bool
	if len(skipHided) > 0 {
		skipHide = skipHided[0]
	}
	defaultOAuth.Config.RangeAccounts(func(a *oauth2.Account) bool {
		if !a.On || (skipHide && a.Extra != nil && a.Extra.Bool(`hide`)) {
			return true
		}
		account := *a
		accounts = append(accounts, account)
		return true
	})
	return accounts
}

func OnInstalled(ctx echo.Context) error {
	return UpdateOAuthAccount()
}
