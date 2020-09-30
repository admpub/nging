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

package cloud

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/model"
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
		row.Watching = backupTasks.Has(row.Id)
		row.FullBackuping = fullBackupIsRunning(row.Id)
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`cloud/backup`, handler.Err(ctx, err))
}

func BackupConfigAdd(ctx echo.Context) error {
	var (
		err error
		id  uint
	)
	m := model.NewCloudBackup(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingCloudBackup)
		if err != nil {
			goto END
		}
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
			ctx.Request().Form().Set(`id`, `0`)
		}
	}

END:
	ctx.Set(`isAdd`, true)
	ctx.Set(`title`, ctx.T(`添加云备份配置`))
	ctx.Set(`activeURL`, `/cloud/backup`)
	return ctx.Render(`cloud/backup_edit`, err)
}

func BackupConfigEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudBackup(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingCloudBackup, echo.ExcludeFieldName(`created`))
		if err != nil {
			goto END
		}
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
			m.Disabled = disabled
			data := ctx.Data()
			err = m.Edit(nil, db.Cond{`id`: id})
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
	}

END:
	ctx.Set(`isAdd`, false)
	ctx.Set(`title`, ctx.T(`修改云备份配置`))
	ctx.Set(`activeURL`, `/cloud/backup`)
	return ctx.Render(`cloud/backup_edit`, err)
}

func BackupConfigDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudBackup(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/cloud/backup`))
}
