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
	"fmt"

	"github.com/webx-top/com"
)

func (m *mySQL) support(feature string) bool {
	switch feature {
	case "scheme", "sequence", "type", "view_trigger":
		return true
	default:
		if com.VersionCompare(m.getVersion(), "5.1") == 1 {
			switch feature {
			case "event", "partitioning":
				return true
			}
		}
		if com.VersionCompare(m.getVersion(), "5") == 1 {
			switch feature {
			case "routine", "trigger", "view":
				return true
			}
		}
		return false
	}
}

func (m *mySQL) showVariables() ([]map[string]string, error) {
	sqlStr := "SHOW VARIABLES"
	return m.kvVal(sqlStr)
}

func (m *mySQL) killProcess(processId int64) error {
	sqlStr := fmt.Sprintf("KILL %d", processId)
	_, err := m.newParam().SetCollection(sqlStr).Exec()
	return err
}

func (m *mySQL) processList() ([]*ProcessList, error) {
	r := []*ProcessList{}
	sqlStr := "SHOW FULL PROCESSLIST"
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return r, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return r, err
	}
	n := len(cols)
	for rows.Next() {
		v := &ProcessList{}
		err = safeScan(rows, n, &v.Id, &v.User, &v.Host, &v.Db, &v.Command, &v.Time, &v.State, &v.Info, &v.Progress)
		if err != nil {
			err = fmt.Errorf(`%v: %v`, sqlStr, err)
			break
		}
		r = append(r, v)
	}
	return r, err
}

func (m *mySQL) showStatus() ([]map[string]string, error) {
	sqlStr := "SHOW STATUS"
	return m.kvVal(sqlStr)
}

func (m *mySQL) getCharsets() (map[string]CharsetData, error) {
	sqlStr := `SHOW CHARSET`
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	n := len(cols)
	ret := map[string]CharsetData{}
	for rows.Next() {
		var v CharsetData
		err = safeScan(rows, n, &v.Charset, &v.Description, &v.DefaultCollation, &v.Maxlen)
		if err != nil {
			return nil, fmt.Errorf(`%v: %v`, sqlStr, err)
		}
		ret[v.Charset.String] = v
	}
	return ret, nil
}

// 获取支持的字符集
func (m *mySQL) getCollations() (*Collations, error) {
	sqlStr := `SHOW COLLATION`
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf(`%v: %v`, sqlStr, err)
	}
	ret := NewCollations()
	for rows.Next() {
		var v Collation
		err = safeScan(rows, len(cols), &v.Collation, &v.Charset, &v.Id, &v.Default, &v.Compiled, &v.Sortlen, &v.PadAttribute)
		if err != nil {
			return nil, fmt.Errorf(`%v: %v`, sqlStr, err)
		}
		coll, ok := ret.Collations[v.Charset.String]
		if !ok {
			coll = []Collation{}
		}
		if v.Default.Valid && len(v.Default.String) > 0 {
			ret.Defaults[v.Charset.String] = len(coll)
		}
		coll = append(coll, v)
		ret.Collations[v.Charset.String] = coll
	}
	return ret, nil
}

func (m *mySQL) getCollation(dbName string, collations *Collations) (string, error) {
	var err error
	if collations == nil {
		collations, err = m.getCollations()
		if err != nil {
			return ``, err
		}
	}
	sqlStr := "SHOW CREATE DATABASE " + quoteCol(dbName)
	row, err := m.newParam().SetCollection(sqlStr).QueryRow()
	if err != nil {
		return ``, err
	}
	var database sql.NullString
	var createDb sql.NullString
	err = row.Scan(&database, &createDb)
	if err != nil {
		return ``, err
	}
	matches := reCollate.FindStringSubmatch(createDb.String)
	if len(matches) > 1 {
		return matches[1], nil
	}
	matches = reCharacter.FindStringSubmatch(createDb.String)
	if len(matches) > 1 {
		if idx, ok := collations.Defaults[matches[1]]; ok {
			return collations.Collations[matches[1]][idx].Collation.String, nil
		}
	}

	return ``, nil
}

func (m *mySQL) getTableStatus(dbName string, tableName string, fast bool) (map[string]*TableStatus, []string, error) {
	sqlStr := `SHOW TABLE STATUS`
	if len(dbName) > 0 {
		sqlStr += " FROM " + quoteCol(dbName)
	}
	if len(tableName) > 0 {
		tableName = quoteVal(tableName, '_', '%')
		sqlStr += ` LIKE ` + tableName
	}
	ret := map[string]*TableStatus{}
	sorts := []string{}
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return ret, sorts, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return ret, sorts, err
	}
	n := len(cols)
	for rows.Next() {
		v := &TableStatus{}
		err := safeScan(rows, n,
			&v.Name, &v.Engine, &v.Version, &v.Row_format, &v.Rows, &v.Avg_row_length, &v.Data_length, &v.Max_data_length, &v.Index_length,
			&v.Data_free, &v.Auto_increment, &v.Create_time, &v.Update_time, &v.Check_time, &v.Collation, &v.Checksum, &v.Create_options,
			&v.Comment, &v.Max_index_length, &v.Temporary)
		if err != nil {
			return ret, sorts, fmt.Errorf(`%v: %v`, sqlStr, err)
		}
		if v.Engine.String == `InnoDB` {
			v.Comment.String = reInnoDBComment.ReplaceAllString(v.Comment.String, `$1`)
		}
		ret[v.Name.String] = v
		sorts = append(sorts, v.Name.String)
		if len(tableName) > 0 {
			return ret, sorts, nil
		}
	}
	return ret, sorts, nil
}

func (m *mySQL) getEngines() ([]*SupportedEngine, error) {
	sqlStr := `SHOW ENGINES`
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	n := len(cols)
	ret := []*SupportedEngine{}
	for rows.Next() {
		v := &SupportedEngine{}
		err := safeScan(rows, n, &v.Engine, &v.Support, &v.Comment, &v.Transactions, &v.XA, &v.Savepoints)
		if err != nil {
			return nil, fmt.Errorf(`%v: %v`, sqlStr, err)
		}
		if v.Support.String == `YES` || v.Support.String == `DEFAULT` {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func (m *mySQL) getVersion() string {
	if len(m.version) > 0 {
		return m.version
	}
	sqlStr := `SELECT version()`
	row, err := m.newParam().SetCollection(sqlStr).QueryRow()
	if err != nil {
		return fmt.Sprintf(`%v: %v`, sqlStr, err)
	}
	var v sql.NullString
	err = row.Scan(&v)
	if err != nil {
		return fmt.Sprintf(`%v: %v`, sqlStr, err)
	}
	m.version = v.String
	return v.String
}

func (m *mySQL) baseInfo() error {
	if m.Get(`dbList`) == nil {
		dbList, err := m.getDatabases()
		if err != nil {
			m.fail(err.Error())
			return m.returnTo(`/db`)
		}
		m.Set(`dbList`, dbList)
	}
	if len(m.dbName) > 0 {
		tableList, err := m.getTables()
		if err != nil {
			m.fail(err.Error())
			return m.returnTo(m.GenURL(`listDb`))
		}
		m.Set(`tableList`, tableList)
	}

	m.Set(`dbVersion`, m.getVersion())
	return nil
}
