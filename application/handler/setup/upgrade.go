/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/nging/application/library/config"
	"github.com/webx-top/com"
)

// Upgrade 通过自动执行配置文件目录下的upgrade.sql来进行升级
func Upgrade() error {
	sqlDir := filepath.Dir(config.DefaultCLIConfig.Conf)
	sqlFile := filepath.Join(sqlDir, `upgrade.sql`)
	if !com.FileExists(sqlFile) {
		return os.ErrNotExist
	}
	//创建数据表
	var sqlStr string
	installer, ok := config.DBInstallers[config.DefaultConfig.DB.Type]
	if !ok {
		return fmt.Errorf(`不支持安装到%s`, config.DefaultConfig.DB.Type)
	}
	err := com.SeekFileLines(sqlFile, func(line string) error {
		if strings.HasPrefix(line, `--`) {
			return nil
		}
		line = strings.TrimSpace(line)
		sqlStr += line
		if strings.HasSuffix(line, `;`) && len(sqlStr) > 0 {
			defer func() {
				sqlStr = ``
			}()
			return installer(sqlStr)
		}
		return nil
	})
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
