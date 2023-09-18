package model

import (
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/webx-top/echo"
)

const (
	// log type
	CloudBackupLogTypeAll   = `all`
	CloudBackupLogTypeError = `error`

	// storage engine
	StorageEngineS3     = `s3`
	StorageEngineSFTP   = `sftp`
	StorageEngineFTP    = `ftp`
	StorageEngineWebDAV = `webdav`
	StorageEngineSMB    = `smb`
)

var CloudBackupLogTypes = echo.NewKVData().Add(CloudBackupLogTypeAll, `全部`).Add(CloudBackupLogTypeError, `报错`)

var CloudBackupStorageEngines = echo.NewKVData()

type CloudBackupExt struct {
	*dbschema.NgingCloudBackup
	Storage       *dbschema.NgingCloudStorage `db:"-,relation=id:dest_storage|gtZero|eq(storage_engine:s3),columns=id&name&type"`
	Watching      bool
	FullBackuping bool
}
