// Package oauth2nging implements the OAuth2 protocol for authenticating users through nging.
// This package can be used as a reference implementation of an OAuth2 provider for Goth.
package oauth2nging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/goth"
	oauth2c "github.com/coscms/oauth2s/client/goth/oauth2"
	"golang.org/x/oauth2"
)

var (
	AuthURL    = "/oauth2/authorize"
	TokenURL   = "/oauth2/token"
	ProfileURL = "/oauth2/profile"
)

// New creates a new Webx provider, and sets up important connection details.
// You should always call `webx.New` to get a new Provider. Never try to create
// one manually.
func New(clientKey, secret, callbackURL, hostURL string, scopes ...string) *Provider {
	return NewCustomisedURL(clientKey, secret, callbackURL, hostURL+AuthURL, hostURL+TokenURL, hostURL+ProfileURL, scopes...)
}

// NewCustomisedURL is similar to New(...) but can be used to set custom URLs to connect to
func NewCustomisedURL(clientKey, secret, callbackURL, authURL, tokenURL, profileURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:         clientKey,
		Secret:            secret,
		CallbackURL:       callbackURL,
		isFullCallbackURL: strings.Contains(callbackURL, `://`),
		HTTPClient:        oauth2c.DefaultClient,
		providerName:      "nging",
		profileURL:        profileURL,
	}
	p.config = newConfig(p, authURL, tokenURL, scopes)
	return p
}

// Provider is the implementation of `goth.Provider` for accessing Github.
type Provider struct {
	ClientKey         string
	Secret            string
	CallbackURL       string
	isFullCallbackURL bool
	HTTPClient        *http.Client
	config            *oauth2.Config
	providerName      string
	profileURL        string
}

// Name is the name used to retrieve this provider later.
func (p *Provider) Name() string {
	return p.providerName
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

// Debug is a no-op for the github package.
func (p *Provider) Debug(debug bool) {}

// BeginAuth asks Github for an authentication end-point.
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	url := p.config.AuthCodeURL(state)
	session := &Session{
		AuthURL: url,
	}
	return session, nil
}

// FetchUser will go to Webx and access basic information about the user.
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*Session)
	user := goth.User{
		AccessToken:  sess.AccessToken,
		RefreshToken: sess.RefreshToken,
		ExpiresAt:    sess.Expiry,
		Provider:     p.Name(),
	}

	if len(user.AccessToken) == 0 {
		// data is not yet retrieved since accessToken is still empty
		return user, fmt.Errorf("%s cannot get user information without accessToken", p.providerName)
	}

	req, err := http.NewRequest("GET", p.profileURL, nil)
	if err != nil {
		return user, err
	}

	req.Header.Add("Authorization", "Bearer "+sess.AccessToken)
	response, err := p.Client().Do(req)
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return user, fmt.Errorf("Nging API responded with a %d trying to fetch user information", response.StatusCode)
	}

	bits, err := io.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	err = json.NewDecoder(bytes.NewReader(bits)).Decode(&user.RawData)
	if err != nil {
		return user, err
	}

	err = userFromReader(bytes.NewReader(bits), &user)
	if err != nil {
		return user, err
	}
	return user, err
}

func userFromReader(reader io.Reader, user *goth.User) error {
	u := struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		AvatarURL   string `json:"avatar"`
		Email       string `json:"email"`
		Description string `json:"description"`
		ExpiresIn   uint32 `json:"expires_in"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Name = u.Name
	user.NickName = ``
	user.Email = u.Email
	user.Description = u.Description
	user.AvatarURL = u.AvatarURL
	user.UserID = strconv.Itoa(u.ID)
	if u.ExpiresIn > 0 {
		user.ExpiresAt = time.Now().Add(time.Duration(u.ExpiresIn) * time.Second)
	}
	return err
}

func newConfig(provider *Provider, authURL, tokenURL string, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	c.Scopes = append(c.Scopes, scopes...)
	if len(c.Scopes) == 0 {
		c.Scopes = append(c.Scopes, `profile`)
	}

	return c
}

// RefreshToken refresh token
func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{RefreshToken: refreshToken}
	ts := p.config.TokenSource(goth.ContextForClient(p.Client()), token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return newToken, err
}

// RefreshTokenAvailable refresh token available
func (p *Provider) RefreshTokenAvailable() bool {
	return true
}
