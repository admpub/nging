package upload

import (
	"io"

	uploadClient "github.com/webx-top/client/upload"

	modelFile "github.com/admpub/nging/application/model/file"
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

func DBSaverUnregister(keys ...string) {
	for _, key := range keys {
		_, ok := dbSavers[key]
		if ok {
			delete(dbSavers, key)
		}
	}
}

func DBSaverGet(key string) DBSaver {
	if dbsaver, ok := dbSavers[key]; ok {
		return dbsaver
	}
	return DefaultDBSaver
}
