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
	"fmt"

	"github.com/markbates/goth"
	"github.com/webx-top/echo"
)

// SessionName is the key used to access the session store.
// we could use the echo's sessions default, but this session should be not confict with the cookie session name defined by the sessions manager
const SessionName = "EchoGothSession"

// GothParams used to convert the context.URLParams to goth's params
type GothParams map[string][]string

// Get returns the value of
func (g GothParams) Get(key string) string {
	if v, y := g[key]; y {
		if len(v) > 0 {
			return v[0]
		}
	}
	return ``
}

var _ goth.Params = GothParams{}

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
	if ctx.Session() == nil {
		fmt.Println("You have to enable sessions")
	}

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

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}
	//fmt.Println(sess.Marshal())
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

	if ctx.Session() == nil {
		fmt.Println("You have to enable sessions")
	}

	providerName, err := GetProviderName(ctx)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}

	sv := ctx.Session().Get(SessionName)
	if sv == nil {
		return goth.User{}, errors.New("could not find a matching session for this request")
	}

	sess, err := provider.UnmarshalSession(sv.(string))
	if err != nil {
		return goth.User{}, err
	}
	_, err = sess.Authorize(provider, GothParams(ctx.Queries()))

	if err != nil {
		return goth.User{}, err
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
	if provider == "" {
		provider = ctx.Query(":provider")
	}
	if provider == "" {
		return provider, errors.New("you must select a provider")
	}
	return provider, nil
}
