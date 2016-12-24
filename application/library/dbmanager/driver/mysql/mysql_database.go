package mysql

import (
	"database/sql"
	"strings"

	"github.com/webx-top/com"
)

func (m *mySQL) createDatabase(dbName, collate string) *Result {
	r := &Result{}
	r.SQL = "CREATE DATABASE `" + strings.Replace(dbName, "`", "", -1) + "`"
	if len(collate) > 0 {
		r.SQL += " COLLATE '" + com.AddSlashes(collate) + "'"
	}
	return r.Exec(m.newParam())
}

func (m *mySQL) dropDatabase(dbName string) *Result {
	r := &Result{}
	r.SQL = "DROP DATABASE `" + strings.Replace(dbName, "`", "", -1) + "`"
	return r.Exec(m.newParam())
}

func (m *mySQL) renameDatabase(newName, collate string) []*Result {
	newName = strings.Replace(newName, "`", "", -1)
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
	rGetTables.ErrorString = err.Error()
	rs = append(rs, rGetTables)
	if err != nil {
		return rs
	}
	var sql string
	for key, table := range tables {
		table = com.AddCSlashes(table, '`')
		if key > 0 {
			sql += ", "
		}
		sql += "`" + table + "` TO `" + newName + "`.`" + table + "`"
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
