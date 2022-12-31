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

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/admpub/archiver"
	loga "github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/background"
	"github.com/admpub/nging/v5/application/library/notice"

	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver/mysql/utils"
)

// SQLTempDir sql文件缓存目录获取函数(用于导入导出SQL)
var SQLTempDir = os.TempDir

func TempDir(op string) string {
	dir := filepath.Join(SQLTempDir(), `dbmanager/cache`, op)
	err := com.MkdirAll(dir, os.ModePerm)
	if err != nil {
		loga.Error(err)
	}
	return dir
}

func (m *mySQL) getCharsetList() []string {
	e, _ := m.getCharsets()
	cs := make([]string, 0)
	for k := range e {
		cs = append(cs, k)
	}
	sort.Strings(cs)
	return cs
}

func (m *mySQL) exporting() error {
	return m.bgExecManage(utils.OpExport)
}

func (m *mySQL) ImportAndOutputOpName(op string) string {
	return `dbmanager.` + m.DbAuth.Driver + `.` + op
}

func (m *mySQL) bgExecManage(op string) error {
	var err error
	if m.IsPost() {
		data := m.Data()
		keys := m.FormValues(`key`)
		background.Cancel(m.ImportAndOutputOpName(op), keys...)
		data.SetInfo(m.T(`操作成功`))
		m.ok(m.T(`操作成功`))
		return m.returnTo(m.GenURL(op) + `&process=1`)
	}
	m.Set(`op`, op)
	group := background.ListBy(m.ImportAndOutputOpName(op))
	bgs := map[string]background.Background{}
	if group != nil {
		bgs = group.Map()
	}
	m.Set(`list`, bgs)
	var title string
	if op == utils.OpExport {
		title = m.T(`导出SQL`)
	} else {
		title = m.T(`导入SQL`)
	}
	m.Set(`title`, title)
	m.Set(`cacheDir`, handler.URLFor(`/download/file?path=dbmanager/cache/`+op))
	return m.Render(`db/mysql/process_store`, m.checkErr(err))
}

func (m *mySQL) Export() error {
	if len(m.Form(`docType`)) > 0 {
		return m.ExportDoc()
	}
	//fmt.Printf("%#v\n", m.getCharsetList())
	process := m.Queryx(`process`).Bool()
	if process {
		return m.exporting()
	}
	var err error
	if m.IsPost() {
		var tables []string
		if len(m.dbName) == 0 {
			m.fail(m.T(`请选择数据库`))
			return m.returnTo(m.GenURL(`listDb`))
		}
		if m.Formx(`all`).Bool() {
			var ok bool
			tables, ok = m.Get(`tableList`).([]string)
			if !ok {
				tables, err = m.getTables()
				if err != nil {
					m.fail(err.Error())
					return m.returnTo(m.GenURL(`export`))
				}
			}
		} else {
			tables = m.FormValues(`table`)
			if len(tables) == 1 && len(tables[0]) > 0 {
				tables = strings.Split(tables[0], `,`)
			}
			views := m.FormValues(`view`)
			if len(views) == 1 && len(views[0]) > 0 {
				views = strings.Split(views[0], `,`)
			}
			if len(views) > 0 {
				tables = append(tables, views...)
			}
		}
		output := m.Form(`output`)
		types := m.FormValues(`type`)
		cacheKey := com.Md5(com.Dump([]interface{}{tables, output, types}, false))
		var (
			structWriter, dataWriter interface{}
			sqlFiles                 []string
			dbSaveDir                string
			async                    bool
			bgExec                   = background.New(context.TODO(), echo.H{
				`database`: m.dbName,
				`tables`:   tables,
				`output`:   output,
				`types`:    types,
			})
			fileInfos = &utils.FileInfos{}
		)
		exports, err := background.Register(m.Context, m.ImportAndOutputOpName(utils.OpExport), cacheKey, bgExec)
		if err != nil {
			return err
		}
		nowTime := time.Now().Format("20060102150405.000")
		saveDir := TempDir(utils.OpExport)
		switch output {
		case `down`:
			m.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
			m.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", m.dbName+"-sql-"+nowTime+".sql"))
			fallthrough
		case `open`:
			m.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
			if com.InSlice(`struct`, types) {
				structWriter = m.Response()
			}
			if com.InSlice(`data`, types) {
				dataWriter = m.Response()
			}
		default:
			async = true
			dbSaveDir = filepath.Join(saveDir, m.dbName)
			com.MkdirAll(dbSaveDir, os.ModePerm)
			if com.InSlice(`struct`, types) {
				structFile := filepath.Join(dbSaveDir, `struct-`+nowTime+`.sql`)
				sqlFiles = append(sqlFiles, structFile)
				structWriter = structFile
				fi := &utils.FileInfo{
					Start: time.Now(),
					Path:  structFile,
				}
				*fileInfos = append(*fileInfos, fi)
			}
			if com.InSlice(`data`, types) {
				dataFile := filepath.Join(dbSaveDir, `data-`+nowTime+`.sql`)
				sqlFiles = append(sqlFiles, dataFile)
				dataWriter = dataFile
				fi := &utils.FileInfo{
					Start: time.Now(),
					Path:  dataFile,
				}
				*fileInfos = append(*fileInfos, fi)
			}
		}
		cfg := *m.DbAuth
		cfg.Db = m.dbName
		coll, err := m.getCollation(m.dbName, nil)
		if err != nil {
			return err
		}
		cfg.Charset = strings.SplitN(coll, `_`, 2)[0]

		user := handler.User(m.Context)
		var username string
		if user != nil {
			username = user.Username
		}
		noticer := notice.New(m.Context, `databaseExport`, username, bgExec.Context())
		gtidMode, _ := m.showVariables(`gtid_mode`)
		var hasGTID bool
		if len(gtidMode) > 0 && len(gtidMode[0]) > 0 {
			if k, y := gtidMode[0][`k`]; y && len(k) > 0 {
				hasGTID = true
			}
		}

		worker := func(c context.Context, cfg driver.DbAuth) error {
			defer func() {
				exports.Cancel(cacheKey)
				if r := recover(); r != nil {
					err = fmt.Errorf(`RECOVER: %v`, r)
				}
			}()
			if utils.SupportedExport() { // 采用 mysqldump 命令导出
				err = utils.Export(c, noticer, &cfg, tables, structWriter, dataWriter, m.getVersion(), hasGTID, true)
			} else {
				if structWriter != nil {
					err = m.exportDBStruct(c, noticer, &cfg, tables, structWriter, m.getVersion(), true)
				}
				if err == nil && dataWriter != nil {
					err = m.exportDBData(c, noticer, &cfg, tables, dataWriter, m.getVersion())
				}
			}
			if err != nil {
				loga.Error(err)
				return err
			}
			if len(sqlFiles) > 0 {
				now := time.Now()
				for _, fi := range *fileInfos {
					fi.End = now
					fi.Size, err = com.FileSize(fi.Path)
					if err != nil {
						fi.Error = err.Error()
					}
					fi.Elapsed = fi.End.Sub(fi.Start)
				}
				zipFile := filepath.Join(dbSaveDir, "sql-"+nowTime+".zip")
				fi := &utils.FileInfo{
					Start:      now,
					Path:       zipFile,
					Compressed: true,
				}
				err = archiver.Archive(sqlFiles, zipFile)
				if err != nil {
					loga.Error(err)
					return err
				}
				for _, sqlFile := range sqlFiles {
					os.Remove(sqlFile)
				}
				fi.Size, err = com.FileSize(zipFile)
				if err != nil {
					fi.Error = err.Error()
				}
				fi.End = time.Now()
				fi.Elapsed = fi.End.Sub(fi.Start)
				fileInfos.Add(fi)
				os.WriteFile(zipFile+`.txt`, com.Str2bytes(com.Dump(fileInfos, false)), os.ModePerm)
			}
			return nil
		}
		if !async {
			done := make(chan struct{})
			ctx := m.StdContext()
			go func() {
				defer exports.Cancel(cacheKey)
				for {
					select {
					case <-ctx.Done():
						return
					case <-done:
						return
					}
				}
			}()
			err = worker(bgExec.Context(), cfg)
			if err != nil {
				noticer(m.T(`导出失败`)+`: `+err.Error(), 0)
			}
			done <- struct{}{}
			return err
		}
		data := m.Data()
		data.SetInfo(m.T(`任务已经在后台成功启动`))
		data.SetURL(handler.URLFor(`/download/file?path=dbmanager/cache/` + utils.OpExport + `/` + m.dbName))
		go worker(bgExec.Context(), cfg)
		return m.JSON(data)
	}
	return m.Redirect(m.GenURL(`listTable`, m.dbName))
}

func (m *mySQL) ExportDoc() error {
	if m.IsPost() {
		var tables []string
		if len(m.dbName) == 0 {
			m.fail(m.T(`请选择数据库`))
			return m.returnTo(m.GenURL(`listDb`))
		}
		tables = m.FormValues(`table`)
		if len(tables) == 1 && len(tables[0]) > 0 {
			tables = strings.Split(tables[0], `,`)
		}
		docType := m.Form(`docType`)
		newExportorDoc, ok := docExportors[docType]
		if !ok {
			return m.NewError(code.InvalidParameter, `不支持导出文档类型: %s`, docType)
		}
		exportor := newExportorDoc(m.dbName)
		err := exportor.Open(m.Context)
		if err != nil {
			return err
		}
		for _, table := range tables {
			origFields, sortFields, err := m.tableFields(table)
			if err != nil {
				return err
			}
			stt, _, err := m.getTableStatus(m.dbName, table, false)
			if err != nil {
				return err
			}
			var tableStatus *TableStatus
			if ts, ok := stt[table]; ok {
				tableStatus = ts
			} else {
				tableStatus = &TableStatus{Name: sql.NullString{Valid: true, String: table}}
			}
			postFields := make([]*Field, len(sortFields))
			for k, v := range sortFields {
				postFields[k] = origFields[v]
			}
			err = exportor.Write(m.Context, tableStatus, postFields)
			if err != nil {
				return err
			}
		}
		return exportor.Close(m.Context)
	}
	return m.Redirect(m.GenURL(`listTable`, m.dbName))
}

type DocExportor interface {
	Open(echo.Context) error
	Write(echo.Context, *TableStatus, []*Field) error
	Close(echo.Context) error
}

var docExportors = map[string]func(dbName string) DocExportor{
	`html`:     newHTMLDocExportor,
	`markdown`: newMarkdownDocExportor,
	`csv`:      newCSVDocExportor,
}
