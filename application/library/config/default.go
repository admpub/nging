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

package config

import (
	"database/sql"
	"errors"
	"fmt"
	stdLog "log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	stdSync "sync"
	"time"

	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/admpub/mysql-schema-sync/sync"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/setup"
)

const (
	LockFileName = `installed.lock`
)

var (
	Installed             sql.NullBool
	installedSchemaVer    float64
	installedTime         time.Time
	defaultConfig         *Config
	defaultCLIConfig      = NewCLIConfig()
	ErrUnknowDatabaseType = errors.New(`unkown database type`)
	onceUpgrade           stdSync.Once
	sqlCollection         = NewSQLCollection().RegisterInstall(`nging`, setup.InstallSQL)
)

func FromCLI() *CLIConfig {
	return defaultCLIConfig
}

func FromFile() *Config {
	return defaultConfig
}

func FromDB(group ...string) echo.H {
	return echo.GetStoreByKeys(common.SettingName, group...)
}

func GetSQLCollection() *SQLCollection {
	return sqlCollection
}

func RegisterInstallSQL(project string, installSQL string) {
	sqlCollection.RegisterInstall(project, installSQL)
}

func RegisterInsertSQL(project string, insertSQL string) {
	sqlCollection.RegisterInsert(project, insertSQL)
}

func RegisterPreupgradeSQL(project string, version, preupgradeSQL string) {
	sqlCollection.RegisterPreupgrade(project, version, preupgradeSQL)
}

func GetInsertSQLs() map[string][]string {
	return sqlCollection.Insert
}

func GetInstallSQLs() map[string][]string {
	return sqlCollection.Install
}

func GetPreupgradeSQLs() map[string]map[string][]string {
	return sqlCollection.Preupgrade
}

func SetInstalled(lockFile string) error {
	now := time.Now()
	err := os.WriteFile(lockFile, []byte(now.Format(`2006-01-02 15:04:05`)+"\n"+fmt.Sprint(Version.DBSchema)), os.ModePerm)
	if err != nil {
		return err
	}
	installedTime = now
	installedSchemaVer = Version.DBSchema
	Installed.Valid = true
	Installed.Bool = true
	return nil
}

func InstalledLockFile() string {
	for _, lockFile := range []string{
		filepath.Join(FromCLI().ConfDir(), LockFileName),
		filepath.Join(echo.Wd(), LockFileName),
	} {
		if info, err := os.Stat(lockFile); err == nil && !info.IsDir() {
			return lockFile
		}
	}
	return ``
}

func IsInstalled() bool {
	if !Installed.Valid {
		lockFile := InstalledLockFile()
		if len(lockFile) > 0 {
			if b, e := os.ReadFile(lockFile); e == nil {
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
	if FromFile() == nil {
		return
	}
	log.Info(`Start to upgrade the database table`)
	if FromFile().DB.Type == `sqlite` {
		// 升级前自动备份当前已安装版本数据库
		log.Info(`Automatically backup the current database`)
		backupName := FromFile().DB.Database + `.` + strings.Replace(fmt.Sprintf(`%v`, installedSchemaVer), `.`, `_`, -1) + `.bak`
		if !com.FileExists(backupName) {
			err := com.Copy(FromFile().DB.Database, backupName)
			if err != nil {
				stdLog.Panicf(`An error occurred while backing up the database "%s" to "%s": %v`, FromFile().DB.Database, backupName, err.Error())
			} else {
				log.Infof(`Backup database "%s" to "%s"`, FromFile().DB.Database, backupName)
			}
		} else {
			log.Infof(`The database backup file "%s" already exists, skip this backup`, backupName)
		}
	}
	eventParams := param.Store{
		`installedSchemaVer`: installedSchemaVer,
		`currentSchemaVer`:   Version.DBSchema,
	}
	echo.PanicIf(echo.FireByNameWithMap(`nging.upgrade.db.before`, eventParams))
	executePreupgrade()
	autoUpgradeDatabase()
	echo.PanicIf(echo.FireByNameWithMap(`nging.upgrade.db.after`, eventParams))
	installedSchemaVer = Version.DBSchema
	err := os.WriteFile(filepath.Join(FromCLI().ConfDir(), LockFileName), []byte(installedTime.Format(`2006-01-02 15:04:05`)+"\n"+fmt.Sprint(Version.DBSchema)), os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	log.Info(`Database table upgrade completed`)
}

func GetSQLInstallFiles() ([]string, error) {
	confDIR := FromCLI().Confd
	sqlFile := filepath.Join(confDIR, `install.sql`)
	var sqlFiles []string
	if com.FileExists(sqlFile) {
		sqlFiles = append(sqlFiles, sqlFile)
	}
	matches, err := filepath.Glob(confDIR + echo.FilePathSeparator + `install.*.sql`)
	if len(matches) > 0 {
		sqlFiles = append(sqlFiles, matches...)
	}
	return sqlFiles, err
}

func GetPreupgradeSQLFiles() []string {
	confDIR := FromCLI().Confd
	sqlFiles := []string{}
	matches, _ := filepath.Glob(confDIR + echo.FilePathSeparator + `preupgrade.*.sql`)
	if len(matches) > 0 {
		sqlFiles = append(sqlFiles, matches...)
	}
	return sqlFiles
}

func GetSQLInsertFiles() []string {
	confDIR := FromCLI().Confd
	sqlFile := filepath.Join(confDIR, `insert.sql`)
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

// 处理自动升级前要执行的sql
func executePreupgrade() {
	preupgradeSQLFiles := GetPreupgradeSQLFiles()
	if len(preupgradeSQLFiles) == 0 && len(GetPreupgradeSQLs()) == 0 {
		return
	}
	installer, ok := DBInstallers[FromFile().DB.Type]
	if !ok {
		stdLog.Panicf(`Does not support installation to database: %s`, FromFile().DB.Type)
	}
	for _, sqlFile := range preupgradeSQLFiles {
		//sqlFile = /your/path/preupgrade.3_0.nging.sql //preupgrade.{versionStr}.{project}.sql
		versionStr := strings.TrimPrefix(filepath.Base(sqlFile), `preupgrade.`)
		versionStr = strings.TrimSuffix(versionStr, `.sql`) // {versionStr}.{project}
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
	for _, sqlVersionContents := range GetPreupgradeSQLs() {
		for versionStr, sqlContents := range sqlVersionContents {
			versionNum, err := strconv.ParseFloat(versionStr, 64)
			if err != nil {
				stdLog.Panicln(versionStr + `: ` + err.Error())
			}
			if versionNum <= installedSchemaVer {
				continue
			}
			for _, sqlContent := range sqlContents {
				log.Info(color.GreenString(`[preupgrade]`), `Execute SQL: `, sqlContent)
				err = common.ParseSQL(sqlContent, false, installer)
				if err != nil {
					stdLog.Panicln(err.Error())
				}
			}
		}
	}
}

// 自动升级数据表
func autoUpgradeDatabase() {
	sqlFiles, err := GetSQLInstallFiles()
	if err != nil && len(GetInstallSQLs()[`nging`]) == 0 {
		stdLog.Panicln(`Attempt to automatically upgrade the database failed! The database installation file does not exist: config/install.sql`)
	}
	var schema string
	for _, sqlFile := range sqlFiles {
		b, err := os.ReadFile(sqlFile)
		if err != nil {
			stdLog.Panicln(err)
		}
		schema += string(b)
	}
	for _, sqlContents := range GetInstallSQLs() {
		for _, sqlContent := range sqlContents {
			schema += sqlContent
		}
	}
	syncConfig := &sync.Config{
		Sync:       true,
		Drop:       true,
		SourceDSN:  ``,
		DestDSN:    ``,
		Tables:     ``,
		SkipTables: ``,
		MailTo:     ``,
	}
	upgrader, ok := DBUpgraders[FromFile().DB.Type]
	if !ok {
		stdLog.Panicf(`Does not support upgrading %s data table`, FromFile().DB.Type)
	}
	dbOperators, err := upgrader(schema, syncConfig, FromFile())
	if err != nil {
		stdLog.Panicln(err)
	}
	r, err := sync.Sync(syncConfig, nil, dbOperators.Source, dbOperators.Destination)
	if err != nil {
		stdLog.Panicln(`Attempt to automatically upgrade the database failed! Error while synchronizing table structure: ` + err.Error())
	}
	nowTime := time.Now().Format(`20060102150405`)
	//写日志
	result := r.Diff(false).String()
	logName := `upgrade_` + fmt.Sprint(installedSchemaVer) + `_` + fmt.Sprint(Version.DBSchema) + `_` + nowTime
	result = `<!doctype html><html><head><meta charset="utf-8"><title>` + logName + `</title></head><body>` + result + `</body></html>`
	confDIR := FromCLI().ConfDir()
	os.WriteFile(filepath.Join(confDIR, logName+`.log.html`), []byte(result), os.ModePerm)
}
