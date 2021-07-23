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

package collector

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/handler"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func ExportLog(ctx echo.Context) error {
	exportID := ctx.Formx(`exportId`).Uint()
	totalRows := ctx.Formx(`rows`).Int()
	m := &dbschema.NgingCollectorExportLog{}
	page, size, totalRows, p := handler.PagingWithPagination(ctx)
	cond := db.Cond{}
	var export *dbschema.NgingCollectorExport
	var err error
	if exportID > 0 {
		export = &dbschema.NgingCollectorExport{}
		err = export.Get(nil, `id`, exportID)
		cond[`export_id`] = exportID
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
	if export == nil {
		export = &dbschema.NgingCollectorExport{}
	}
	ctx.Set(`export`, export)
	ret := handler.Err(ctx, err)
	ctx.Set(`activeURL`, `/collector/export`)
	return ctx.Render(`collector/export_log`, ret)
}

func renderLogViewData(ctx echo.Context, m *dbschema.NgingCollectorExportLog, err error) error {
	ctx.Set(`data`, m)
	ctx.Set(`activeURL`, `/collector/export`)
	var export *dbschema.NgingCollectorExport
	if m.ExportId > 0 {
		export = &dbschema.NgingCollectorExport{}
		err = export.Get(nil, `id`, m.ExportId)
	}
	ctx.Set(`export`, export)
	return ctx.Render(`collector/export_log_view`, handler.Err(ctx, err))
}

func ExportLogView(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint()
	m := &dbschema.NgingCollectorExportLog{}
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/collector/export_log`))
	}
	return renderLogViewData(ctx, m, err)
}

func ExportLogDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	exportId := ctx.Formx(`exportId`).Uint()
	m := &dbschema.NgingCollectorExportLog{}
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
		if exportId > 0 {
			cond[`export_id`] = exportId
		}
	}
	err = m.Delete(nil, cond)
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

END:
	return ctx.Redirect(handler.URLFor(`/collector/export_log?exportId=`) + fmt.Sprint(exportId))
}
