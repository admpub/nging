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
package sqlite

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/admpub/nging/application/library/config"
)

var (
	sqlComment   = regexp.MustCompile("(?is) COMMENT '[^']*'")
	sqlPK        = regexp.MustCompile("(?is),PRIMARY KEY \\(([^)]+)\\)([,]?)")
	sqlEngine    = regexp.MustCompile("(?is)\\) ENGINE=InnoDB [^;]*;")
	sqlEnum      = regexp.MustCompile("(?is) enum\\(([^)]+)\\) ")
	sqlUnsigned  = regexp.MustCompile("(?is) unsigned ")
	sqlTableName = regexp.MustCompile("CREATE TABLE [^`]*`([^`]+)` \\(")
	sqlUnique    = regexp.MustCompile("(?is),UNIQUE KEY `([\\w]+)` \\(([^)]+)\\)([,]?)")
	sqlIndex     = regexp.MustCompile("(?is),KEY `([\\w]+)` \\(([^)]+)\\)([,]?)")
)

func Exec(sqlStr string) error {
	if strings.HasPrefix(sqlStr, `SET `) {
		return nil
	}
	if strings.HasPrefix(sqlStr, `CREATE TABLE `) {
		matches := sqlTableName.FindStringSubmatch(sqlStr)
		if matches == nil {
			return errors.New(`Can not find table name`)
		}
		tableName := matches[1]
		sqlStr = sqlComment.ReplaceAllString(sqlStr, ``)
		sqlStr = sqlEngine.ReplaceAllString(sqlStr, `);`)
		matches = sqlPK.FindStringSubmatch(sqlStr)
		if len(matches) > 1 {
			sqlStr = sqlPK.ReplaceAllString(sqlStr, `$2`)
			items := strings.Split(matches[1], `,`)
			for _, item := range items {
				item = strings.Trim(item, "`")
				sqlPKCol := regexp.MustCompile("(?is)(`" + item + "`) [^ ]+ (unsigned )?(NOT NULL )?AUTO_INCREMENT")
				sqlStr = sqlPKCol.ReplaceAllString(sqlStr, `$1 integer PRIMARY KEY $3`)
			}
		}
		matches2 := sqlEnum.FindAllStringSubmatch(sqlStr, -1)
		if matches2 != nil {
			for _, matches := range matches2 {
				items := strings.Split(matches[1], `,`)
				var maxSize int
				for _, item := range items {
					size := len(item)
					if size > maxSize {
						maxSize = size
					}
				}
				if maxSize > 1 {
					maxSize -= 2
				}
				sqlStr = sqlEnum.ReplaceAllString(sqlStr, ` char(`+strconv.Itoa(maxSize)+`) `)
			}
		}

		matches2 = sqlUnique.FindAllStringSubmatch(sqlStr, -1)
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
		sqlStr = sqlUnsigned.ReplaceAllString(sqlStr, ``)
		err := config.ExecMySQL(sqlStr)
		if err != nil {
			return err
		}
		for _, v := range indexes {
			sql := fmt.Sprintf("CREATE INDEX `IDX_%[2]s_%[1]s` ON `%[2]s`(%[3]s)", v["name"], v["table"], v["columns"])
			err = config.ExecMySQL(sql)
			if err != nil {
				return err
			}
		}
		for _, v := range uniqueIndexes {
			sql := fmt.Sprintf("CREATE UNIQUE INDEX `UNQ_%[2]s_%[1]s` ON `%[2]s`(%[3]s)", v["name"], v["table"], v["columns"])
			err = config.ExecMySQL(sql)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return config.ExecMySQL(sqlStr)
}
