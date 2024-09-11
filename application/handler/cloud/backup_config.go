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
	"context"
	"encoding/json"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"
	"github.com/webx-top/echo/param"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/background"
	"github.com/coscms/webcore/library/cloudbackup"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webcore/model"
)

func BackupConfigList(ctx echo.Context) error {
	if ctx.IsAjax() && ctx.Form(`checkStatus`) == `fullbackup` {
		id := ctx.Formx(`id`).Uint()
		data := ctx.Data()
		if id < 1 {
			data.SetError(ctx.NewError(code.InvalidParameter, ``).SetZone(`id`))
		} else {
			data.SetData(echo.H{`backuping`: fullBackupIsRunning(id)})
		}
		return ctx.JSON(data)
	}
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
	return ctx.Render(`cloud/backup`, common.Err(ctx, err))
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
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/cloud/backup`))
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
	ctx.Set(`fileSources`, GetFileSources())
	ctx.Set(`engineForms`, cloudbackup.Forms)
	ctx.Set(`activeURL`, `/cloud/backup`)
	return ctx.Render(`cloud/backup_edit`, err)
}

func BackupConfigEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudBackup(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
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
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/cloud/backup`))
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
	ctx.Set(`fileSources`, GetFileSources())
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
				log.Errorf(`failed to cloudbackup.LevelDB().RemoveDB(%v): %v`, m.Id, rerr.Error())
			}
		}
	}
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/cloud/backup`))
}

func BackupRestore(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudBackup(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if m.Disabled == common.BoolN {
		return ctx.NewError(code.DataStatusIncorrect, `必须停用后才能进行还原操作`)
	}
	if ctx.IsPost() {
		cfg := *m.NgingCloudBackup
		localSavePath := ctx.Formx(`localSavePath`).String()
		if len(localSavePath) == 0 {
			err = ctx.NewError(code.InvalidParameter, `请指定本机保存路径`)
			return err
		}
		actionIdent := `cloudbackup`
		bgKey := `restore.` + param.AsString(cfg.Id)
		bg := background.New(context.Background(), nil)
		group, err := background.Register(ctx, actionIdent, bgKey, bg)
		if err != nil {
			return err
		}
		finishMsg := ctx.T(`恭喜，文件恢复完毕`)
		user := backend.User(ctx)
		noticer := notice.NewP(ctx, actionIdent, user.Username, bg.Context()).AutoComplete(true)
		defer group.Cancel(bgKey)
		cfg.SourcePath = localSavePath
		err = cloudbackup.Restore(ctx, cfg, func(from, to string) {
			noticer.Send(from+` => `+to, notice.StateSuccess)
		}, noticer)
		if err != nil {
			noticer.Send(err.Error(), notice.StateFailure)
			common.SendErr(ctx, err)
		} else {
			noticer.Send(finishMsg, notice.StateSuccess)
			noticer.Complete()
			common.SendOk(ctx, finishMsg)
		}
		return ctx.Redirect(backend.URLFor(`/cloud/backup`))
	}

	ctx.Set(`title`, ctx.T(`还原备份文件`))
	ctx.Set(`data`, m.NgingCloudBackup)
	ctx.Set(`activeURL`, `/cloud/backup`)
	return ctx.Render(`cloud/backup_restore`, err)
}
