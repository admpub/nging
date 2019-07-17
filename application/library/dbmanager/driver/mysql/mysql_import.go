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
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/admpub/nging/application/library/dbmanager/driver"
)

// Import 导入SQL文件
func Import(cfg *driver.DbAuth, sqlFiles []string, outWriter io.Writer, asyncs ...bool) error {
	if len(sqlFiles) == 0 {
		return nil
	}
	log.Println(`Starting import:`, sqlFiles)
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
		"SET FOREIGN_KEY_CHECKS=0;SET UNIQUE_CHECKS=0;source %s;SET FOREIGN_KEY_CHECKS=1;SET UNIQUE_CHECKS=1;",
	}
	for _, sqlFile := range sqlFiles {
		if len(sqlFile) == 0 {
			continue
		}
		lastIndex := len(args) - 1
		args[lastIndex] = fmt.Sprintf(args[lastIndex], sqlFile)
		cmd := exec.Command("mysql", args...)
		if outWriter != nil {
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				return fmt.Errorf(`Failed to import: %v`, err)
			}
			_, err = io.Copy(outWriter, stdout)
			if err != nil {
				stdout.Close()
				return fmt.Errorf(`Failed to import: %v`, err)
			}
			stdout.Close()
		}
		if err := cmd.Start(); err != nil {
			return fmt.Errorf(`Failed to import: %v`, err)
		}
		if !async {
			if err := cmd.Wait(); err != nil {
				return err
			}
		}
	}
	return nil
}
