/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package cloud

import (
	"encoding/json"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"
	"github.com/webx-top/echo/param"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/cloudbackup"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/model"
)

func BackupConfigList(ctx echo.Context) error {
	m := model.NewCloudBackup(ctx)
	cond := db.NewCompounds()
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`name`, db.Like(`%`+q+`%`))
	}
	list, err := m.ListPage(cond, `-id`)
	for _, row := range list {
		row.Watching = cloudbackup.BackupTasks.Has(row.Id)
		row.FullBackuping = fullBackupIsRunning(row.Id)
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`cloud/backup`, handler.Err(ctx, err))
}

func backupFormFilter(opts ...formfilter.Options) echo.FormDataFilter {
	opts = append(opts,
		formfilter.Exclude(`storageConfig`),
	)
	return formfilter.Build(opts...)
}

func getStorageConfig(ctx echo.Context, engineName string) string {
	storageConfig := echo.H{}
	for k, v := range ctx.Forms() {
		if !strings.HasPrefix(k, `storageConfig.`) || len(v) == 0 {
			continue
		}
		if !cloudbackup.HasForm(engineName, k) {
			continue
		}
		name := strings.TrimPrefix(k, `storageConfig.`)
		if name == `password` {
			v[0] = common.Crypto().Encode(v[0])
		}
		storageConfig.Set(name, v[0])
	}
	b, _ := json.Marshal(storageConfig)
	return string(b)
}

func setStorageConfigForm(ctx echo.Context, conf string) error {
	if len(conf) == 0 {
		return nil
	}
	storageConfig := echo.H{}
	err := json.Unmarshal([]byte(conf), &storageConfig)
	if err != nil {
		return err
	}
	for k, v := range storageConfig {
		value := param.AsString(v)
		if k == `password` {
			value = common.Crypto().Decode(value)
		}
		ctx.Request().Form().Set(`storageConfig.`+k, value)
	}
	return nil
}

func BackupConfigAdd(ctx echo.Context) error {
	var (
		err error
		id  uint
	)
	m := model.NewCloudBackup(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingCloudBackup, backupFormFilter())
		if err != nil {
			goto END
		}
		m.StorageConfig = getStorageConfig(ctx, m.StorageEngine)
		_, err = m.Add()
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(handler.URLFor(`/cloud/backup`))
	}
	id = ctx.Formx(`copyId`).Uint()
	if id > 0 {
		err = m.Get(nil, `id`, id)
		if err == nil {
			echo.StructToForm(ctx, m.NgingCloudBackup, ``, func(topName, fieldName string) string {
				return echo.LowerCaseFirstLetter(topName, fieldName)
			})
			setStorageConfigForm(ctx, m.StorageConfig)
			ctx.Request().Form().Set(`id`, `0`)
		}
	}

END:
	ctx.Set(`isAdd`, true)
	ctx.Set(`title`, ctx.T(`添加云备份配置`))
	ctx.Set(`engines`, model.CloudBackupStorageEngines.Slice())
	ctx.Set(`engineForms`, cloudbackup.Forms)
	ctx.Set(`activeURL`, `/cloud/backup`)
	return ctx.Render(`cloud/backup_edit`, err)
}

func BackupConfigEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudBackup(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingCloudBackup, backupFormFilter(formfilter.Exclude(`created`)))
		if err != nil {
			goto END
		}
		m.StorageConfig = getStorageConfig(ctx, m.StorageEngine)
		m.Id = id
		err = m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(handler.URLFor(`/cloud/backup`))
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			if !common.IsBoolFlag(disabled) {
				return ctx.NewError(code.InvalidParameter, ``).SetZone(`disabled`)
			}
			m.Disabled = disabled
			data := ctx.Data()
			err = m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err == nil {
				if m.Disabled == `Y` {
					err = allBackupStop(m.Id)
				}
			}
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
	}
	if err == nil {
		echo.StructToForm(ctx, m.NgingCloudBackup, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
		setStorageConfigForm(ctx, m.StorageConfig)
	}

END:
	ctx.Set(`isAdd`, false)
	ctx.Set(`title`, ctx.T(`修改云备份配置`))
	ctx.Set(`engines`, model.CloudBackupStorageEngines.Slice())
	ctx.Set(`engineForms`, cloudbackup.Forms)
	ctx.Set(`activeURL`, `/cloud/backup`)
	return ctx.Render(`cloud/backup_edit`, err)
}

func BackupConfigDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudBackup(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		if err = allBackupStop(m.Id); err == nil {
			if rerr := cloudbackup.LevelDB().RemoveDB(m.Id); rerr != nil {
				log.Error(rerr.Error())
			}
		}
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/cloud/backup`))
}
