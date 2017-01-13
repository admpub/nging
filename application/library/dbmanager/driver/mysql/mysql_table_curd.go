package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/webx-top/com"
)

func (m *mySQL) selectTable(rows *sql.Rows, limit int, textLength ...int) (columns []string, r []map[string]*sql.NullString, err error) {
	columns, err = rows.Columns()
	r = []map[string]*sql.NullString{}
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
		r = append(r, val)
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
			return "SELECT COUNT(DISTINCT " + strings.Join(groups, ", ") + ")"
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
