/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package oauth2

import (
	"github.com/imdario/mergo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/heroku"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/paypal"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/soundcloud"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/stripe"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/wepay"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yammer"
	"github.com/webx-top/echo"
)

const (
	// DefaultPath /oauth
	DefaultPath = "/oauth"

	// DefaultContextKey oauth_user
	DefaultContextKey = "oauth_user"
)

type Account struct {
	On          bool // on / off
	Name        string
	Key         string
	Secret      string `json:"-" xml:"-"`
	Extra       echo.H
	LoginURL    string
	CallbackURL string
	Constructor func(*Account) goth.Provider `json:"-" xml:"-"`
}

func (a *Account) SetConstructor(constructor func(*Account) goth.Provider) {
	a.Constructor = constructor
}

func (a *Account) Instance() goth.Provider {
	return a.Constructor(a)
}

// Config the configs for the gothic oauth/oauth2 authentication for third-party websites
// All Key and Secret values are empty by default strings. Non-empty will be registered as Goth Provider automatically, by Iris
// the users can still register their own providers using goth.UseProviders
// contains the providers' keys  (& secrets) and the relative auth callback url path(ex: "/auth" will be registered as /auth/:provider/callback)
//
type Config struct {
	Host, Path string
	Accounts   []*Account

	// defaults to 'oauth_user' used by plugin to give you the goth.User, but you can take this manually also by `context.Get(ContextKey).(goth.User)`
	ContextKey string
}

// DefaultConfig returns OAuth config, the fields of the iteral are zero-values ( empty strings)
func DefaultConfig() *Config {
	return &Config{
		Path:       DefaultPath,
		Accounts:   []*Account{},
		ContextKey: DefaultContextKey,
	}
}

// MergeSingle merges the default with the given config and returns the result
func (c *Config) MergeSingle(cfg *Config) (config *Config) {
	config = cfg
	mergo.Merge(config, c)
	return
}

func (c *Config) CallbackURL(providerName string) string {
	return c.Host + c.Path + "/callback/" + providerName
}

func (c *Config) LoginURL(providerName string) string {
	return c.Host + c.Path + "/login/" + providerName
}

// GenerateProviders returns the valid goth providers and the relative url paths (because the goth.Provider doesn't have a public method to get the Auth path...)
// we do the hard-core/hand checking here at the configs.
//
// receives one parameter which is the host from the server,ex: http://localhost:3000, will be used as prefix for the oauth callback
func (c *Config) GenerateProviders() *Config {
	goth.ClearProviders()
	var providers []goth.Provider
	//we could use a map but that's easier for the users because of code completion of their IDEs/editors
	for _, account := range c.Accounts {
		if !account.On {
			continue
		}
		if provider := c.NewProvider(account); provider != nil {
			providers = append(providers, provider)
		}
	}
	goth.UseProviders(providers...)
	return c
}

func (c *Config) NewProvider(account *Account) goth.Provider {
	if len(account.LoginURL) == 0 {
		account.LoginURL = c.LoginURL(account.Name)
	}
	if len(account.CallbackURL) == 0 {
		account.CallbackURL = c.CallbackURL(account.Name)
	}
	if account.Constructor != nil {
		return account.Instance()
	}
	switch account.Name {
	case "twitter":
		return twitter.New(account.Key, account.Secret, account.CallbackURL)
	case "facebook":
		return facebook.New(account.Key, account.Secret, account.CallbackURL)
	case "gplus":
		return gplus.New(account.Key, account.Secret, account.CallbackURL)
	case "github":
		return github.New(account.Key, account.Secret, account.CallbackURL)
	case "spotify":
		return spotify.New(account.Key, account.Secret, account.CallbackURL)
	case "linkedin":
		return linkedin.New(account.Key, account.Secret, account.CallbackURL)
	case "lastfm":
		return lastfm.New(account.Key, account.Secret, account.CallbackURL)
	case "twitch":
		return twitch.New(account.Key, account.Secret, account.CallbackURL)
	case "dropbox":
		return dropbox.New(account.Key, account.Secret, account.CallbackURL)
	case "digitalocean":
		return digitalocean.New(account.Key, account.Secret, account.CallbackURL)
	case "bitbucket":
		return bitbucket.New(account.Key, account.Secret, account.CallbackURL)
	case "instagram":
		return instagram.New(account.Key, account.Secret, account.CallbackURL)
	case "box":
		return box.New(account.Key, account.Secret, account.CallbackURL)
	case "salesforce":
		return salesforce.New(account.Key, account.Secret, account.CallbackURL)
	case "amazon":
		return amazon.New(account.Key, account.Secret, account.CallbackURL)
	case "yammer":
		return yammer.New(account.Key, account.Secret, account.CallbackURL)
	case "onedrive":
		return onedrive.New(account.Key, account.Secret, account.CallbackURL)
	case "yahoo":
		return yahoo.New(account.Key, account.Secret, account.CallbackURL)
	case "slack":
		return slack.New(account.Key, account.Secret, account.CallbackURL)
	case "stripe":
		return stripe.New(account.Key, account.Secret, account.CallbackURL)
	case "wepay":
		return wepay.New(account.Key, account.Secret, account.CallbackURL)
	case "paypal":
		return paypal.New(account.Key, account.Secret, account.CallbackURL)
	case "steam":
		return steam.New(account.Key, account.CallbackURL)
	case "heroku":
		return heroku.New(account.Key, account.Secret, account.CallbackURL)
	case "uber":
		return uber.New(account.Key, account.Secret, account.CallbackURL)
	case "soundcloud":
		return soundcloud.New(account.Key, account.Secret, account.CallbackURL)
	case "gitlab":
		return gitlab.New(account.Key, account.Secret, account.CallbackURL)
	}
	return nil
}

func (c *Config) AddAccount(accounts ...*Account) *Config {
	c.Accounts = append(c.Accounts, accounts...)
	return c
}

func (c *Config) SetAccount(newAccount *Account) *Config {
	var exists bool
	for index, account := range c.Accounts {
		if account.Name != newAccount.Name {
			continue
		}
		isOff := account.On && !newAccount.On
		account.On = newAccount.On
		account.Key = newAccount.Key
		account.Secret = newAccount.Secret
		account.Extra = newAccount.Extra
		account.Constructor = newAccount.Constructor
		account.LoginURL = newAccount.LoginURL
		account.CallbackURL = newAccount.CallbackURL
		c.Accounts[index] = account
		if isOff {
			c.GenerateProviders()
		} else if account.On {
			if provider := c.NewProvider(account); provider != nil {
				goth.UseProviders(provider)
			}
		}
		exists = true
		break
	}
	if !exists {
		c.Accounts = append(c.Accounts, newAccount)
		if newAccount.On {
			if provider := c.NewProvider(newAccount); provider != nil {
				goth.UseProviders(provider)
			}
		}
	}
	return c
}
