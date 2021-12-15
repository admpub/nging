package model

import "github.com/nging-plugins/frpmanager/pkg/dbschema"

type FrpUserAndServer struct {
	*dbschema.NgingFrpUser
	Server *dbschema.NgingFrpServer `db:"-,relation=id:server_id|gtZero"`
}
