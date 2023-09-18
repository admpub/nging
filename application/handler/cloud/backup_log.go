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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/model"
)

func Log(ctx echo.Context) error {
	backupID := ctx.Formx(`backupId`).Uint()
	m := model.NewCloudBackupLog(ctx)
	page, size, totalRows, p := handler.PagingWithPagination(ctx)
	cond := db.Cond{}
	var backupM *model.CloudBackup
	var err error
	if backupID > 0 {
		backupM = model.NewCloudBackup(ctx)
		err = backupM.Get(nil, `id`, backupID)
		cond[`backup_id`] = backupID
	}
	cnt, err2 := m.List(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, page, size, cond)
	if err2 != nil {
		err = err2
	}
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	ctx.Set(`listData`, m.Objects())
	ctx.Set(`pagination`, p)
	if backupM == nil {
		backupM = model.NewCloudBackup(ctx)
	}
	ctx.Set(`data`, backupM)
	ctx.SetFunc(`getOperationName`, model.CloudBackupOperations.Get)
	ctx.SetFunc(`getBackupTypeName`, model.CloudBackupTypes.Get)
	ret := handler.Err(ctx, err)
	ctx.Set(`activeURL`, `/cloud/backup`)
	return ctx.Render(`cloud/backup_log`, ret)
}

func LogDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	backupID := ctx.Formx(`backupId`).Uint()
	m := model.NewCloudBackupLog(ctx)
	var (
		cond db.Cond
		err  error
		n    int
	)
	if id > 0 {
		cond = db.Cond{`id`: id}
	} else {
		ago := ctx.Form(`ago`)
		if len(ago) < 2 {
			handler.SendFail(ctx, ctx.T(`missing param`))
			goto END
		}

		switch ago[len(ago)-1] {
		case 'd': //删除几天前的。例如：7d
			n, err = strconv.Atoi(strings.TrimSuffix(ago, `d`))
			if err != nil {
				handler.SendFail(ctx, err.Error()+`:`+ago)
				goto END
			}

			cond = db.Cond{`created`: db.Lt(time.Now().AddDate(0, 0, -n).Unix())}
		case 'm': //删除几个月前的。例如：1m
			n, err = strconv.Atoi(strings.TrimSuffix(ago, `m`))
			if err != nil {
				handler.SendFail(ctx, err.Error()+`:`+ago)
				goto END
			}

			cond = db.Cond{`created`: db.Lt(time.Now().AddDate(0, -n, 0).Unix())}
		case 'y': //删除几年前的。例如：1y
			n, err = strconv.Atoi(strings.TrimSuffix(ago, `y`))
			if err != nil {
				handler.SendFail(ctx, err.Error()+`:`+ago)
				goto END
			}

			cond = db.Cond{`created`: db.Lt(time.Now().AddDate(-n, 0, 0).Unix())}
		default:
			handler.SendFail(ctx, ctx.T(`invalid param`))
			goto END
		}
		if backupID > 0 {
			cond[`backup_id`] = backupID
		}
	}
	err = m.Delete(nil, cond)
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

END:
	return ctx.Redirect(handler.URLFor(`/cloud/backup_log?backupId=`) + fmt.Sprint(backupID))
}
