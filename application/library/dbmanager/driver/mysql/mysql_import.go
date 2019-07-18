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
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/archiver"
	"github.com/admpub/errors"
	loga "github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/cron"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/admpub/nging/application/library/notice"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// Import 导入SQL文件
func Import(cfg *driver.DbAuth, files []string, asyncs ...bool) error {
	if len(files) == 0 {
		return nil
	}
	log.Println(`Starting import:`, files)
	var (
		port, host string
		async      bool
	)
	if len(asyncs) > 0 {
		async = asyncs[0]
	}
	if p := strings.LastIndex(cfg.Host, `:`); p > 0 {
		host = cfg.Host[0:p]
		port = cfg.Host[p+1:]
	} else {
		host = cfg.Host
	}
	if len(port) == 0 {
		port = `3306`
	}
	args := []string{
		"-h" + host,
		"-P" + port,
		"-u" + cfg.Username,
		"-p" + cfg.Password,
		cfg.Db,
		"-e",
		`SET FOREIGN_KEY_CHECKS=0;SET UNIQUE_CHECKS=0;source %s;SET FOREIGN_KEY_CHECKS=1;SET UNIQUE_CHECKS=1;`,
	}
	var delDirs []string
	sqlFiles := []string{}
	defer func() {
		for _, delDir := range delDirs {
			os.RemoveAll(delDir)
		}
		for _, sqlFile := range sqlFiles {
			if !com.FileExists(sqlFile) {
				continue
			}
			os.Remove(sqlFile)
		}
	}()
	nowTime := com.String(time.Now().Unix())
	for index, sqlFile := range files {
		switch strings.ToLower(filepath.Ext(sqlFile)) {
		case `.sql`:
			sqlFiles = append(sqlFiles, sqlFile)
		case `.zip`:
			dir := filepath.Join(os.TempDir(), fmt.Sprintf("upload-"+nowTime+"-%d", index))
			err := archiver.Zip.Open(sqlFile, dir)
			if err != nil {
				loga.Error(err)
				continue
			}
			delDirs = append(delDirs, dir)
			err = os.Remove(sqlFile)
			if err != nil {
				loga.Error(err)
			}
			ifiles := []string{}
			err = filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return err
				}
				if strings.ToLower(filepath.Ext(fpath)) != `.sql` {
					return nil
				}
				if strings.Contains(info.Name(), `struct`) {
					sqlFiles = append(sqlFiles, fpath)
					return nil
				}
				ifiles = append(ifiles, fpath)
				return nil
			})
			sqlFiles = append(sqlFiles, ifiles...)
		}
	}
	rec := cron.NewCmdRec(1000)
	for _, sqlFile := range sqlFiles {
		if len(sqlFile) == 0 {
			continue
		}
		sqlFile = filepath.ToSlash(sqlFile)
		lastIndex := len(args) - 1
		args[lastIndex] = fmt.Sprintf(args[lastIndex], sqlFile)
		log.Println(`mysql`, strings.Join(args, ` `))
		cmd := exec.Command("mysql", args...)
		cmd.Stderr = rec
		if err := cmd.Start(); err != nil {
			return fmt.Errorf(`Failed to import: %v`, err)
		}
		if !async {
			if err := cmd.Wait(); err != nil {
				return errors.New(err.Error() + `: ` + rec.String())
			}
		}
	}
	return nil
}

func responseDropzone(err error, ctx echo.Context) error {
	if err != nil {
		user := handler.User(ctx)
		if user != nil {
			notice.OpenMessage(user.Username, `upload`)
			notice.Send(user.Username, notice.NewMessageWithValue(`upload`, ctx.T(`文件上传出错`), err.Error()))
		}
		return ctx.JSON(echo.H{`error`: err.Error()}, 500)
	}
	return ctx.String(`OK`)
}

func (m *mySQL) Import() error {
	var err error
	if m.IsPost() {
		if len(m.dbName) == 0 {
			m.fail(m.T(`请选择数据库`))
			return m.returnTo(m.GenURL(`listDb`))
		}
		async := m.Formx(`async`).Bool()
		var sqlFiles []string
		saveDir := os.TempDir()
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
		cfg := *m.DbAuth
		cfg.Db = m.dbName
		err = Import(&cfg, sqlFiles, async)
		return responseDropzone(err, m.Context)
	}

	return m.Render(`db/mysql/import`, m.checkErr(err))
}
