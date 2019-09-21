package upload

import (
	"io"

	uploadClient "github.com/webx-top/client/upload"
	modelFile "github.com/admpub/nging/application/model/file"
	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/db"
)

type (
	DBSaver func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) error
)

var (
	dbSavers = map[string]DBSaver{
		`user-avatar`: func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) error {
			if fileM.OwnerId <= 0 {
				return fileM.Add(reader)
			}
			fileM.UsedTimes = 1
			m := &dbschema.File{}
			err := m.Get(nil, db.And(
				db.Cond{`table_name`:`user-avatar`},
				db.Cond{`table_id`:fileM.OwnerId},
			))
			if err != nil {
				if err != db.ErrNoMoreRows {
					return err
				}
				// 不存在
				return fileM.Add(reader)
			}
			// 已存在
			fileM.Id = m.Id
			fileM.Created = m.Created
			fileM.FillData(reader, true)
			return fileM.Edit(nil, db.Cond{`id`: m.Id})
		},
	}
	DefaultDBSaver = func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) error {
		return fileM.Add(reader)
	}
)

func DBSaverRegister(key string,dbsaver DBSaver){
	dbSavers[key] = dbsaver
}

func DBSaverGet(key string) DBSaver {
	if dbsaver, ok := dbSavers[key]; ok {
		return dbsaver
	}
	return DefaultDBSaver
}