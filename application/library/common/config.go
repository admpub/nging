package common

import (
	"strings"

	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/config/subconfig/scookie"
)

const ConfigName = `DefaultConfig`

// APIKeyGetter API Key
type APIKeyGetter interface {
	APIKey() string
}

type CookieConfigGetter interface {
	CookieConfig() scookie.Config
}

func CookieConfig() scookie.Config {
	return echo.Get(ConfigName).(CookieConfigGetter).CookieConfig()
}

func Setting(group ...string) echo.H {
	return echo.GetStoreByKeys(`NgingConfig`, group...)
}

func BackendURL(ctx echo.Context) string {
	backendURL := Setting(`base`).String(`backendURL`)
	if len(backendURL) == 0 {
		if ctx == nil {
			return backendURL
		}
		backendURL = ctx.Site()
	}
	backendURL = strings.TrimSuffix(backendURL, `/`)
	return backendURL
}

func SystemAPIKey() string {
	apiKey := Setting(`base`).String(`apiKey`)
	return apiKey
}
