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

package database

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/mysql-schema-sync/sync"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/cron"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/db`, func(g echo.RouteRegister) {
		e := handler.Echo()
		g.Route(`GET,POST`, `/schema_sync`, e.MetaHandler(echo.H{`name`: `同步方案列表`}, SchemaSync))
		g.Route(`GET,POST`, `/schema_sync_add`, e.MetaHandler(echo.H{`name`: `添加同步方案`}, SchemaSyncAdd))
		g.Route(`GET,POST`, `/schema_sync_edit`, e.MetaHandler(echo.H{`name`: `编辑同步方案`}, SchemaSyncEdit))
		g.Route(`GET`, `/schema_sync_delete`, e.MetaHandler(echo.H{`name`: `删除同步方案`}, SchemaSyncDelete))
		g.Route(`GET`, `/schema_sync_preview`, e.MetaHandler(echo.H{`name`: `预览要同步的项`}, SchemaSyncPreview))
		g.Route(`GET`, `/schema_sync_run`, e.MetaHandler(echo.H{`name`: `执行同步方案`}, SchemaSyncRun))
		g.Route(`GET`, `/schema_sync_log/:id`, e.MetaHandler(echo.H{`name`: `日志列表`}, SchemaSyncLog))
		g.Route(`GET`, `/schema_sync_log_view/:id`, e.MetaHandler(echo.H{`name`: `日志详情`}, SchemaSyncLogView))
		g.Route(`GET`, `/schema_sync_log_delete`, e.MetaHandler(echo.H{`name`: `删除日志`}, SchemaSyncLogDelete))
	})
	cron.AddSYSJob(`mysql_schema_sync`, SchemaSyncJob, `>mysql_schema_sync:1`, `同步MySQL数据表结构`)
}

//SchemaSyncJob 计划任务调用方式
func SchemaSyncJob(id string) cron.Runner {
	return func(timeout time.Duration) (out string, runingErr string, onErr error, isTimeout bool) {
		m := model.NewDbSync(nil)
		err := m.Get(nil, db.Cond{`id`: id})
		if err == nil {
			var result *dbschema.DbSyncLog
			result, err = execSync(m, false)
			if result.Failed > 0 {
				onErr = cron.ErrFailure
				runingErr = fmt.Sprintf(`有%d个错误`, result.Failed)
			}
		}
		if err != nil {
			runingErr = err.Error()
		} else {
			out = `日志详情：/db/schema_sync_log/` + id
		}
		return
	}
}

//SchemaSync 表结构同步
func SchemaSync(ctx echo.Context) error {
	m := model.NewDbSync(ctx)
	page, size, totalRows, p := handler.PagingWithPagination(ctx)
	cnt, err := m.List(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, page, size)
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	ret := handler.Err(ctx, err)
	ctx.Set(`pagination`, p)
	listData := m.Objects()
	for k, v := range listData {
		v.DsnSource = m.HidePassword(v.DsnSource)
		v.DsnDestination = m.HidePassword(v.DsnDestination)
		listData[k] = v
	}
	ctx.Set(`listData`, listData)
	ctx.Set(`activeURL`, `/db/schema_sync`)
	return ctx.Render(`db/schema_sync`, ret)
}

//SchemaSyncAdd 添加表结构同步方案
func SchemaSyncAdd(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := model.NewDbSync(ctx)
		err = ctx.MustBind(m.DbSync)
		if err == nil {
			user := m.Form(`dsn_source_user`)
			passwd := m.Form(`dsn_source_passwd`)
			host := m.Form(`dsn_source_host`)
			dbName := m.Form(`dsn_source_database`)
			m.DsnSource = m.ToDSN(user, passwd, host, dbName)
			user = m.Form(`dsn_destination_user`)
			passwd = m.Form(`dsn_destination_passwd`)
			host = m.Form(`dsn_destination_host`)
			dbName = m.Form(`dsn_destination_database`)
			m.DsnDestination = m.ToDSN(user, passwd, host, dbName)
			_, err = m.Add()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/db/schema_sync`))
			}
		}
	}
	id := ctx.Formx(`copyId`).Uint()
	if id > 0 {
		m := model.NewDbSync(ctx)
		cond := db.Cond{`id`: id}
		err = m.Get(nil, cond)
		if err == nil {
			setFormData(ctx, m)
			ctx.Request().Form().Set(`id`, `0`)
		}
	}
	ctx.Set(`activeURL`, `/db/schema_sync`)
	return ctx.Render(`db/schema_sync_edit`, handler.Err(ctx, err))
}

func setFormData(ctx echo.Context, m *model.DbSync) {
	echo.StructToForm(ctx, m.DbSync, ``, echo.LowerCaseFirstLetter)
	user, passwd, host, dbName := m.ParseDSN(m.DsnSource)
	ctx.Request().Form().Set(`dsn_source_user`, user)
	ctx.Request().Form().Set(`dsn_source_passwd`, passwd)
	ctx.Request().Form().Set(`dsn_source_host`, host)
	ctx.Request().Form().Set(`dsn_source_database`, dbName)
	user, passwd, host, dbName = m.ParseDSN(m.DsnDestination)
	ctx.Request().Form().Set(`dsn_destination_user`, user)
	ctx.Request().Form().Set(`dsn_destination_passwd`, passwd)
	ctx.Request().Form().Set(`dsn_destination_host`, host)
	ctx.Request().Form().Set(`dsn_destination_database`, dbName)
}

//SchemaSyncEdit 编辑表结构同步方案
func SchemaSyncEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewDbSync(ctx)
	cond := db.Cond{`id`: id}
	err = m.Get(nil, cond)
	if ctx.IsPost() {
		err = ctx.MustBind(m.DbSync, func(k string, v []string) (string, []string) {
			switch strings.ToLower(k) {
			case `created`: //禁止修改创建时间
				return ``, v
			}
			return k, v
		})
		if err == nil {
			user := m.Form(`dsn_source_user`)
			passwd := m.Form(`dsn_source_passwd`)
			host := m.Form(`dsn_source_host`)
			dbName := m.Form(`dsn_source_database`)
			m.DsnSource = m.ToDSN(user, passwd, host, dbName)
			user = m.Form(`dsn_destination_user`)
			passwd = m.Form(`dsn_destination_passwd`)
			host = m.Form(`dsn_destination_host`)
			dbName = m.Form(`dsn_destination_database`)
			m.DsnDestination = m.ToDSN(user, passwd, host, dbName)
			err = m.Edit(nil, cond)
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/db/schema_sync`))
			}
		}
	} else if err == nil {
		setFormData(ctx, m)
	}
	ctx.Set(`activeURL`, `/db/schema_sync`)
	return ctx.Render(`db/schema_sync_edit`, handler.Err(ctx, err))
}

func SchemaSyncDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewDbSync(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		logM := model.NewDbSyncLog(ctx)
		err := logM.Delete(nil, db.Cond{`sync_id`: id})
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
		} else {
			handler.SendFail(ctx, err.Error())
		}
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/db/schema_sync`))
}

func execSync(a *model.DbSync, preview bool) (*dbschema.DbSyncLog, error) {
	mc := &sync.EmailConfig{
		From: config.DefaultConfig.Email.From,
		To:   a.DbSync.MailTo,
	}
	if config.DefaultConfig.Email.SMTPConfig != nil {
		mc.SMTPHost = config.DefaultConfig.Email.SMTPConfig.Address()
		mc.Password = config.DefaultConfig.Email.SMTPConfig.Password
		mc.On = len(a.DbSync.MailTo) > 0
	}
	logM := model.NewDbSyncLog(a.Base.Context)
	logM.SyncId = a.Id
	r, err := sync.Sync(&sync.Config{
		Sync:        preview == false,
		Drop:        a.DbSync.Drop > 0,
		SourceDSN:   a.DbSync.DsnSource,
		DestDSN:     a.DbSync.DsnDestination,
		AlterIgnore: a.DbSync.AlterIgnore,
		Tables:      a.DbSync.Tables,
		SkipTables:  a.DbSync.SkipTables,
		MailTo:      a.DbSync.MailTo,
	}, mc)
	if err != nil {
		return logM.DbSyncLog, err
	}
	result := r.Diff(false).String()
	logM.Result = result
	logM.ChangeTableNum = uint(r.ChangeTableNum())
	logM.ChangeTables = strings.Join(r.ChangeTables(), `,`)
	logM.Failed = uint(r.FailedNum())
	logM.Elapsed = uint64(r.Elapsed().Seconds())
	if !preview {
		logM.Add()
	}
	return logM.DbSyncLog, err
}

func SchemaSyncPreview(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewDbSync(ctx)
	err := m.Get(nil, db.Cond{`id`: id})
	var previewData string
	if err != nil {
		handler.SendFail(ctx, err.Error())
		previewData = err.Error()
	} else {
		var result *dbschema.DbSyncLog
		result, err = execSync(m, true)
		previewData = result.Result
	}

	ctx.Set(`previewData`, previewData)
	ctx.Set(`activeURL`, `/db/schema_sync`)
	return ctx.Render(`db/schema_sync_preview`, handler.Err(ctx, err))
}

func SchemaSyncRun(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewDbSync(ctx)
	err := m.Get(nil, db.Cond{`id`: id})
	var previewData string
	if err != nil {
		handler.SendFail(ctx, err.Error())
		previewData = err.Error()
	} else {
		var result *dbschema.DbSyncLog
		result, err = execSync(m, false)
		previewData = result.Result
	}

	ctx.Set(`previewData`, previewData)
	ctx.Set(`activeURL`, `/db/schema_sync`)
	return ctx.Render(`db/schema_sync_preview`, handler.Err(ctx, err))
}

func SchemaSyncLog(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint()
	syncM := model.NewDbSync(ctx)
	err := syncM.Get(nil, db.Cond{`id`: id})
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/db/schema_sync`))
	}
	syncM.DbSync.DsnSource = syncM.HidePassword(syncM.DbSync.DsnSource)
	syncM.DbSync.DsnDestination = syncM.HidePassword(syncM.DbSync.DsnDestination)
	ctx.Set(`data`, syncM.DbSync)

	m := model.NewDbSyncLog(ctx)
	page, size, totalRows, p := handler.PagingWithPagination(ctx)
	cnt, err := m.List(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, page, size, `sync_id`, id)
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	ret := handler.Err(ctx, err)
	ctx.Set(`pagination`, p)
	ctx.Set(`listData`, m.Objects())
	ctx.Set(`activeURL`, `/db/schema_sync`)
	return ctx.Render(`db/schema_sync_log`, ret)
}

func SchemaSyncLogView(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint()
	m := model.NewDbSyncLog(ctx)
	err := m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/db/schema_sync`))
	}
	ctx.Set(`data`, m.DbSyncLog)
	ctx.Set(`activeURL`, `/db/schema_sync`)
	return ctx.Render(`db/schema_sync_log_view`, handler.Err(ctx, err))
}

func SchemaSyncLogDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	syncId := ctx.Formx(`syncId`).Uint()
	m := model.NewDbSyncLog(ctx)
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
		if syncId > 0 {
			cond[`sync_id`] = syncId
		}
	}
	err = m.Delete(nil, cond)
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

END:
	return ctx.Redirect(handler.URLFor(`/db/schema_sync_log/`) + fmt.Sprint(syncId))
}
