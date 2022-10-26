/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

	"github.com/admpub/color"
	"github.com/webx-top/com"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"
)

var (
	upgradeSQLs = map[string]map[float64][]string{}
)

func RegisterUpgradeSQL(project string, minNgingDBVer float64, upgradeSQL string) {
	if _, ok := upgradeSQLs[project]; !ok {
		upgradeSQLs[project] = map[float64][]string{
			minNgingDBVer: {upgradeSQL},
		}
		return
	}
	if _, ok := upgradeSQLs[project][minNgingDBVer]; !ok {
		upgradeSQLs[project][minNgingDBVer] = []string{upgradeSQL}
		return
	}
	upgradeSQLs[project][minNgingDBVer] = append(upgradeSQLs[project][minNgingDBVer], upgradeSQL)
}

// Upgrade 通过自动执行配置文件目录下的upgrade.sql来进行升级
func Upgrade() error {
	if err := config.OnceUpgradeDB(); err != nil {
		return err
	}
	//创建数据表
	installer, ok := config.DBInstallers[config.FromFile().DB.Type]
	if !ok {
		return fmt.Errorf(`不支持安装到%s`, config.FromFile().DB.Type)
	}
	for _, upgrades := range upgradeSQLs {
		for versionNum, sqlContents := range upgrades {
			if versionNum <= config.Version.DBSchema {
				continue
			}
			for _, sqlContent := range sqlContents {
				err := common.ParseSQL(sqlContent, false, installer)
				if err != nil {
					return err
				}
			}
		}
	}
	sqlDir := config.FromCLI().Confd
	sqlFile := filepath.Join(sqlDir, `upgrade.sql`)
	if !com.FileExists(sqlFile) {
		return os.ErrNotExist
	}
	log.Info(color.GreenString(`[upgrader]`), `Execute SQL file: `, sqlFile)
	err := com.SeekFileLines(sqlFile, common.SQLLineParser(installer))
	if err != nil {
		return err
	}
	bakFile := filepath.Join(sqlDir, `upgraded.sql`)
	if com.FileExists(bakFile) {
		bakFile += `.` + time.Now().Format(`20060102150405`)
	}
	err = com.Rename(sqlFile, bakFile)
	return err
}
