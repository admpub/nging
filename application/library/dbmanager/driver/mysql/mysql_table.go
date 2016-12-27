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
	"errors"
	"strings"
)

// 获取数据表列表
func (m *mySQL) getTables() ([]string, error) {
	sqlStr := `SHOW TABLES`
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

func (m *mySQL) optimizeTables(tables []string, operation string) error {
	r := &Result{}
	defer m.AddResults(r)
	var op string
	switch operation {
	case `optimize`, `check`, `analyze`, `repair`:
		op = strings.ToUpper(operation)
	default:
		return errors.New(m.T(`不支持的操作: %s`, operation))
	}
	for _, table := range tables {
		table = quoteCol(table)
		r.SQL = op + ` TABLE ` + table
		r.Execs(m.newParam())
		if r.err != nil {
			return r.err
		}
	}
	r.end()
	return r.err
}

func (m *mySQL) moveTables(tables []string, targetDb string) error {
	r := &Result{}
	r.SQL = `RENAME TABLE `
	targetDb = quoteCol(targetDb)
	for i, table := range tables {
		table = quoteCol(table)
		if i > 0 {
			r.SQL += `,`
		}
		r.SQL += table + " TO " + targetDb + "." + table
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

//删除表
func (m *mySQL) dropTables(tables []string) error {
	r := &Result{}
	r.SQL = `DROP TABLE `
	for i, table := range tables {
		table = quoteCol(table)
		if i > 0 {
			r.SQL += `,`
		}
		r.SQL += table
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

//清空表
func (m *mySQL) truncateTables(tables []string) error {
	r := &Result{}
	defer m.AddResults(r)
	for _, table := range tables {
		table = quoteCol(table)
		r.SQL = `TRUNCATE TABLE ` + table
		r.Execs(m.newParam())
		if r.err != nil {
			return r.err
		}
	}
	r.end()
	return nil
}

type viewCreateInfo struct {
	View                 sql.NullString
	CreateView           sql.NullString
	Character_set_client sql.NullString
	Collation_connection sql.NullString
	Select               string
}

func (m *mySQL) tableView(name string) (*viewCreateInfo, error) {
	sqlStr := `SHOW CREATE VIEW ` + quoteCol(name)
	row, err := m.newParam().SetCollection(sqlStr).QueryRow()
	if err != nil {
		return nil, err
	}
	info := &viewCreateInfo{}
	err = row.Scan(&info.View, &info.CreateView, &info.Character_set_client, &info.Collation_connection)
	if err != nil {
		return info, err
	}
	info.Select = reView.ReplaceAllString(info.CreateView.String, ``)
	return info, nil
}

func (m *mySQL) copyTables(tables []string, targetDb string, isView bool) error {
	r := &Result{}
	r.SQL = `SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO'`
	r.Execs(m.newParam())
	m.AddResults(r)
	if r.err != nil {
		return r.err
	}
	same := m.dbName == targetDb
	targetDb = quoteCol(targetDb)
	for _, table := range tables {
		var name string
		quotedTable := quoteCol(table)
		if same {
			name = `copy_` + table
			name = quoteCol(name)
		} else {
			name = targetDb + "." + quotedTable
		}
		if isView {
			r.SQL = `DROP VIEW IF EXISTS ` + name
			r.Execs(m.newParam())
			if r.err != nil {
				return r.err
			}

			viewInfo, err := m.tableView(table)
			if err != nil {
				return err
			}
			r.SQL = `CREATE VIEW ` + name + ` AS ` + viewInfo.Select
			r.Execs(m.newParam())
			if r.err != nil {
				return r.err
			}
			continue

		}
		r.SQL = `DROP TABLE IF EXISTS ` + name
		r.Execs(m.newParam())
		if r.err != nil {
			return r.err
		}

		r.SQL = `CREATE TABLE ` + name + ` LIKE ` + quotedTable
		r.Execs(m.newParam())
		if r.err != nil {
			return r.err
		}

		r.SQL = `INSERT INTO ` + name + ` SELECT * FROM ` + quotedTable
		r.Execs(m.newParam())
		if r.err != nil {
			return r.err
		}
	}
	r.end()
	return r.err
}
