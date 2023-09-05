package config

import (
	"fmt"
	stdLog "log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/admpub/mysql-schema-sync/sync"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

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
		Sync:           true,
		Drop:           true,
		SourceDSN:      ``,
		DestDSN:        ``,
		Tables:         ``,
		SkipTables:     ``,
		MailTo:         ``,
		MySQLOnlineDDL: false,
	}
	autoUpgradeDB := FromFile().Extend.GetStore(`upgradeDB`)
	if autoUpgradeDB.Bool(`mySQLOnlineDDL`) {
		syncConfig.MySQLOnlineDDL = true
	}
	upgrader, ok := DBUpgraders[FromFile().DB.Type]
	if !ok {
		stdLog.Panicf(`Does not support upgrading %s data table`, FromFile().DB.Type)
	}
	dbOperators, err := upgrader(schema, syncConfig, FromFile().DB)
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
