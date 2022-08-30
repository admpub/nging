/*
   Nging is a toolbox for webmasters
   Copyright (C) 2019-present  Wenhui Shen <swh@admpub.com>

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

package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/archiver"
	"github.com/admpub/log"
	"github.com/fatih/color"
	"github.com/webx-top/com"

	"github.com/admpub/nging/v4/application/library/common"
	writerPkg "github.com/admpub/nging/v4/application/library/cron/writer"
	"github.com/admpub/nging/v4/application/library/notice"

	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver"
)

// Import 导入SQL文件
func Import(ctx context.Context, noticer *notice.NoticeAndProgress, cfg *driver.DbAuth, cacheDir string, files []string) error {
	if len(files) == 0 {
		return nil
	}
	importorPath, err := common.LookPath(Importor, MySQLBinPaths...)
	if err != nil {
		return err
	}
	if len(importorPath) == 0 {
		importorPath = Importor
	}
	names := make([]string, len(files))
	for i, file := range files {
		names[i] = filepath.Base(file)
	}
	noticer.Success(`开始导入: ` + strings.Join(names, ", "))
	var (
		port string
		host string
	)
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
		``,
	}
	if !com.InSlice(cfg.Charset, Charsets) {
		return errors.New(`字符集charset值无效`)
	}
	sqls := `SET NAMES ` + cfg.Charset + `;SET FOREIGN_KEY_CHECKS=0;SET UNIQUE_CHECKS=0;source %s;SET FOREIGN_KEY_CHECKS=1;SET UNIQUE_CHECKS=1;`
	var (
		delDirs  []string
		sqlFiles []string
		onFilish = func() {
			for _, delDir := range delDirs {
				os.RemoveAll(delDir)
			}
			for _, sqlFile := range sqlFiles {
				if !com.FileExists(sqlFile) {
					continue
				}
				os.Remove(sqlFile)
			}
		}
	)
	defer onFilish()
	nowTime := com.String(time.Now().Unix())
	dataFiles := []string{}
	for index, sqlFile := range files {
		switch strings.ToLower(filepath.Ext(sqlFile)) {
		case `.sql`:
			if strings.Contains(filepath.Base(sqlFile), `struct`) {
				sqlFiles = append(sqlFiles, sqlFile)
			} else {
				dataFiles = append(dataFiles, sqlFile)
			}
		case `.zip`:
			dir := filepath.Join(cacheDir, fmt.Sprintf("upload-"+nowTime+"-%d", index))
			err := com.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
			err = archiver.Unarchive(sqlFile, dir)
			if err != nil {
				log.Error(err)
				continue
			}
			delDirs = append(delDirs, dir)
			err = os.Remove(sqlFile)
			if err != nil {
				log.Error(err)
			}
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
				dataFiles = append(dataFiles, fpath)
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	sqlFiles = append(sqlFiles, dataFiles...)
	lastIndex := len(args) - 1
	noticer.Add(int64(len(sqlFiles)))
	for _, sqlFile := range sqlFiles {
		if len(sqlFile) == 0 {
			continue
		}
		sqlFile = filepath.ToSlash(sqlFile)
		args := args[:]
		args[lastIndex] = fmt.Sprintf(sqls, sqlFile)
		//log.Println(importorPath, strings.Join(args, ` `))
		//log.Debug(args[lastIndex])
		cmd := exec.CommandContext(ctx, importorPath, args...)
		rec := writerPkg.New(1000)
		cmd.Stderr = rec
		if err := cmd.Start(); err != nil {
			return fmt.Errorf(`failed to import: %v`, err)
		}
		if err := cmd.Wait(); err != nil {
			noticer.Done(1)
			noticer.Failure(`[FAILURE] ` + err.Error() + `: ` + rec.String() + `: ` + filepath.Base(sqlFile))
			log.Debug(color.RedString(`[FAILURE]`), ` `, err.Error(), `: `+rec.String()+`: `, args[lastIndex])
		} else {
			noticer.Done(1)
			noticer.Success(`[SUCCESS] ` + filepath.Base(sqlFile))
			log.Debug(color.GreenString(`[SUCCESS]`), ` `, args[lastIndex])
		}
	}
	return nil
}
