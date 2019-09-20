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
	"bytes"
	"database/sql"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

func (m *mySQL) kvVal(sqlStr string) ([]map[string]string, error) {
	r := []map[string]string{}
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return r, err
	}
	defer rows.Close()
	for rows.Next() {
		var k sql.NullString
		var v sql.NullString
		err = rows.Scan(&k, &v)
		if err != nil {
			break
		}
		r = append(r, map[string]string{
			"k": k.String,
			"v": v.String,
		})
	}
	return r, err
}

func (m *mySQL) newParam() *factory.Param {
	return factory.NewParam(m.db)
}

func (m *mySQL) ok(msg string) {
	m.SetOk(msg)
}

func (m *mySQL) checkErr(err error) interface{} {
	return m.CheckErr(err)
}

func (m *mySQL) fail(msg string) {
	m.SetFail(msg)
}

func (m *mySQL) getScopeGrant(object string) *Grant {
	g := &Grant{Value: object}
	if object == `*.*` {
		g.Scope = `all`
		return g
	}
	if strings.Contains(object, `@`) {
		g.Scope = `proxy`
		return g
	}
	strs := strings.SplitN(object, `.`, 2)
	for i, v := range strs {
		v = strings.Trim(v, "`")
		switch i {
		case 0:
			g.Database = v
		case 1:
			if v == `*` {
				g.Scope = `database`
			} else if strings.HasSuffix(v, `)`) {
				vs := strings.SplitN(v, `(`, 2)
				switch len(vs) {
				case 2:
					g.Table = strings.TrimSpace(vs[0])
					g.Table = strings.TrimSuffix(g.Table, "`")
					g.Columns = strings.TrimSuffix(vs[1], `)`)
					g.Scope = `column`
				}
			} else {
				g.Table = strings.TrimSpace(v)
				g.Table = strings.TrimSuffix(g.Table, "`")
				g.Scope = `table`
			}
		}
	}
	return g
}

func quoteCol(col string) string {
	return "`" + strings.Replace(com.AddSlashes(col), "`", "``", -1) + "`"
}

func quoteVal(val string, otherChars ...rune) string {
	return "'" + com.AddSlashes(val, otherChars...) + "'"
}

func convertFields(columns []string, fields map[string]*Field, selects []string) string {
	var r string
	l := len(selects)
	for _, colName := range columns {
		quotedName := quoteCol(colName)
		if l > 0 {
			found := false
			for _, val := range selects {
				if quotedName == val {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		field, ok := fields[colName]
		if !ok {
			continue
		}
		as := convertField(field)
		if len(as) > 0 {
			r += ", " + as + " AS " + quotedName
		}
	}
	return r
}

/** Convert field in select and edit
* @param array one element from fields()
* @return string
 */
func convertField(field *Field) string {
	if strings.Contains(field.Type, "binary") {
		return "HEX(" + quoteCol(field.Field) + ")"
	}
	if field.Type == "bit" {
		return "BIN(" + quoteCol(field.Field) + " + 0)" // + 0 is required outside MySQLnd
	}
	switch {
	case strings.Contains(field.Type, "geometry"),
		strings.Contains(field.Type, "point"),
		strings.Contains(field.Type, "linestring"),
		strings.Contains(field.Type, "polygon"):
		return "AsWKT(" + quoteCol(field.Field) + ")"
	}
	return ``
}

/** Convert value in edit after applying functions back
* @param array one element from fields()
* @param string
* @return string
 */
func unconvertField(field *Field, ret string) string {

	if strings.Contains(field.Type, "binary") {
		return "UNHEX(" + ret + ")"
	}
	if field.Type == "bit" {
		return "CONV(" + ret + ", 2, 10) + 0"
	}
	switch {
	case strings.Contains(field.Type, "geometry"),
		strings.Contains(field.Type, "point"),
		strings.Contains(field.Type, "linestring"),
		strings.Contains(field.Type, "polygon"):
		return "GeomFromText(" + ret + ")"
	}
	return ret
}

/** Process edit input field
* @param one field from fields()
* @return string or false to leave the original value
 */
func (m *mySQL) processInputFieldValue(field *Field) (string, bool) {
	idf := bracketEscape(field.Field, false)
	if field.Type == "set" {
		total := 0
		for _, v := range m.FormValues("value[" + idf + "][]") {
			i, _ := strconv.Atoi(v)
			total += i
		}
		return strconv.Itoa(total), true
	}
	function := m.Form("function[" + idf + "]")
	value := m.Form("value[" + idf + "]")
	if field.Type == "enum" {
		i, _ := strconv.Atoi(value)
		if i == -1 {
			return ``, false
		}
		if len(value) == 0 {
			return "NULL", true
		}
		return strconv.Itoa(i), true
	}
	if field.AutoIncrement.Valid && len(value) == 0 {
		return ``, false
	}
	if function == "orig" {
		if field.On_update == "CURRENT_TIMESTAMP" {
			return quoteCol(field.Field), true
		}
		return ``, false
	}
	if function == "NULL" {
		return "NULL", true
	}
	if function == "json" {
		return value, true
	}
	if reFieldTypeBlob.MatchString(field.Type) {
		buf := new(bytes.Buffer)
		_, e := m.SaveUploadedFileToWriter("value["+idf+"]", buf)
		if e != nil {
			return ``, false
		}
		return quoteVal(buf.String()), true
	}
	return processInput(field, value, function), true
}

func processInput(field *Field, value string, function string) string {
	if function == "SQL" {
		return value // SQL injection
	}
	ret := quoteVal(value)
	switch function {
	case `now`, `getdate`, `uuid`:
		ret = function + `()`
	case `current_date`, `current_timestamp`:
		return function
	case `addtime`, `subtime`, `concat`:
		return function + `(` + quoteCol(field.Field) + `,` + ret + `)`
	case `md5`, `sha1`, `password`, `encrypt`:
		return function + `(` + ret + `)`
	default:
		if reFunctionAddOrSubOr.MatchString(function) {
			ret = quoteCol(field.Field) + ` ` + function + ` ` + ret
		} else if reFunctionInterval.MatchString(function) {
			ret2 := ret
			ret = quoteCol(field.Field) + ` ` + function + ` `
			if reSQLValue.MatchString(value) {
				ret += value
			} else {
				ret += ret2
			}
		}
	}
	return unconvertField(field, ret)
}

func getCharset(version string) string {
	if com.VersionCompare(version, `5.5.3`) >= 0 {
		return "utf8mb4"
	}
	return "utf8" // SHOW CHARSET would require an extra query
}

func applySQLFunction(function, column string) string {
	if len(function) > 0 {
		switch function {
		case `unixepoch`:
			return `DATETIME(` + column + `, '` + function + `')`
		case `count distinct`:
			return `COUNT(DISTINCT ` + column + `)`
		default:
			return strings.ToUpper(function) + `(` + column + `)`
		}
	}
	return column
}

/** Find unique identifier of a row
* @param array
* @param array result of indexes()
* @return array or null if there is no unique identifier
 */
func uniqueArray(row map[string]*sql.NullString, indexes map[string]*Indexes) map[string]*sql.NullString {
	ret := map[string]*sql.NullString{}
	for _, index := range indexes {
		switch index.Type {
		case `PRIMARY`, `UNIQUE`:
			for _, key := range index.Columns {
				v, y := row[key]
				if y {
					ret[key] = v
					continue
				}
				break
			}
		}
	}
	return ret
}

/** Escape or unescape string to use inside form []
* @param string
* @param bool
* @return string
 */
func bracketEscape(idf string, back bool) string {
	// escape brackets inside name="x[]"
	if back {
		for k, v := range trans {
			idf = strings.Replace(idf, v, k, -1)
		}
		return idf
	}
	for k, v := range trans {
		idf = strings.Replace(idf, k, v, -1)
	}
	return idf
}

/** Escape column key used in where()
* @param string
* @return string
 */
func escapeKey(key string) string {
	if match := reFieldName.FindAllString(key, 1); len(match) > 3 {
		return match[1] + quoteCol(match[2]) + match[3] //! SQL injection
	}
	return quoteCol(key)
}

func (m *mySQL) editFunctions(field *Field) []string {
	var r string
	if field.AutoIncrement.Valid {
		r = m.T(`自动增量`)
	} else {
		if field.Null {
			r += "NULL/"
		}
		for key, functions := range editFunctions {
			if key == 0 {
				for pattern, value := range functions {
					if len(pattern) == 0 {
						r += "/" + value
					} else {
						re, err := regexp.Compile(pattern)
						if err != nil {
							m.Logger().Error(err)
							continue
						}
						if !re.MatchString(field.Type) {
							continue
						}

						r += "/" + value
					}
				}
				continue
			}
			switch field.Type {
			case `set`, `enum`:
			default:
				if !reFieldTypeBlob.MatchString(field.Type) {
					r += `/SQL`
				}
			}
		}
	}
	if len(r) > 0 {
		return strings.Split(r, `/`)
	}
	return []string{}
}

func (m *mySQL) whereByMapx(where *echo.Mapx, null *echo.Mapx, fields map[string]*Field) string {
	wheres := map[string]*echo.Mapx{}
	nulls := map[string]*echo.Mapx{}
	if where != nil {
		wheres = where.Map
	}
	if null != nil {
		nulls = null.Map
	}
	return m.where(wheres, nulls, fields)
}

func (m *mySQL) where(wheres map[string]*echo.Mapx, nulls map[string]*echo.Mapx, fields map[string]*Field) string {
	r := []string{}
	for key, mapx := range wheres {
		if mapx == nil {
			continue
		}
		key = bracketEscape(key, true)
		column := escapeKey(key)
		field, ok := fields[key]
		if !ok {
			continue
		}
		val := mapx.Value()
		if (m.DbAuth.Driver == `mssql`) || (m.supportSQL && reOnlyFloatOrEmpty.MatchString(val)) {
			r = append(r, column+" LIKE "+quoteVal(val, '%', '_'))
		} else {
			r = append(r, column+"="+unconvertField(field, quoteVal(val)))
		}
		/*
			if m.supportSQL &&
				(strings.Contains(field.Type, `char`) || strings.Contains(field.Type, `text`)) &&
				reNotSpaceOrDashOrAt.MatchString(val) {
				r = append(r, column+"="+quoteVal(val)+" COLLATE "+getCharset(m.getVersion())+"_bin")
			}
		*/
	}
	for key, mapx := range nulls {
		if mapx == nil {
			continue
		}
		key = mapx.Value()
		r = append(r, escapeKey(key)+" IS NULL")
	}
	return strings.Join(r, " AND ")

}

func enumValues(field *Field) []*Enum {
	r := []*Enum{}
	matches := reFieldEnumValue.FindAllStringSubmatch(field.Length, -1)
	//com.Dump(matches)
	if len(matches) > 0 {
		for i, val := range matches {
			val[1] = strings.Replace(val[1], `''`, `'`, -1)
			val[1] = strings.Replace(val[1], `\`, ``, -1)
			r = append(r, &Enum{
				Int:    enumNumber(i),
				String: val[1],
			})
		}
	}
	return r
}

func enumNumber(i int) int {
	return 1 << uint64(math.Abs(float64(i)))
}
