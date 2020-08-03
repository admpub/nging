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

package config

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	stdLog "log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	stdSync "sync"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/admpub/mysql-schema-sync/sync"
	"github.com/admpub/nging/application/library/common"
)

var (
	Installed             sql.NullBool
	installedSchemaVer    float64
	installedTime         time.Time
	reload                bool
	DefaultConfig         *Config
	DefaultCLIConfig      = NewCLIConfig()
	OAuthUserSessionKey   = `oauthUser`
	ErrUnknowDatabaseType = errors.New(`unkown database type`)
	onceUpgrade           stdSync.Once
)

func SetInstalled(lockFile string) error {
	now := time.Now()
	err := ioutil.WriteFile(lockFile, []byte(now.Format(`2006-01-02 15:04:05`)+"\n"+fmt.Sprint(Version.DBSchema)), os.ModePerm)
	if err != nil {
		return err
	}
	installedTime = now
	installedSchemaVer = Version.DBSchema
	Installed.Valid = true
	Installed.Bool = true
	return nil
}

func IsInstalled() bool {
	if !Installed.Valid {
		lockFile := filepath.Join(echo.Wd(), `installed.lock`)
		if info, err := os.Stat(lockFile); err == nil && info.IsDir() == false {
			if b, e := ioutil.ReadFile(lockFile); e == nil {
				content := string(b)
				content = strings.TrimSpace(content)
				lines := strings.Split(content, "\n")
				switch len(lines) {
				case 2:
					installedSchemaVer, _ = strconv.ParseFloat(strings.TrimSpace(lines[1]), 64)
					fallthrough
				case 1:
					installedTime, _ = time.Parse(`2006-01-02 15:04:05`, strings.TrimSpace(lines[0]))
				}
			}
			Installed.Valid = true
			Installed.Bool = true
		}
	}
	return Installed.Bool
}

func OnceUpgradeDB() error {
	onceUpgrade.Do(UpgradeDB)
	return nil
}

func UpgradeDB() {
	if !Installed.Bool {
		return
	}
	if Version.DBSchema <= installedSchemaVer {
		return
	}
	if DefaultConfig == nil {
		return
	}
	var upgraded bool
	if DefaultConfig.DB.Type == `mysql` {
		executePreupgrade()
		autoUpgradeDatabase()
		upgraded = true
	} else {
		stdLog.Panicln(`数据库表结构需要升级！`)
	}
	if !upgraded {
		return
	}
	installedSchemaVer = Version.DBSchema
	ioutil.WriteFile(filepath.Join(echo.Wd(), `installed.lock`), []byte(installedTime.Format(`2006-01-02 15:04:05`)+"\n"+fmt.Sprint(Version.DBSchema)), os.ModePerm)
}

func GetSQLInstallFiles() ([]string, error) {
	confDIR := filepath.Dir(DefaultCLIConfig.Conf)
	sqlFile := confDIR + echo.FilePathSeparator + `install.sql`
	if !com.FileExists(sqlFile) {
		return nil, os.ErrNotExist
	}
	sqlFiles := []string{sqlFile}
	matches, err := filepath.Glob(confDIR + echo.FilePathSeparator + `install.*.sql`)
	if len(matches) > 0 {
		sqlFiles = append(sqlFiles, matches...)
	}
	return sqlFiles, err
}

func GetPreupgradeSQLFiles() []string {
	confDIR := filepath.Dir(DefaultCLIConfig.Conf)
	sqlFiles := []string{}
	matches, _ := filepath.Glob(confDIR + echo.FilePathSeparator + `preupgrade.*.sql`)
	if len(matches) > 0 {
		sqlFiles = append(sqlFiles, matches...)
	}
	return sqlFiles
}

func GetSQLInsertFiles() []string {
	confDIR := filepath.Dir(DefaultCLIConfig.Conf)
	sqlFile := confDIR + echo.FilePathSeparator + `insert.sql`
	sqlFiles := []string{}
	if com.FileExists(sqlFile) {
		sqlFiles = append(sqlFiles, sqlFile)
	}
	matches, _ := filepath.Glob(confDIR + echo.FilePathSeparator + `insert.*.sql`)
	if len(matches) > 0 {
		sqlFiles = append(sqlFiles, matches...)
	}
	return sqlFiles
}

//处理自动升级前要执行的sql
func executePreupgrade() {
	preupgradeSQLFiles := GetPreupgradeSQLFiles()
	if len(preupgradeSQLFiles) < 1 {
		return
	}
	installer, ok := DBInstallers[DefaultConfig.DB.Type]
	if !ok {
		stdLog.Panicf(`不支持安装到%s`, DefaultConfig.DB.Type)
	}
	for _, sqlFile := range preupgradeSQLFiles {
		//sqlFile = /your/path/preupgrade.3_0.nging.sql
		versionStr := strings.TrimPrefix(filepath.Base(sqlFile), `preupgrade.`)
		versionStr = strings.TrimSuffix(versionStr, `.sql`)
		versionStr = strings.ReplaceAll(strings.SplitN(versionStr, `.`, 2)[0], `_`, `.`)
		versionNum, err := strconv.ParseFloat(versionStr, 64)
		if err != nil {
			stdLog.Panicln(versionStr + `: ` + err.Error())
		}
		if versionNum <= installedSchemaVer {
			continue
		}
		log.Info(color.GreenString(`[preupgrade]`), `Execute SQL file: `, sqlFile)
		err = common.ParseSQL(sqlFile, true, installer)
		if err != nil {
			stdLog.Panicln(err.Error())
		}
	}
}

//自动升级数据表
func autoUpgradeDatabase() {
	sqlFiles, err := GetSQLInstallFiles()
	if err != nil {
		panic(`尝试自动升级数据库失败！数据库安装文件不存在：config/install.sql`)
	}
	var schema string
	for _, sqlFile := range sqlFiles {
		b, err := ioutil.ReadFile(sqlFile)
		if err != nil {
			panic(err)
		}
		schema += string(b)
	}
	r, err := sync.Sync(&sync.Config{
		Sync:       true,
		Drop:       true,
		SourceDSN:  ``,
		DestDSN:    DefaultConfig.DB.User + `:` + DefaultConfig.DB.Password + `@(` + DefaultConfig.DB.Host + `)/` + DefaultConfig.DB.Database,
		Tables:     ``,
		SkipTables: ``,
		MailTo:     ``,
	}, nil, sync.NewMySchemaData(schema, `source`))
	if err != nil {
		panic(`尝试自动升级数据库失败！同步表结构时出错：` + err.Error())
	}
	nowTime := time.Now().Format(`20060102150405`)
	//写日志
	result := r.Diff(false).String()
	logName := `upgrade_` + fmt.Sprint(installedSchemaVer) + `_` + fmt.Sprint(Version.DBSchema) + `_` + nowTime
	result = `<!doctype html><html><head><meta charset="utf-8"><title>` + logName + `</title></head><body>` + result + `</body></html>`
	confDIR := filepath.Dir(DefaultCLIConfig.Conf)
	ioutil.WriteFile(filepath.Join(confDIR, logName+`.log.html`), []byte(result), os.ModePerm)
}
