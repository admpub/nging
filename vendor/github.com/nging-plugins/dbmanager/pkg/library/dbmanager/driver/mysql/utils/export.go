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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/errors"
	"github.com/webx-top/com"
	"github.com/webx-top/db"

	"github.com/admpub/nging/v4/application/library/common"
	writerPkg "github.com/admpub/nging/v4/application/library/cron/writer"
	"github.com/admpub/nging/v4/application/library/notice"

	"github.com/nging-plugins/dbmanager/pkg/library/dbmanager/driver"
)

/*
mysqldump 参数说明：
-d 			结构(--no-data:不导出任何数据，只导出数据库表结构)
-t 			数据(--no-create-info:只导出数据，而不添加CREATE TABLE 语句)
-n 			(--no-create-db:只导出数据，而不添加CREATE DATABASE 语句）
-R 			(--routines:导出存储过程以及自定义函数)
-E 			(--events:导出事件)
--triggers 	(默认导出触发器，使用--skip-triggers屏蔽导出)
-B 			(--databases:导出数据库列表，单个库时可省略）
--tables 	表列表（单个表时可省略）
*/

var (
	cleanRegExp = regexp.MustCompile(` AUTO_INCREMENT=[0-9]*\s*`)
)

// Export 导出SQL文件
func Export(ctx context.Context, noticer notice.Noticer,
	cfg *driver.DbAuth, tables []string, structWriter, dataWriter interface{},
	mysqlVersion string, hasGTID bool, resetAutoIncrements ...bool) error {
	if len(tables) == 0 {
		return errors.New(`No table selected for export`)
	}
	if noticer == nil {
		noticer = notice.DefaultNoticer
	}
	exportorPath, err := common.LookPath(Exportor, MySQLBinPaths...)
	if err != nil {
		return err
	}
	if len(exportorPath) == 0 {
		exportorPath = Exportor
	}
	noticer(`开始备份: `+strings.Join(tables, ","), 1)
	var (
		port, host         string
		resetAutoIncrement bool
	)
	if len(resetAutoIncrements) > 0 {
		resetAutoIncrement = resetAutoIncrements[0]
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
	if !com.InSlice(cfg.Charset, Charsets) {
		return errors.New(`字符集charset值无效`)
	}
	args := []string{
		"--default-character-set=" + cfg.Charset,
		"--single-transaction",
	}
	if com.VersionComparex(mysqlVersion, `8.0.0`, `>=`) {
		args = append(args, "--column-statistics=0") // 低版本不支持
	}
	if hasGTID {
		args = append(args, "--set-gtid-purged=OFF")
	}
	args = append(args, []string{
		"--no-autocommit",
		//"--ignore-table="+cfg.Db+".details",
		//"--ignore-table="+cfg.Db+".large_table2",
		"--opt",
		"-d", //加上此参数代表只导出表结构，不导出数据
		"-h" + host,
		"-P" + port,
		"-u" + cfg.Username,
		"-p" + cfg.Password,
		cfg.Db,
		//"--result-file=/root/backup.sql",
	}...)
	clean := func(w io.Writer, r io.ReadCloser) {
		r.Close()
		if c, y := w.(io.Closer); y {
			c.Close()
		}
	}
	var typeOptIndex int
	for index, value := range args {
		if value == `-d` {
			typeOptIndex = index
			break
		}
	}
	//args = append(args, `--tables`)
	args = append(args, tables...)
	rec := writerPkg.New(1000)
	for index, writer := range []interface{}{structWriter, dataWriter} {
		if writer == nil {
			continue
		}
		if index > 0 {
			args[typeOptIndex] = `-t` //导出数据
		}
		//log.Println(exportorPath, strings.Join(args, ` `))
		cmd := exec.CommandContext(ctx, exportorPath, args...)
		cmd.Stderr = rec
		var (
			w        io.Writer
			err      error
			onFinish func() error
		)
		switch v := writer.(type) {
		case io.Writer:
			w = v
		case string:
			dir := filepath.Dir(v)
			err = com.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return fmt.Errorf(`failed to backup: %v`, err)
			}
			w, err = os.Create(v)
			if err != nil {
				return fmt.Errorf(`failed to backup: %v`, err)
			}
			onFinish = func() error {
				if index > 0 {
					return nil
				}
				if resetAutoIncrement {
					return ResetAutoIncrement(v)
				}
				return nil
			}
		default:
			return errors.Wrapf(db.ErrUnsupported, `SQL Writer Error: %T`, v)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			clean(w, stdout)
			return fmt.Errorf(`failed to backup (cmd.StdoutPipe): %v`, err)
		}
		if err := cmd.Start(); err != nil {
			clean(w, stdout)
			return fmt.Errorf(`failed to backup (cmd.Start): %v`, err)
		}
		if _, err := io.Copy(w, stdout); err != nil {
			clean(w, stdout)
			return fmt.Errorf(`failed to backup (io.Copy): %v`, err)
		}
		if err := cmd.Wait(); err != nil {
			clean(w, stdout)
			return errors.New(err.Error() + ` (cmd.Wait): ` + rec.String())
		}
		clean(w, stdout)
		if onFinish != nil {
			if err = onFinish(); err != nil {
				return err
			}
		}
	}
	noticer(`结束备份`, 1)
	return nil
}

// ResetAutoIncrement 重置AUTO_INCREMENT值为0
func ResetAutoIncrement(sqlStructFile string) error {
	b, err := os.ReadFile(sqlStructFile)
	if err != nil {
		return err
	}
	b = cleanRegExp.ReplaceAll(b, []byte(` `))
	return os.WriteFile(sqlStructFile, b, os.ModePerm)
}
