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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	stdSync "sync"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/setup"
)

const (
	LockFileName = `installed.lock`
)

var _ = FixWd()

var (
	Installed             sql.NullBool
	installedSchemaVer    float64
	installedTime         time.Time
	defaultConfig         *Config
	defaultConfigMu       stdSync.RWMutex
	defaultCLIConfig      = NewCLIConfig()
	ErrUnknowDatabaseType = errors.New(`unkown database type`)
	onceUpgrade           stdSync.Once
	sqlCollection         = NewSQLCollection().RegisterInstall(`nging`, setup.InstallSQL)
)

func FromCLI() *CLIConfig {
	return defaultCLIConfig
}

func FromFile() *Config {
	defaultConfigMu.RLock()
	v := defaultConfig
	defaultConfigMu.RUnlock()
	return v
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
