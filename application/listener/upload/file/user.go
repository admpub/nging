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

// Package file 监听文件更改，并记录到file表
package file

import (
	"fmt"
	"io"
	"time"

	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/fileupdater/listener"
	modelFile "github.com/admpub/nging/application/model/file"
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
	}, ``)

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
		var m *dbschema.NgingFile
		m, err = fileM.GetAvatar()
		defer func() {
			if err != nil {
				return
			}
			userID := param.AsUint64(fileM.TableId)
			if userID == 0 { // 新增用户时
				return
			}
			userM := &dbschema.NgingUser{}
			userM.CPAFrom(fileM.NgingFile)
			err = userM.SetField(nil, `avatar`, fileM.ViewUrl, db.Cond{`id`: userID})
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
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.NgingUser)
		tableID = fmt.Sprint(fm.Id)
		content = fm.Avatar
		property = listener.NewProUP(fm, db.Cond{`id`: fm.Id})
		return
	}, false).SetTable(`user`, `avatar`).ListenDefault()

	// - config 表事件
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.NgingConfig)
		tableID = fm.Group + `.` + fm.Key
		content = fm.Value
		property = getConfigEventAttrs(fm).GenUpdater(fm, db.And(
			db.Cond{`key`: fm.Key},
			db.Cond{`group`: fm.Group},
		))
		return
	}, false).SetTable(`config`, `value`).ListenDefault()
}

func getConfigEventAttrs(confM *dbschema.NgingConfig) *listener.Property {
	property := &listener.Property{}
	switch confM.Type {
	case `html`:
		property.SetEmbedded(true)
	case `image`, `video`, `audio`, `file`:
	case `list`:
		property.SetSeperator(`,`)
	default:
		property.SetExit(true)
	}
	return property
}
