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
	defer rows.Close()
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
