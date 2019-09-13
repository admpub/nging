package manager

import (
	"fmt"

	"github.com/admpub/events"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/registry/upload"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func init() {
	config.Emitter.On(`user-file-deleted`, events.Callback(func(e events.Event) error {
		data := e.Context.Get(`data`).(*dbschema.File)
		ownerID := e.Context.Uint64(`ownerID`)
		userM := &dbschema.User{}
		err := userM.Get(nil, db.Cond{`id`: ownerID})
		if err != nil {
			if err == db.ErrNoMoreRows {
				return nil
			}
			return err
		}
		err = userM.SetFields(nil, map[string]interface{}{
			`file_size`: db.Raw(`file_size-` + fmt.Sprintf(`%d`, data.Size)),
			`file_num`:  db.Raw(`file_num-1`),
		}, db.And(
			db.Cond{`id`: ownerID},
			db.Cond{`file_size`: db.Gte(data.Size)},
			db.Cond{`file_num`: db.Gt(0)},
		))
		if err != nil {
			fileM := &dbschema.File{}
			recv := echo.H{}
			err = fileM.NewParam().SetMW(func(r db.Result) db.Result {
				return r.Select(db.Raw(`SUM(size) AS c`), db.Raw(`COUNT(1) AS n`))
			}).SetRecv(&recv).SetArgs(db.And(
				db.Cond{`owner_type`: `user`},
				db.Cond{`owner_id`: ownerID},
			)).One()
			if err != nil {
				return err
			}
			totalNum := recv.Uint64(`n`)
			totalSize := recv.Uint64(`c`)
			err = userM.SetFields(nil, map[string]interface{}{
				`file_size`: totalSize,
				`file_num`:  totalNum,
			}, db.Cond{`id`: ownerID})
		}
		return err
	}))
	config.Emitter.On(`file-deleted`, events.Callback(func(e events.Event) error {
		ctx := e.Context.Get(`ctx`).(echo.Context)
		files := e.Context.Get(`files`).([]string)
		data := e.Context.Get("data").(*dbschema.File)
		newStore := upload.StorerGet(data.StorerName)
		if newStore == nil {
			return ctx.E(`存储引擎“%s”未被登记`, data.StorerName)
		}
		storer := newStore(``)
		var errs common.Errors
		for _, file := range files {
			if err := storer.Delete(file); err != nil {
				errs = append(errs, err)
			}
		}
		return errs
	}))
}
