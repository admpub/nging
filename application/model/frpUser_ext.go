package model

import "github.com/admpub/nging/v3/application/dbschema"

type FrpUserAndServer struct {
	*dbschema.NgingFrpUser
	Server *dbschema.NgingFrpServer `db:"-,relation=id:server_id|gtZero"`
}
