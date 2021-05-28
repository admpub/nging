package common

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/library/config/subconfig/scookie"
)

// APIKeyGetter API Key
type APIKeyGetter interface {
	APIKey() string
}

type CookieConfigGetter interface {
	CookieConfig() scookie.Config
}

func CookieConfig() scookie.Config {
	return echo.Get(`DefaultConfig`).(CookieConfigGetter).CookieConfig()
}

func Setting(group ...string) echo.H {
	return echo.GetStoreByKeys(`NgingConfig`, group...)
}
