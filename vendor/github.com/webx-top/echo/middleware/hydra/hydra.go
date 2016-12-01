package hydra

import (
	"fmt"

	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/firewall"
	hydraSDK "github.com/ory-am/hydra/sdk"
	"github.com/ory-am/ladon"
	"github.com/webx-top/echo"
)

var DefaultClient *hydraSDK.Client

type Options struct {
	Skipper      echo.Skipper
	ClientID     string
	ClientSecret string
	ClusterURL   string
}

func Connect(val Options) (hc *hydraSDK.Client, err error) {
	hc, err = hydraSDK.Connect(
		hydraSDK.ClientID(val.ClientID),
		hydraSDK.ClientSecret(val.ClientSecret),
		hydraSDK.ClusterURL(val.ClusterURL),
	)
	return
}

func GetClient(hc *hydraSDK.Client, id string) (fosite.Client, error) {
	return hc.Clients.GetClient(id)
}

func GetContext(c echo.Context) *firewall.Context {
	ctx, _ := c.Get("hydra").(*firewall.Context)
	return ctx
}

func NewTokenAccessRequest(resource string, action string, context map[string]interface{}) *firewall.TokenAccessRequest {
	return &firewall.TokenAccessRequest{
		Resource: resource,
		Action:   action,
		Context:  ladon.Context(context),
	}
}

func ScopesRequired(opt interface{}, tokenAccessRequest *firewall.TokenAccessRequest, scopes ...string) echo.MiddlewareFunc {
	var hc *hydraSDK.Client
	var err error
	var skipper echo.Skipper
	if client, ok := opt.(*hydraSDK.Client); ok {
		hc = client
	} else if val, ok := opt.(*Options); ok {
		skipper = val.Skipper
		hc, err = Connect(*val)
	} else if val, ok := opt.(Options); ok {
		skipper = val.Skipper
		hc, err = Connect(val)
	} else if DefaultClient != nil {
		hc = DefaultClient
	} else {
		err = fmt.Errorf("invalid parameter: %T", opt)
	}

	if err != nil {
		panic(err.Error())
	}
	if DefaultClient == nil {
		DefaultClient = hc
	}
	if skipper == nil {
		skipper = echo.DefaultSkipper
	}
	if tokenAccessRequest == nil {
		tokenAccessRequest = NewTokenAccessRequest("matrix", "create", map[string]interface{}{})
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if skipper(c) {
				return h.Handle(c)
			}
			ctx, err := hc.Warden.TokenAllowed(
				c,
				hc.Warden.TokenFromRequest(c.Request().StdRequest()),
				tokenAccessRequest,
				scopes...,
			)
			if err != nil {
				return err
			}
			// All required scopes are found
			c.Set("hydra", ctx)
			return h.Handle(c)
		})
	}
}
