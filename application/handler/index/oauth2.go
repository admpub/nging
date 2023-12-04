package index

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/handler/setup"
	"github.com/admpub/nging/v5/application/library/backend/oauth2client"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/extend"
)

func init() {
	handler.Register(func(e echo.RouteRegister) {
		oauth2client.InitOauth(handler.IRegister().Echo())
	})

	setup.OnInstalled(oauth2client.OnInstalled)
	extend.Register(`oauth2backend`, func() interface{} {
		return &oauth2client.OAuth2Config{}
	})
	config.OnKeySetSettings(`base.backendURL`, oauth2client.OnChangeBackendURL)
}
