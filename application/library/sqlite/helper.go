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
package sqlite

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/nging/application/library/config"
	"github.com/webx-top/com"
)

var (
	sqlComment   = regexp.MustCompile("(?is) COMMENT '[^']*'")
	sqlPK        = regexp.MustCompile("(?is),[\\s]*PRIMARY KEY \\(([^)]+)\\)([,]?)")
	sqlEngine    = regexp.MustCompile("(?is) ENGINE=[^ ]+ [^;]*;")
	sqlEnum      = regexp.MustCompile("(?is) (enum|set)\\(([^)]+)\\) ")
	sqlUnsigned  = regexp.MustCompile("(?is) unsigned ")
	sqlTableName = regexp.MustCompile("CREATE TABLE [^`]*`([^`]+)` \\(")
	sqlUnique    = regexp.MustCompile("(?is),[\\s]*UNIQUE KEY `([\\w]+)` \\(([^)]+)\\)([,]?)")
	sqlIndex     = regexp.MustCompile("(?is),[\\s]*KEY `([\\w]+)` \\(([^)]+)\\)([,]?)")
	sqlInteger   = regexp.MustCompile("(?is) (smallint|tinyint|bigint|int)\\([0-9]+\\) ")
	sqlCharset   = regexp.MustCompile("(?is) character set [^ ]* ")
	sqlOnUpdate  = regexp.MustCompile("(?is) on update [^,]*")
	sqlAutoIncr  = regexp.MustCompile("(?is) (unsigned )?(NOT NULL )?AUTO_INCREMENT")

	alterSQLComment           = regexp.MustCompile("(?is),[\\r\\n\\s]*COMMENT '[^']*'")
	alterSQLColumnB           = regexp.MustCompile("(?is) (AFTER|BEFORE) `[^`]+`")
	alterSQLColumnA           = regexp.MustCompile("(?is) (FIRST|LAST)$")
	alterSQLTableName         = regexp.MustCompile("ALTER TABLE `([^`]+)`")
	alterSQLRenameTableName   = regexp.MustCompile("ALTER TABLE `([^`]+)` RENAME TO ")
	alterSQLOperate           = regexp.MustCompile("(?is)[\\s]+(DROP|CHANGE|ADD) (INDEX )?([^,;]+)[,;]")
	alterSQLFieldChange       = regexp.MustCompile("(?is)[\\s]*`([^ ]+)` `([^ ]+)` ([^ ]+)") //旧字段名 新字段名 新字段数据类型
	alterSQLFieldAdd          = regexp.MustCompile("(?is)[\\s]*`([^ ]+)` ([^ ]+)")           //新字段名 新字段数据类型
	alterSQLFieldNULL         = regexp.MustCompile("(?is) (NOT )?NULL")
	alterSQLFieldDefaultValue = regexp.MustCompile("(?is) DEFAULT ([^ ]+)")
	alterSQLFieldUnsigned     = regexp.MustCompile("(?is) unsigned")
	alterSQLFieldAutoIncr     = regexp.MustCompile("(?is) AUTO_INCREMENT")
	alterSQLFieldUnique       = regexp.MustCompile("(?is) UNIQUE")
	alterSQLFieldCollate      = regexp.MustCompile("(?is) COLLATE [']?([^' ]+)[']?")
	sqlDDLParseSingle         = regexp.MustCompile("`([^`]+)` ([^,]+)")
	sqlDDLSeperator           = regexp.MustCompile(",[\\r\\n\\s]*`")
)

func execIntall(sqlStr string) error {
	sqls, err := covertCreateTableSQL(sqlStr)
	if err != nil {
		return err
	}
	for _, sql := range sqls {
		err = config.ExecMySQL(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

//CREATE TABLE `db_sync` (`id` integer PRIMARY KEY NOT NULL ,`dsn_source` varchar(255) NOT NULL,`dsn_destination` varchar(255) NOT NULL,`tables` text NOT NULL,`skip_tables` text NOT NULL,`alter_ignore` text NOT NULL,`drop` integer NOT NULL DEFAULT '0',`mail_to` varchar(200) NOT NULL DEFAULT '',`created` integer NOT NULL,`updated` integer NOT NULL DEFAULT '0')
func createTableSQL(table string) string {
	sqlStr := `select sql from SQLite_Master where tbl_name = '` + table + `' and type='table'`
	rows := []map[string]string{}
	_, err := config.QueryTo(sqlStr, &rows)
	if err != nil {
		log.Println(err.Error(), `->SQL:`, sqlStr)
	}
	if len(rows) > 0 {
		return rows[0]["sql"]
	}
	return ``
}

func covertCreateTableSQL(sqlStr string) ([]string, error) {
	matches := sqlTableName.FindStringSubmatch(sqlStr)
	if matches == nil {
		return nil, errors.New(`Can not find table name`)
	}
	tableName := matches[1]
	sqlStr = mySQLField2SQLite(sqlStr)
	var sqls []string
	matches = sqlPK.FindStringSubmatch(sqlStr)
	if len(matches) > 1 {
		sqlStr = sqlPK.ReplaceAllString(sqlStr, `$2`)
		items := strings.Split(matches[1], `,`)
		for _, item := range items {
			item = strings.Trim(item, "`")
			sqlPKCol := regexp.MustCompile("(?is)(`" + item + "`) [^ ]+ ((NOT )?NULL )?AUTO_INCREMENT")
			sqlStr = sqlPKCol.ReplaceAllString(sqlStr, `$1 integer PRIMARY KEY $4`)
		}
	}
	sqlStr = replaceEnum(sqlStr)
	matches2 := sqlUnique.FindAllStringSubmatch(sqlStr, -1)
	uniqueIndexes := []map[string]string{}
	for matches2 != nil {
		for _, matches := range matches2 {
			sqlStr = sqlUnique.ReplaceAllString(sqlStr, `$3`)
			uniqueIndexes = append(uniqueIndexes, map[string]string{
				`name`:    matches[1],
				`table`:   tableName,
				`columns`: matches[2],
			})
		}
		matches2 = sqlUnique.FindAllStringSubmatch(sqlStr, -1)
	}
	matches2 = sqlIndex.FindAllStringSubmatch(sqlStr, -1)
	indexes := []map[string]string{}
	for matches2 != nil {
		for _, matches := range matches2 {
			sqlStr = sqlIndex.ReplaceAllString(sqlStr, `$3`)
			indexes = append(indexes, map[string]string{
				`name`:    matches[1],
				`table`:   tableName,
				`columns`: matches[2],
			})
		}
		matches2 = sqlIndex.FindAllStringSubmatch(sqlStr, -1)
	}
	sqls = append(sqls, sqlStr)
	for _, v := range indexes {
		sql := fmt.Sprintf("CREATE INDEX `IDX_%[2]s_%[1]s` ON `%[2]s`(%[3]s);", v["name"], v["table"], v["columns"])
		sqls = append(sqls, sql)
	}
	for _, v := range uniqueIndexes {
		sql := fmt.Sprintf("CREATE UNIQUE INDEX `UNQ_%[2]s_%[1]s` ON `%[2]s`(%[3]s);", v["name"], v["table"], v["columns"])
		sqls = append(sqls, sql)
	}
	return sqls, nil
}

func replaceEnum(sqlStr string) string {
	return sqlEnum.ReplaceAllStringFunc(sqlStr, func(s string) string {
		match := sqlEnum.FindStringSubmatch(s)
		items := strings.Split(match[2], `,`)
		var maxSize int
		for _, item := range items {
			size := len(item)
			if size > 0 {
				switch item[0] {
				case '"', '\'':
					size -= 2
				default:
				}
			}
			if size > maxSize {
				maxSize = size
			}
		}
		if maxSize < 1 {
			maxSize = 1
		}
		return ` char(` + strconv.Itoa(maxSize) + `) `
	})
}

func foreignKeysState() string {
	sqlStr := `PRAGMA foreign_keys`
	rows := []map[string]string{}
	_, err := config.QueryTo(sqlStr, &rows)
	if err != nil {
		log.Println(err.Error(), `->SQL:`, sqlStr)
	}
	if len(rows) > 0 {
		return rows[0]["foreign_keys"]
	}
	return `0`
}

func indexSQL(table string) []map[string]string {
	sqlStr := `select sql from SQLite_Master where tbl_name = '` + table + `' and type='index'`
	rows := []map[string]string{}
	_, err := config.QueryTo(sqlStr, &rows)
	if err != nil {
		log.Println(err.Error(), `->SQL:`, sqlStr)
	}
	return rows
}

func mySQLField2SQLite(sqlStr string) string {
	sqlStr = sqlComment.ReplaceAllString(sqlStr, ``)
	sqlStr = sqlEngine.ReplaceAllString(sqlStr, `;`)

	sqlStr = sqlInteger.ReplaceAllString(sqlStr, ` integer `)
	sqlStr = sqlCharset.ReplaceAllString(sqlStr, ` `)
	sqlStr = sqlOnUpdate.ReplaceAllString(sqlStr, ``)

	sqlStr = sqlUnsigned.ReplaceAllString(sqlStr, ` `)

	sqlStr = alterSQLFieldCollate.ReplaceAllStringFunc(sqlStr, func(k string) string {
		if strings.HasSuffix(k, `_ci`) {
			return ` COLLATE NOCASE`
		}
		if strings.HasSuffix(k, `_bin`) {
			return ` COLLATE BINARY`
		}
		return ` `
	})
	return sqlStr
}

//TODO：development
//ALTER TABLE `task_log`
//DROP `task_id`,
//CHANGE `elapsed` `elapsed` int(12) NOT NULL DEFAULT '0' COMMENT '消耗时间(毫秒)' AFTER `status`,
//ADD `test` int(11) NOT NULL DEFAULT '0' COMMENT 'test',
//COMMENT='任务日志';
//
//ALTER TABLE `task_log`
//ADD INDEX `status` (`status`),
//DROP INDEX `idx_task_id`;
//
//ALTER TABLE `task_log`
//CHANGE `id` `id` int(11) unsigned NOT NULL FIRST,
//CHANGE `created` `created` int(11) NOT NULL COMMENT '创建时间' AUTO_INCREMENT UNIQUE AFTER `elapsed`;
func execAlter(sqlStr string) error {
	if alterSQLRenameTableName.MatchString(sqlStr) {
		return config.ExecMySQL(sqlStr)
	}
	matches := alterSQLTableName.FindStringSubmatch(sqlStr)
	if matches == nil {
		return errors.New(`Can not find table name`)
	}
	tableName := matches[1]
	ddlString := createTableSQL(tableName)
	ddlString = strings.TrimSpace(ddlString)
	if len(ddlString) < 1 {
		return nil
	}
	position := strings.Index(ddlString, `(`)
	if position < 1 {
		return nil
	}

	var ddlFieldsDef string
	sqlFieldsToInsert := []string{}

	fieldsDef := ddlString[position+1:]
	fieldsDef = strings.TrimRight(fieldsDef, `)`)
	fields := map[string]string{}
	fieldk := []string{}
	for _, fieldDef := range sqlDDLSeperator.Split(fieldsDef, -1) {
		match := sqlDDLParseSingle.FindStringSubmatch("`" + fieldDef)
		if len(match) < 3 {
			continue
		}
		fields[match[1]] = strings.TrimSpace(match[2])
		fieldk = append(fieldk, match[1])
		sqlFieldsToInsert = append(sqlFieldsToInsert, match[1])
	}

	matches2 := alterSQLOperate.FindAllStringSubmatch(sqlStr, -1)
	if matches2 != nil {
		for _, match := range matches2 {
			isIndex := len(match[2]) > 0
			match[3] = strings.TrimSpace(match[3])
			if isIndex {
				//TODO
				continue
			}
			switch strings.ToUpper(match[1]) {
			case `ADD`:
				findItems := alterSQLFieldAdd.FindStringSubmatch(match[3])
				if len(findItems) > 0 {
					newFieldName := strings.TrimSpace(findItems[1])
					newFieldType := strings.TrimSpace(findItems[2])
					newFieldType = mySQLField2SQLite(newFieldType)
					fields[newFieldName] = newFieldType
					fieldk = append(fieldk, newFieldName)
				}
			case `CHANGE`:
				findItems := alterSQLFieldChange.FindStringSubmatch(match[3])
				if len(findItems) > 0 {
					oldFieldName := strings.TrimSpace(findItems[1])
					newFieldName := strings.TrimSpace(findItems[2])
					newFieldType := strings.TrimSpace(findItems[3])
					newFieldType = mySQLField2SQLite(newFieldType)
					if oldFieldName == newFieldName {
						fields[oldFieldName] = newFieldType
					} else {
						fields[newFieldName] = newFieldType
						fieldk = append(fieldk, newFieldName)

						fieldName := oldFieldName
						_, exists := fields[fieldName]
						if exists {
							delete(fields, fieldName)
							com.SliceRemoveCallback(len(fieldk), func(i int) func(bool) error {
								if fieldk[i] != fieldName {
									return nil
								}
								return func(inside bool) error {
									if inside {
										fieldk = append(fieldk[0:i], fieldk[i+1:]...)
									} else {
										fieldk = fieldk[0:i]
									}
									return nil
								}
							})
							com.SliceRemoveCallback(len(sqlFieldsToInsert), func(i int) func(bool) error {
								if sqlFieldsToInsert[i] != fieldName {
									return nil
								}
								return func(inside bool) error {
									if inside {
										sqlFieldsToInsert = append(sqlFieldsToInsert[0:i], sqlFieldsToInsert[i+1:]...)
									} else {
										sqlFieldsToInsert = sqlFieldsToInsert[0:i]
									}
									return nil
								}
							})
						}
					}
				}
			case `DROP`:
				fieldName := strings.Trim(match[3], "`")
				_, exists := fields[fieldName]
				if exists {
					delete(fields, fieldName)
					com.SliceRemoveCallback(len(fieldk), func(i int) func(bool) error {
						if fieldk[i] != fieldName {
							return nil
						}
						return func(inside bool) error {
							if inside {
								fieldk = append(fieldk[0:i], fieldk[i+1:]...)
							} else {
								fieldk = fieldk[0:i]
							}
							return nil
						}
					})
					com.SliceRemoveCallback(len(sqlFieldsToInsert), func(i int) func(bool) error {
						if sqlFieldsToInsert[i] != fieldName {
							return nil
						}
						return func(inside bool) error {
							if inside {
								sqlFieldsToInsert = append(sqlFieldsToInsert[0:i], sqlFieldsToInsert[i+1:]...)
							} else {
								sqlFieldsToInsert = sqlFieldsToInsert[0:i]
							}
							return nil
						}
					})
				}
			}
		}
	}
	tempTable := "_" + tableName + "_old_" + time.Now().Local().Format("20060102_150405")
	newTableFields := "`" + strings.Join(sqlFieldsToInsert, "`,`") + "`"
	queryies := []string{
		"SAVEPOINT alter_column_" + tableName,
		"PRAGMA foreign_keys = 0",
		"PRAGMA triggers = NO",
		"ALTER TABLE `" + tableName + "` RENAME TO `" + tempTable + "`",
		//"CREATE TABLE `" + tempTable + "` AS SELECT * FROM `" + tableName + "`",
		//"DROP TABLE `" + tableName + "`",
		"CREATE TABLE `" + tableName + "` (" + strings.Trim(ddlFieldsDef, " \n\r\t,") + ")",
		"INSERT INTO `" + tableName + "` SELECT " + newTableFields + " FROM `" + tempTable + "`",
		"DROP TABLE `" + tempTable + "`",
	}
	// Create indexes for the new table
	indexes := indexSQL(tableName)
	for _, index := range indexes {
		queryies = append(queryies, index["sql"])
	}

	/// @todo add views
	queryies = append(queryies, "PRAGMA triggers = YES")
	queryies = append(queryies, "PRAGMA foreign_keys = "+foreignKeysState())
	queryies = append(queryies, "RELEASE alter_column_"+tableName)

	return config.ExecMySQL(strings.Join(queryies, `;`))
}

func ExecSQL(sqlStr string) error {
	if strings.HasPrefix(sqlStr, `SET `) {
		return nil
	}
	if strings.HasPrefix(sqlStr, `CREATE TABLE `) {
		return execIntall(sqlStr)
	}
	// if strings.HasPrefix(sqlStr, `ALTER TABLE `) {
	// 	return execAlter(sqlStr)
	// }
	return config.ExecMySQL(sqlStr)
}
