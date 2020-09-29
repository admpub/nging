package model

import "github.com/admpub/nging/application/dbschema"

type CloudBackupListItem struct {
	*dbschema.NgingCloudBackup
	Storage  *dbschema.NgingCloudStorage `db:"-,relation=id:dest_storage"`
	Watching bool
}

type CloudBackupExt struct {
	*dbschema.NgingCloudBackup
	Storage *dbschema.NgingCloudStorage `db:"-,relation=id:dest_storage"`
}
