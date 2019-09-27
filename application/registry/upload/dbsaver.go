package upload

import (
	"io"

	modelFile "github.com/admpub/nging/application/model/file"
	uploadClient "github.com/webx-top/client/upload"
)

type (
	DBSaver func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) error
)

var (
	dbSavers       = map[string]DBSaver{}
	DefaultDBSaver = func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) error {
		return nil
	}
)

func DBSaverRegister(key string, dbsaver DBSaver) {
	dbSavers[key] = dbsaver
}

func DBSaverGet(key string) DBSaver {
	if dbsaver, ok := dbSavers[key]; ok {
		return dbsaver
	}
	return DefaultDBSaver
}
