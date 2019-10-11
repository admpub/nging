package initialize

import (
	"fmt"
	"io"

	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/common"
	modelFile "github.com/admpub/nging/application/model/file"
	"github.com/admpub/nging/application/registry/upload"
)

func init() {
	common.OnUpdateOwnerFilePath = OnUpdateOwnerFilePath
	common.OnRemoveOwnerFile = OnRemoveOwnerFile
	upload.DefaultDBSaver = func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) error {
		return fileM.Add(reader)
	}
	upload.DBSaverRegister(`user-avatar`, func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) (err error) {
		if len(fileM.TableId) == 0 {
			return fileM.Add(reader)
		}
		fileM.UsedTimes = 0
		m := &dbschema.File{}
		m.CPAFrom(fileM.File)
		err = m.Get(nil, db.And(
			db.Cond{`table_id`: fileM.TableId},
			db.Cond{`table_name`: fileM.TableName},
			db.Cond{`field_name`: fileM.FieldName},
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
		return fileM.Updater(`user`, `avatar`, fmt.Sprint(userM.Id)).Add(userM.Avatar, false)
	})
	dbschema.DBI.On(`user:updating`, func(m factory.Model, editColumns ...string) error {
		if len(editColumns) > 0 && !com.InSlice(`avatar`, editColumns) {
			return nil
		}
		fileM := modelFile.NewEmbedded(m.Context())
		userM := m.(*dbschema.User)
		return fileM.Updater(`user`, `avatar`, fmt.Sprint(userM.Id)).Edit(userM.Avatar, false)
	})
	dbschema.DBI.On(`user:deleted`, func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		userM := m.(*dbschema.User)
		return fileM.Updater(`user`, `avatar`, fmt.Sprint(userM.Id)).Delete()
	})

	// - config
	dbschema.DBI.On(`config:created`, func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		confM := m.(*dbschema.Config)
		embedded, seperator, exit := getConfigEventAttrs(confM)
		if exit {
			return nil
		}
		return fileM.Updater(`config`, `value`, confM.Group+`.`+confM.Key).SetSeperator(seperator).Add(confM.Value, embedded)
	})
	dbschema.DBI.On(`config:updating`, func(m factory.Model, editColumns ...string) error {
		if len(editColumns) > 0 && !com.InSlice(`value`, editColumns) {
			return nil
		}
		fileM := modelFile.NewEmbedded(m.Context())
		confM := m.(*dbschema.Config)
		embedded, seperator, exit := getConfigEventAttrs(confM)
		if exit {
			return nil
		}
		return fileM.Updater(`config`, `value`, confM.Group+`.`+confM.Key).SetSeperator(seperator).Edit(confM.Value, embedded)
	})
	dbschema.DBI.On(`config:deleted`, func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		confM := m.(*dbschema.Config)
		return fileM.Updater(`config`, `value`, confM.Group+`.`+confM.Key).Delete()
	})
}

func getConfigEventAttrs(confM *dbschema.Config) (embedded bool, seperator string, exit bool) {
	switch confM.Type {
	case `html`:
		embedded = true
	case `image`, `video`, `audio`, `file`:
	case `list`:
		seperator = `,`
	default:
		exit = true
	}
	return
}
