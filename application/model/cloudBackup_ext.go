package model

import (
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/webx-top/echo"
)

const (
	// log type
	CloudBackupLogTypeAll   = `all`
	CloudBackupLogTypeError = `error`
)

var CloudBackupLogTypes = echo.NewKVData().Add(CloudBackupLogTypeAll, `全部`).Add(CloudBackupLogTypeError, `报错`)

type CloudBackupExt struct {
	*dbschema.NgingCloudBackup
	Storage       *dbschema.NgingCloudStorage `db:"-,relation=id:dest_storage"`
	Watching      bool
	FullBackuping bool
}
