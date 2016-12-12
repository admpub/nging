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
	"net/http"

	"github.com/markbates/goth"
	"github.com/webx-top/echo"
)

// OAuth is a plugin which helps you to use OAuth/OAuth2 apis from famous websites
type OAuth struct {
	Config          Config
	HostURL         string
	successHandlers []interface{}
	failHandler     echo.HTTPErrorHandler
}

// New returns a new OAuth plugin
// receives one parameter of type 'Config'
func New(hostURL string, cfg Config) *OAuth {
	c := DefaultConfig().MergeSingle(cfg)
	return &OAuth{Config: c, HostURL: hostURL}
}

// Success registers handler(s) which fires when the user logged in successfully
func (p *OAuth) Success(handlersFn ...interface{}) {
	p.successHandlers = append(p.successHandlers, handlersFn...)
}

// Fail registers handler which fires when the user failed to logged in
// underhood it justs registers an error handler to the StatusUnauthorized(400 status code), same as 'iris.OnError(400,handler)'
func (p *OAuth) Fail(handler echo.HTTPErrorHandler) {
	p.failHandler = handler
}

// User returns the user for the particular client
// if user is not validated  or not found it returns nil
// same as 'ctx.Get(config's ContextKey field).(goth.User)'
func (p *OAuth) User(ctx echo.Context) (u goth.User) {
	return ctx.Get(p.Config.ContextKey).(goth.User)
}

// Wrapper register the oauth route
func (p *OAuth) Wrapper(e *echo.Echo) {
	oauthProviders := p.Config.GenerateProviders(p.HostURL)
	if len(oauthProviders) > 0 {
		goth.UseProviders(oauthProviders...)
		// set the mux path to handle the registered providers
		e.Get(p.Config.Path+"/login/:provider", func(ctx echo.Context) error {
			err := BeginAuthHandler(ctx)
			if err != nil {
				e.Logger().Error("[OAUTH2] Error: " + err.Error())
			}
			return err
		})

		authMiddleware := func(h echo.Handler) echo.Handler {
			return echo.HandlerFunc(func(ctx echo.Context) error {
				user, err := CompleteUserAuth(ctx)
				if err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
				}
				ctx.Set(p.Config.ContextKey, user)
				return h.Handle(ctx)
			})
		}

		p.successHandlers = append([]interface{}{authMiddleware}, p.successHandlers...)
		lastIndex := len(p.successHandlers) - 1
		if lastIndex == 0 {
			e.Get(p.Config.Path+"/callback/:provider", func(ctx echo.Context) error {
				return ctx.String(`Success Handler is not set`)
			}, p.successHandlers...)
		} else {
			e.Get(p.Config.Path+"/callback/:provider", p.successHandlers[lastIndex], p.successHandlers[0:lastIndex]...)
		}
		// register the error handler
		if p.failHandler != nil {
			e.SetHTTPErrorHandler(p.failHandler)
		}
	}
}
