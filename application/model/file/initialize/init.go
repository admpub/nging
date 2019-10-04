package initialize

import (
	"io"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/common"
	modelFile "github.com/admpub/nging/application/model/file"
	"github.com/admpub/nging/application/registry/upload"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
)

func init() {
	common.OnUpdateOwnerFilePath = OnUpdateOwnerFilePath
	common.OnRemoveOwnerFile = OnRemoveOwnerFile
	upload.DefaultDBSaver = func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) error {
		return fileM.Add(reader)
	}
	upload.DBSaverRegister(`user-avatar`, func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) (err error) {
		if fileM.OwnerId <= 0 {
			return fileM.Add(reader)
		}
		fileM.UsedTimes = 0
		m := &dbschema.File{}
		m.CPAFrom(fileM.File)
		err = m.Get(nil, db.And(
			db.Cond{`table_id`: fileM.OwnerId},
			db.Cond{`table_name`: `user`},
			db.Cond{`field_name`: `avatar`},
		))
		defer func() {
			if err != nil {
				return
			}
			userM := &dbschema.User{}
			userM.CPAFrom(fileM.File)
			err = userM.SetField(nil, `avatar`, fileM.ViewUrl, db.Cond{`id`: fileM.OwnerId})
		}()
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
		fileM.UsedTimes = m.UsedTimes
		fileM.FillData(reader, true)
		return fileM.Edit(nil, db.Cond{`id`: m.Id})
	})

	// - user
	dbschema.DBI.On(`user:created`, func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		userM := m.(*dbschema.User)
		return fileM.Updater(`user`, `avatar`, uint64(userM.Id)).Add(userM.Avatar, false)
	})
	dbschema.DBI.On(`user:updating`, func(m factory.Model, editColumns ...string) error {
		if len(editColumns) > 0 && !com.InSlice(`avatar`, editColumns) {
			return nil
		}
		fileM := modelFile.NewEmbedded(m.Context())
		userM := m.(*dbschema.User)
		return fileM.Updater(`user`, `avatar`, uint64(userM.Id)).Edit(userM.Avatar, false)
	})
	dbschema.DBI.On(`user:deleted`, func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		userM := m.(*dbschema.User)
		return fileM.Updater(`user`, `avatar`, uint64(userM.Id)).Delete()
	})

	// - config
	dbschema.DBI.On(`config:created`, func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		confM := m.(*dbschema.Config)
		return fileM.Updater(`config`, confM.Group+`.`+confM.Key, 0).Add(confM.Value, true)
	})
	dbschema.DBI.On(`config:updating`, func(m factory.Model, editColumns ...string) error {
		if len(editColumns) > 0 && !com.InSlice(`value`, editColumns) {
			return nil
		}
		fileM := modelFile.NewEmbedded(m.Context())
		confM := m.(*dbschema.Config)
		return fileM.Updater(`config`, confM.Group+`.`+confM.Key, 0).Edit(confM.Value, true)
	})
	dbschema.DBI.On(`config:deleted`, func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		confM := m.(*dbschema.Config)
		return fileM.Updater(`config`, confM.Group+`.`+confM.Key, 0).Delete()
	})
}
