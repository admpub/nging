package webauthn

import (
	"github.com/admpub/nging/v5/application/library/module"
	"github.com/admpub/nging/v5/application/library/route"

	"github.com/nging-plugins/webauthn/application/handler/backend"
	//"github.com/nging-plugins/webauthn/application/library/customer"
	"github.com/nging-plugins/webauthn/application/library/user"
)

const ID = `webauthn`

var Module = module.Module{
	TemplatePath: map[string]string{
		ID: `webauthn/template/backend`,
	},
	AssetsPath: []string{},
	//Navigate: ,
	Route: func(r *route.Collection) {
		user.RegisterBackend(r.Backend.Echo().Group(`/user`))
		user.RegisterLogin(r.Backend.Echo())
		backend.Register(r.Backend.Echo())
		//customer.RegisterFrontend(r.Frontend.Echo())
		//customer.RegisterLogin(r.Frontend.Echo())
	},
	DBSchemaVer: 0.0000,
}
