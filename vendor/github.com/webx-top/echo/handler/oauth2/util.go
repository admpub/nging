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
	"errors"
	"net/url"

	"github.com/markbates/goth"
	"github.com/webx-top/echo"
)

// SessionName is the key used to access the session store.
// we could use the echo's sessions default, but this session should be not confict with the cookie session name defined by the sessions manager
const SessionName = "EchoGothSession"

var (
	_         goth.Params = url.Values{}
	EmptyUser             = goth.User{}
)

/*
BeginAuthHandler is a convienence handler for starting the authentication process.
It expects to be able to get the name of the provider from the named parameters
as either "provider" or url query parameter ":provider".
BeginAuthHandler will redirect the user to the appropriate authentication end-point
for the requested provider.
*/
func BeginAuthHandler(ctx echo.Context) error {
	url, err := GetAuthURL(ctx)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	next := ctx.Form(`next`)
	if len(next) > 0 {
		ctx.Session().Set(`next`, next)
	}
	return ctx.Redirect(url)
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var SetState = func(ctx echo.Context) string {
	state := ctx.Query("state")
	if len(state) > 0 {
		return state
	}

	return "state"

}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(ctx echo.Context) string {
	return ctx.Query("state")
}

/*
GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or url query parameter ":provider".
I would recommend using the BeginAuthHandler instead of doing all of these steps
yourself, but that's entirely up to you.
*/
func GetAuthURL(ctx echo.Context) (string, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return "", err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth(SetState(ctx))
	if err != nil {
		return "", err
	}

	if cr, ok := sess.(echo.ContextRegister); ok {
		cr.SetContext(ctx)
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}
	length := len(url)
	if length > 0 {
		switch url[0] {
		case '/':
			url = ctx.Site() + url
		case '.':
			url = ctx.Site() + `/` + url
		default:
			if length > 7 {
				switch url[0:7] {
				case `https:/`, `http://`:
				default:
					url = ctx.Site() + `/` + url
				}
			}
		}
	}
	//panic(sess.Marshal())
	err = ctx.Session().Set(SessionName, sess.Marshal()).Save()
	return url, err
}

/*
CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all of the basic information about the user from the provider.
It expects to be able to get the name of the provider from the named parameters
as either "provider" or url query parameter ":provider".
*/
var CompleteUserAuth = func(ctx echo.Context) (goth.User, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return EmptyUser, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return EmptyUser, err
	}

	//error=invalid_request&error_description=The provided value for the input parameter 'redirect_uri' is not valid. The scope 'openid offline_access user.read' requires that the request must be sent over a secure connection using SSL.&state=state
	errorDescription := ctx.Query(`error_description`)
	if len(errorDescription) > 0 {
		return EmptyUser, errors.New(providerName + `: ` + errorDescription)
	}

	sv, ok := ctx.Session().Get(SessionName).(string)
	if !ok || len(sv) == 0 {
		return EmptyUser, errors.New("could not find a matching session for this request")
	}

	sess, err := provider.UnmarshalSession(sv)
	if err != nil {
		return EmptyUser, err
	}

	if cr, ok := sess.(echo.ContextRegister); ok {
		cr.SetContext(ctx)
	}

	_, err = sess.Authorize(provider, url.Values(ctx.Queries()))

	if err != nil {
		return EmptyUser, err
	}

	return provider.FetchUser(sess)
}

// GetProviderName is a function used to get the name of a provider
// for a given request. By default, this provider is fetched from
// the URL query string. If you provide it in a different way,
// assign your own function to this variable that returns the provider
// name for your request.
var GetProviderName = getProviderName

func getProviderName(ctx echo.Context) (string, error) {
	provider := ctx.Param("provider")
	if len(provider) == 0 {
		provider = ctx.Query("provider")
	} else {
		return provider, nil
	}
	if len(provider) == 0 {
		return provider, errors.New("you must select a provider")
	}
	return provider, nil
}
