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
	"strings"

	"github.com/webx-top/com"
)

func (m *mySQL) listData(
	callback func(columns []string, row map[string]*sql.NullString) error,
	table string, selectFuncs []string, selectCols []string,
	wheres []string, orderFields []string, descs []string,
	page int, limit int, totalRows int, textLength ...int) (columns []string, values []map[string]*sql.NullString, total int, err error) {
	var (
		groups  []string
		selects []string
		orders  []string
	)
	total = totalRows
	descNum := len(descs)
	funcNum := len(selectFuncs)
	for index, colName := range orderFields {
		if len(colName) == 0 {
			continue
		}
		if index >= descNum {
			continue
		}
		var order string
		if reSQLValue.MatchString(colName) {
			order = colName
		} else {
			order = quoteCol(colName)
		}
		if descs[index] == `1` {
			order += ` DESC`
		}
		orders = append(orders, order)
	}
	for index, colName := range selectCols {
		var (
			fn         string
			isGrouping bool
			sel        string
		)
		if index < funcNum {
			for _, f := range functions {
				if f == selectFuncs[index] {
					fn = f
					break
				}
			}
			for _, f := range grouping {
				if f == selectFuncs[index] {
					fn = f
					isGrouping = true
					break
				}
			}
		}
		if len(fn) == 0 && len(colName) == 0 {
			continue
		}
		if len(colName) == 0 {
			colName = `*`
		}
		sel = applySQLFunction(fn, quoteCol(colName))
		if !isGrouping {
			groups = append(groups, sel)
		}
		selects = append(selects, sel)
	}
	var fieldStr string
	if len(selects) > 0 {
		fieldStr = strings.Join(selects, `, `)
	} else {
		fieldStr = `*`
	}
	r := &Result{}
	var whereStr string
	if len(wheres) > 0 {
		whereStr += "\nWHERE " + strings.Join(wheres, ` AND `)
	}
	isGroup := len(groups) > 0 && len(groups) < len(selects)
	if isGroup {
		whereStr += "\nGROUP BY " + strings.Join(groups, `, `)
	}
	if len(orders) > 0 {
		whereStr += "\nORDER BY " + strings.Join(orders, `, `)
	}
	r.SQL = `SELECT` + withLimit(fieldStr+` FROM `+quoteCol(table), whereStr, limit, (page-1)*limit, "\n")
	if totalRows < 1 {
		countSQL := m.countRows(table, wheres, isGroup, groups)
		row := m.newParam().SetCollection(countSQL).QueryRow()
		err = row.Scan(&totalRows)
		if err != nil {
			return
		}
		total = totalRows
	}
	r.Query(m.newParam(), func(rows *sql.Rows) error {
		if callback == nil {
			columns, values, err = m.selectTable(rows, limit, textLength...)
		} else {
			columns, err = m.selectNext(rows, callback, limit, textLength...)
		}
		return err
	})
	m.AddResults(r)
	return
}

func (m *mySQL) selectTable(rows *sql.Rows, limit int, textLength ...int) (columns []string, r []map[string]*sql.NullString, err error) {
	r = []map[string]*sql.NullString{}
	columns, err = m.selectNext(rows, func(_ []string, row map[string]*sql.NullString) error {
		r = append(r, row)
		return nil
	}, limit, textLength...)
	return
}

func (m *mySQL) selectNext(rows *sql.Rows, callback func(columns []string, row map[string]*sql.NullString) error, limit int, textLength ...int) (columns []string, err error) {
	columns, err = rows.Columns()
	if err != nil {
		return
	}
	size := len(columns)
	var maxLen int
	if len(textLength) > 0 {
		maxLen = textLength[0]
	}
	for i := 0; i < limit && rows.Next(); i++ {
		values := make([]interface{}, size)
		for k := range columns {
			values[k] = &sql.NullString{}
		}
		err = rows.Scan(values...)
		if err != nil {
			return
		}
		val := map[string]*sql.NullString{}
		for k, colName := range columns {
			val[colName] = values[k].(*sql.NullString)
			if maxLen > 0 {
				val[colName].String = com.Substr(val[colName].String, ` ...`, maxLen)
			}
		}
		err = callback(columns, val)
		if err != nil {
			return
		}
	}
	return
}

func (m *mySQL) countRows(table string, wheres []string, isGroup bool, groups []string) string {
	query := " FROM " + quoteCol(table)
	if len(wheres) > 0 {
		query += " WHERE " + strings.Join(wheres, " AND ")
	}
	var groupBy string
	if isGroup {
		if m.supportSQL || len(groups) > 0 {
			return "SELECT COUNT(DISTINCT " + strings.Join(groups, ", ") + ")" + query
		}
		return "SELECT COUNT(*) FROM (SELECT 1" + query + groupBy + ") x"
	}
	return "SELECT COUNT(*)" + query
}

func withLimit(query string, where string, limit int, offset int, separator string) string {
	r := " " + query + where
	if limit > -1 {
		r += separator + fmt.Sprintf("LIMIT %d", limit)
		if offset > 0 {
			r += fmt.Sprintf(" OFFSET %d", offset)
		}
	}
	return r
}

/** Formulate SQL modification query with limit 1
* @param string everything after UPDATE or DELETE
* @param string
* @return string
 */
func withLimit1(query, where string) string {
	return withLimit(query, where, 1, 0, " ")
}

/** Delete data from table
* @param string
* @param string " WHERE ..."
* @param int 0 or 1
* @return bool
 */
func (m *mySQL) delete(table, queryWhere string, limit int) error {
	query := "FROM " + quoteCol(table)
	r := &Result{}
	r.SQL = `DELETE`
	if limit > 0 {
		r.SQL += withLimit1(query, queryWhere)
	} else {
		r.SQL += " " + query + queryWhere
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

func (m *mySQL) dumpHeaders(exportFormat string, multiTable bool) string {
	output := m.Form(`output`)
	var ext string
	if strings.Contains(exportFormat, `sql`) {
		ext = "sql"
	} else if multiTable {
		ext = "tar"
	} else {
		ext = "csv"
	}
	// multiple CSV packed to TAR
	var contentType string
	if output == `gz` {
		contentType = "application/x-gzip"
	} else if ext == `tar` {
		contentType = "application/x-tar"
	} else if ext == "sql" || output != `file` {
		contentType = "text/plain"
	} else {
		contentType = "text/csv"
	}
	m.Response().Header().Set("Content-Type", contentType+"; charset=utf-8")
	return ext
}

func (m *mySQL) exportData(fields map[string]*Field, table string, selectFuncs []string, selectCols []string, wheres []string, orderFields []string, descs []string, page int, limit int, totalRows int, textLength ...int) error {
	exportFormat := m.Form(`exportFormat`)
	exportStyle := m.Form(`exportStyle`)
	if exportFormat == `sql` {
		if exportStyle == `TRUNCATE+INSERT` {
			m.Response().Write(com.Str2bytes("TRUNCATE " + quoteCol(table) + ";\n"))
		}
	}
	var insert string
	var buffer string
	var suffix string
	var maxPacket int
	if m.Driver != `sqlite` {
		maxPacket = 1048576 // default, minimum is 1024
	}
	download := m.Formx(`download`).String()
	if len(download) > 0 {
		switch download {
		case `1`, `true`:
			m.Response().Header().Set("Content-Disposition", "attachment; filename="+friendlyURL(table+`-`))
		case `gzip`:
			//TODO
		}
	}
	ext := m.dumpHeaders(exportFormat, false)
	_ = ext
	_, _, _, err := m.listData(func(cols []string, row map[string]*sql.NullString) error {
		if exportFormat != `sql` {
			if exportStyle == `table` {
				dumpCSV(true, fields, cols, row, exportFormat, m.Response())
				exportStyle = `insert`
			}
			dumpCSV(false, fields, cols, row, exportFormat, m.Response())
			return nil
		}
		if len(insert) == 0 {
			keys := make([]string, len(cols))
			vals := make([]string, len(cols))
			for idx, key := range cols {
				key = quoteCol(key)
				keys[idx] = key
				vals[idx] = key + " = VALUES(" + key + ")"
			}
			if exportStyle == `INSERT+UPDATE` {
				suffix = "\nON DUPLICATE KEY UPDATE " + strings.Join(vals, ", ")
			}
			suffix += ";\n"
			insert = "INSERT INTO " + quoteCol(table) + " (" + strings.Join(keys, `, `) + ") VALUES"
		}
		var values, sep string
		for _, col := range cols {
			val := row[col]
			if !val.Valid {
				values += sep + `NULL`
			} else {
				field, ok := fields[col]
				var v string
				if ok && reFieldTypeNumber.MatchString(field.Type) && len(val.String) > 0 && !strings.HasPrefix(field.Full_type, `[`) {
					v = val.String
				} else {
					v = field.Format(val.String)
					v = quoteVal(v)
				}
				values += sep + unconvertField(field, v)
			}
			sep = `, `
		}
		var s string
		if maxPacket > 0 {
			s = "\n"
		} else {
			s = " "
		}
		s += "(" + values + ")"
		if len(buffer) == 0 {
			buffer = insert + s
		} else if len(buffer)+4+len(s)+len(suffix) < maxPacket { // 4 - length specification
			buffer += "," + s
		} else {
			m.Response().Write(com.Str2bytes(buffer + suffix))
			buffer = insert + s
		}
		return nil
	}, table, selectFuncs, selectCols, wheres, orderFields, descs, page, limit, totalRows, textLength...)

	if len(buffer) > 0 {
		m.Response().Write(com.Str2bytes(buffer + suffix))
	}
	return err
}

/** Update data in table
* @param string
* @param array escaped columns in keys, quoted data in values
* @param string " WHERE ..."
* @param int 0 or 1
* @param string
* @return bool
 */
func (m *mySQL) update(table string, set map[string]string, queryWhere string, limit int, separators ...string) error {
	separator := "\n"
	if len(separators) > 0 {
		separator = separators[0]
	}
	values := []string{}
	for key, val := range set {
		values = append(values, quoteCol(key)+"="+val)
	}
	query := quoteCol(table) + " SET" + separator + strings.Join(values, ","+separator)

	r := &Result{}
	r.SQL = "UPDATE"
	if limit > 0 {
		r.SQL += withLimit1(query, queryWhere)
	} else {
		r.SQL += " " + query + queryWhere
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

func (m *mySQL) set(table, queryWhere string, key string, value string, limit int) error {
	r := &Result{}
	query := quoteCol(table) + " SET " + quoteCol(key) + "=" + quoteVal(value)
	r.SQL = `UPDATE`
	if limit > 0 {
		r.SQL += withLimit1(query, queryWhere)
	} else {
		r.SQL += " " + query + queryWhere
	}
	r.Exec(m.newParam())
	m.Logger().Debug(r.SQL)
	//m.AddResults(r)
	return r.err
}

/** Insert data into table
* @param string
* @param array escaped columns in keys, quoted data in values
* @return bool
 */
func (m *mySQL) insert(table string, set map[string]string) error {
	r := &Result{}
	r.SQL = "INSERT INTO " + quoteCol(table)
	keys := []string{}
	vals := []string{}
	for key, val := range set {
		keys = append(keys, quoteCol(key))
		vals = append(vals, val)
	}
	if len(keys) > 0 {
		r.SQL += " (" + strings.Join(keys, ", ") + ")\nVALUES(" + strings.Join(vals, ", ") + ")"
	} else {
		r.SQL += " DEFAULT VALUES"
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}
