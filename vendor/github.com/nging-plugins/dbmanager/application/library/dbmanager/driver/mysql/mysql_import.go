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
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/errors"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/background"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/admpub/nging/v5/application/library/respond"

	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver/mysql/utils"
)

func responseDropzone(err error, ctx echo.Context) error {
	if err != nil {
		if user := handler.User(ctx); user != nil {
			notice.OpenMessage(user.Username, `upload`)
			notice.Send(user.Username, notice.NewMessageWithValue(`upload`, ctx.T(`文件上传出错`), err.Error()))
		}
	}
	return respond.Dropzone(ctx, err, nil)
}

func (m *mySQL) importing() error {
	return m.bgExecManage(utils.OpImport)
}

func (m *mySQL) Import() error {
	process := m.Queryx(`process`).Bool()
	if process {
		return m.importing()
	}
	var err error
	if m.IsPost() {
		if len(m.dbName) == 0 {
			m.fail(m.T(`请选择数据库`))
			return m.returnTo(m.GenURL(`listDb`))
		}
		user := handler.User(m.Context)
		var username string
		if user != nil {
			username = user.Username
		}
		async := m.Formx(`async`, `true`).Bool()
		var sqlFiles []string
		saveDir := TempDir(utils.OpImport)
		err = m.SaveUploadedFiles(`file`, func(fdr *multipart.FileHeader) (string, error) {
			extension := filepath.Ext(fdr.Filename)
			switch strings.ToLower(extension) {
			case `.sql`:
			case `.zip`:
			default:
				return ``, errors.New(`只能上传扩展名为“.sql”和“.zip”的文件`)
			}
			sqlFile := filepath.Join(saveDir, fdr.Filename)
			sqlFiles = append(sqlFiles, sqlFile)
			return sqlFile, nil
		})
		if err != nil {
			return responseDropzone(err, m.Context)
		}
		bgExec := background.New(context.TODO(), echo.H{
			`database`: m.dbName,
			`sqlFiles`: sqlFiles,
			`async`:    async,
		})
		cacheKey := bgExec.Started.Format(`20060102150405`)
		imports, err := background.Register(m.Context, m.ImportAndOutputOpName(utils.OpImport), cacheKey, bgExec)
		if err != nil {
			return responseDropzone(err, m.Context)
		}
		noticer := notice.NewP(m.Context, `databaseImport`, username, bgExec.Context())
		noticer.Success(m.T(`文件上传成功`))
		coll, err := m.getCollation(m.dbName, nil)
		if err != nil {
			return responseDropzone(err, m.Context)
		}
		cfg := *m.DbAuth
		cfg.Db = m.dbName
		cfg.Charset = strings.SplitN(coll, `_`, 2)[0]
		importor := func(c context.Context, noticer *notice.NoticeAndProgress, cfg *driver.DbAuth, cacheDir string, files []string) (err error) {
			if len(files) == 0 {
				return
			}
			if utils.SupportedImport() { // 采用 mysql 命令导入
				err = utils.Import(c, noticer, cfg, cacheDir, files)
			} else {
				err = m.importDB(c, noticer, cfg, cacheDir, files)
			}
			return
		}
		if async {
			go func(cfg driver.DbAuth) {
				done := make(chan error)
				go func(cfg driver.DbAuth) {
					err := importor(bgExec.Context(), noticer, &cfg, TempDir(utils.OpImport), sqlFiles)
					if err != nil {
						noticer.Failure(m.T(`导入失败`) + `: ` + err.Error())
						noticer.Complete().Failure(m.T(`导入结束 :(`))
					} else {
						noticer.Complete().Success(m.T(`导入结束 :)`))
					}
					imports.Cancel(cacheKey)
					done <- err
					close(done)
				}(cfg)
				t := time.NewTicker(24 * time.Hour)
				defer t.Stop()
				for {
					select {
					case <-t.C:
						imports.Cancel(cacheKey)
						return
					case <-done:
						return
					}
				}
			}(cfg)
			noticer.Success(m.T(`正在后台导入，请稍候...`))
		} else {
			done := make(chan struct{})
			ctx := m.StdContext()
			go func() {
				defer imports.Cancel(cacheKey)
				for {
					select {
					case <-ctx.Done():
						return
					case <-done:
						return
					}
				}
			}()
			err = importor(bgExec.Context(), noticer, &cfg, TempDir(utils.OpImport), sqlFiles)
			if err != nil {
				noticer.Failure(m.T(`导入失败`) + `: ` + err.Error())
			}
			done <- struct{}{}
			close(done)
		}
		return responseDropzone(err, m.Context)
	}

	return m.Render(`db/mysql/import`, m.checkErr(err))
}
