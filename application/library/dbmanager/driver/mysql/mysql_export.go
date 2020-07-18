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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/errors"
	loga "github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/admpub/nging/application/library/dbmanager/driver/mysql/utils"
	"github.com/mholt/archiver/v3"
)

// SQLTempDir sql文件缓存目录获取函数(用于导入导出SQL)
var SQLTempDir = os.TempDir

func TempDir(op string) string {
	dir := filepath.Join(SQLTempDir(), `dbmanager/cache`, op)
	err := os.MkdirAll(dir, os.ModePerm)
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

func (m *mySQL) bgExecManage(op utils.OP) error {
	var err error
	if m.IsPost() {
		data := m.Data()
		keys := m.FormValues(`key`)
		for _, key := range keys {
			err = utils.Cancel(op, key)
			if err != nil {
				return m.JSON(data.SetError(err))
			}
		}
		data.SetInfo(m.T(`操作成功`))
		m.ok(m.T(`操作成功`))
		return m.returnTo(m.GenURL(op.String()) + `&process=1`)
	}
	m.Set(`op`, op)
	m.Set(`list`, utils.ListBy(op))
	var title string
	if op == utils.OpExport {
		title = m.T(`导出SQL`)
	} else {
		title = m.T(`导入SQL`)
	}
	m.Set(`title`, title)
	m.Set(`cacheDir`, handler.URLFor(`/download/file?path=dbmanager/cache/`+op.String()))
	return m.Render(`db/mysql/process_store`, m.checkErr(err))
}

func (m *mySQL) Export() error {
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
				views = append(tables, strings.Split(views[0], `,`)...)
			}
			if len(views) > 0 {
				tables = append(tables, views...)
			}
		}
		output := m.Form(`output`)
		types := m.FormValues(`type`)
		cacheKey := com.Md5(com.Dump([]interface{}{tables, output, types}, false))

		var exports utils.Exec
		if old, exists := utils.Backgrounds.Load(utils.OpExport); exists {
			exports = old.(utils.Exec)
		} else {
			exports = utils.Exec{}
		}
		if exports.Exists(cacheKey) {
			return errors.New(m.T(`任务正在后台处理中，请稍候...`))
		}

		var (
			structWriter, dataWriter interface{}
			sqlFiles                 []string
			async                    bool
			bgExec                   = utils.NewGBExec(nil, echo.H{
				`database`: m.dbName,
				`tables`:   tables,
				`output`:   output,
				`types`:    types,
			})
			fileInfos = bgExec.Procs
		)
		exports[cacheKey] = bgExec
		nowTime := com.String(time.Now().Unix())
		saveDir := TempDir(`export`)
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
			if com.InSlice(`struct`, types) {
				structFile := filepath.Join(saveDir, m.dbName+`-struct-`+nowTime+`.sql`)
				sqlFiles = append(sqlFiles, structFile)
				structWriter = structFile
				fi := &utils.FileInfo{
					Start: time.Now(),
					Path:  structFile,
				}
				*fileInfos = append(*fileInfos, fi)
			}
			if com.InSlice(`data`, types) {
				dataFile := filepath.Join(saveDir, m.dbName+`-data-`+nowTime+`.sql`)
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
		worker := func(ctx context.Context, cfg driver.DbAuth) error {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf(`RECOVER: %v`, r)
				}
			}()
			err = utils.Export(ctx, &cfg, tables, structWriter, dataWriter, true)
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
				zipFile := filepath.Join(saveDir, cfg.Db+"-sql-"+nowTime+".zip")
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
				*fileInfos = append(*fileInfos, fi)
				ioutil.WriteFile(zipFile+`.txt`, com.Str2bytes(com.Dump(fileInfos, false)), os.ModePerm)
				exports.Cancel(cacheKey)
			}
			return nil
		}
		if !async {
			done := make(chan struct{})
			clientGone := m.Response().StdResponseWriter().(http.CloseNotifier).CloseNotify()
			go func() {
				for {
					select {
					case <-clientGone:
						bgExec.Cancel()()
						return
					case <-done:
						return
					}
				}
			}()
			err = worker(bgExec.Context(), cfg)
			done <- struct{}{}
			return err
		}
		data := m.Data()
		data.SetInfo(m.T(`任务已经在后台成功启动`))
		data.SetURL(handler.URLFor(`/download/file?path=dbmanager/cache/export`))
		go worker(m, cfg)
		exports.Add(utils.OpExport, cacheKey, bgExec)
		return m.JSON(data)
	}
	return m.Redirect(m.GenURL(`listTable`, m.dbName))
}
