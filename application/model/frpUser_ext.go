package model

import "github.com/admpub/nging/application/dbschema"

type FrpUserAndServer struct {
	*dbschema.NgingFrpUser
	Server *dbschema.NgingFrpServer `db:"-,relation=id:server_id|gtZero"`
}
