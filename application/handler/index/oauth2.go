package index

import (
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/backend/oauth2client"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/extend"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/registry/route"
)

func init() {
	route.Register(func(e echo.RouteRegister) {
		oauth2client.InitOauth(route.IRegister().Echo(), httpserver.SearchEngineNoindex())
	})

	backend.OnInstalled(oauth2client.OnInstalled)
	extend.Register(`oauth2backend`, func() interface{} {
		return &oauth2client.OAuth2Config{}
	})
	config.OnKeySetSettings(`base.backendURL`, oauth2client.OnChangeBackendURL)
}
