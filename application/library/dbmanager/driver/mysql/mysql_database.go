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
package mysql

import (
	"database/sql"

	"github.com/webx-top/com"
)

func (m *mySQL) createDatabase(dbName, collate string) *Result {
	r := &Result{}
	r.SQL = "CREATE DATABASE " + quoteCol(dbName)
	if len(collate) > 0 {
		r.SQL += " COLLATE " + quoteVal(collate)
	}
	return r.Exec(m.newParam())
}

func (m *mySQL) dropDatabase(dbName string) *Result {
	r := &Result{}
	r.SQL = "DROP DATABASE " + quoteCol(dbName)
	return r.Exec(m.newParam())
}

func (m *mySQL) alterDatabase(dbName string, collation string) *Result {
	r := &Result{}
	r.SQL = "ALTER DATABASE " + quoteCol(dbName)
	if reOnlyWord.MatchString(collation) {
		r.SQL += " COLLATE " + collation
	}
	return r.Exec(m.newParam())
}

func (m *mySQL) renameDatabase(newName, collate string) []*Result {
	rs := []*Result{}
	r := m.createDatabase(newName, collate)
	rs = append(rs, r)
	if r.err != nil {
		return rs
	}
	rGetTables := &Result{}
	rGetTables.start()
	tables, err := m.getTables()
	rGetTables.end()
	rGetTables.SQL = `SHOW TABLES`
	rGetTables.err = err
	if err != nil {
		rGetTables.ErrorString = err.Error()
	}
	rs = append(rs, rGetTables)
	if err != nil {
		return rs
	}
	var sql string

	newName = quoteCol(newName)
	for key, table := range tables {
		table = quoteCol(table)
		if key > 0 {
			sql += ", "
		}
		sql += table + " TO " + newName + "." + table
	}
	if len(sql) > 0 {
		rRename := &Result{}
		rRename.SQL = "RENAME TABLE " + sql
		rRename = rRename.Exec(m.newParam())
		err = rRename.err
		rs = append(rs, rRename)
	}
	if err == nil {
		rDrop := m.dropDatabase(m.dbName)
		rs = append(rs, rDrop)
	}
	return rs
}

// 获取数据库列表
func (m *mySQL) getDatabases() ([]string, error) {
	sqlStr := `SELECT SCHEMA_NAME FROM information_schema.SCHEMATA`
	if com.VersionCompare(m.getVersion(), `5`) < 0 {
		sqlStr = `SHOW DATABASES`
	}
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for rows.Next() {
		var v sql.NullString
		err := rows.Scan(&v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v.String)
	}
	return ret, nil
}
