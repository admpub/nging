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
	loga "github.com/admpub/log"
	writerPkg "github.com/admpub/nging/application/library/cron/writer"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/admpub/nging/application/library/notice"
	"github.com/webx-top/com"
)

// Import 导入SQL文件
func Import(ctx context.Context, noticer notice.Noticer, cfg *driver.DbAuth, cacheDir string, files []string, asyncs ...bool) error {
	if len(files) == 0 {
		return nil
	}
	if noticer == nil {
		noticer = notice.DefaultNoticer
	}
	names := make([]string, len(files))
	for i, file := range files {
		names[i] = filepath.Base(file)
	}
	noticer(`开始导入: `+strings.Join(names, ", "), 1)
	var (
		port, host string
		async      = true
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
		``,
	}
	sqls := `SET FOREIGN_KEY_CHECKS=0;SET UNIQUE_CHECKS=0;source %s;SET FOREIGN_KEY_CHECKS=1;SET UNIQUE_CHECKS=1;`
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
			err := archiver.Unarchive(sqlFile, dir)
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
	sqlFiles = append(sqlFiles, dataFiles...)
	rec := writerPkg.New(1000)
	for _, sqlFile := range sqlFiles {
		if len(sqlFile) == 0 {
			continue
		}
		sqlFile = filepath.ToSlash(sqlFile)
		lastIndex := len(args) - 1
		args[lastIndex] = fmt.Sprintf(sqls, sqlFile)
		//log.Println(`mysql`, strings.Join(args, ` `))
		cmd := exec.CommandContext(ctx, "mysql", args...)
		cmd.Stderr = rec
		if err := cmd.Start(); err != nil {
			return fmt.Errorf(`Failed to import: %v`, err)
		}
		if !async { //非异步，需阻塞
			if err := cmd.Wait(); err != nil {
				return errors.New(err.Error() + `: ` + rec.String())
			}
		}
	}
	noticer(`结束导入`, 1)
	return nil
}
