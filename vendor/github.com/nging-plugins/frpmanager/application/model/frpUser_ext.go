package model

import "github.com/nging-plugins/frpmanager/application/dbschema"

type FrpUserAndServer struct {
	*dbschema.NgingFrpUser
	Server *dbschema.NgingFrpServer `db:"-,relation=id:server_id|gtZero"`
}
