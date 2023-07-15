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
		r.Backend.Register(user.RegisterLogin)
		r.Backend.Register(backend.Register)
		r.Backend.RegisterToGroup(`/user`, user.RegisterBackend)
		//r.Frontend.RegisterToGroup(`/user`, customer.RegisterFrontend)
		//r.Frontend.Register(customer.RegisterLogin)
	},
	DBSchemaVer: 0.0000,
}
