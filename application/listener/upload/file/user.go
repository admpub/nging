/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package file

import (
	"fmt"
	"io"
	"time"

	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	modelFile "github.com/admpub/nging/application/model/file"
	_ "github.com/admpub/nging/application/model/file/initialize"
	"github.com/admpub/nging/application/registry/upload"
	"github.com/admpub/nging/application/registry/upload/table"
)

func init() {
	// 用户上传个人文件时的文件命名方式
	upload.CheckerRegister(`user`, func(ctx echo.Context, tis table.TableInfoStorer) (subdir string, name string, err error) {
		user := handler.User(ctx)
		if user == nil {
			err = ctx.E(`登录信息获取失败，请重新登录`)
			return
		}
		userID := uint64(user.Id)
		timestamp := ctx.Formx(`time`).Int64()
		// 验证签名（避免上传接口被滥用）
		if ctx.Form(`token`) != upload.Token(ctx.Queries()) {
			err = ctx.E(`令牌错误`)
			return
		}
		if time.Now().Local().Unix()-timestamp > upload.UploadLinkLifeTime {
			err = ctx.E(`上传网址已过期`)
			return
		}
		uid := fmt.Sprint(userID)
		subdir = uid + `/`
		subdir += time.Now().Format(`2006/01/02/`)
		tis.SetTableID(uid)
		tis.SetTableName(`user`)
		tis.SetFieldName(``)
		return
	})

	// 文件信息默认保存方式
	upload.DefaultDBSaver = func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) error {
		return fileM.Add(reader)
	}

	// 后台用户头像文件信息保存方式
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

	// - 监听数据表事件

	// - user 表事件
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

	// - config 表事件
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
