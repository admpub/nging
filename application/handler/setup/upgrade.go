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

package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/webx-top/com"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
)

// Upgrade 通过自动执行配置文件目录下的upgrade.sql来进行升级
func Upgrade() error {
	if err := config.OnceUpgradeDB(); err != nil {
		return err
	}
	sqlDir := filepath.Dir(config.DefaultCLIConfig.Conf)
	sqlFile := filepath.Join(sqlDir, `upgrade.sql`)
	if !com.FileExists(sqlFile) {
		return os.ErrNotExist
	}
	log.Info(color.GreenString(`[upgrader]`), `Execute SQL file: `, sqlFile)
	//创建数据表
	installer, ok := config.DBInstallers[config.DefaultConfig.DB.Type]
	if !ok {
		return fmt.Errorf(`不支持安装到%s`, config.DefaultConfig.DB.Type)
	}
	err := com.SeekFileLines(sqlFile, common.SQLLineParser(installer))
	if err != nil {
		return err
	}
	bakFile := filepath.Join(sqlDir, `upgraded.sql`)
	if com.FileExists(bakFile) {
		bakFile += `.` + time.Now().Format(`20060102150405`)
	}
	err = os.Rename(sqlFile, bakFile)
	return err
}
