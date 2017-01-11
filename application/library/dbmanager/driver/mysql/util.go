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

	"strings"

	"regexp"

	"github.com/admpub/nging/application/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
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
	m.Session().AddFlash(common.Ok(msg))
}

func (m *mySQL) checkErr(err error) interface{} {
	return common.Err(m.Context, err)
}

func (m *mySQL) fail(msg string) {
	m.Session().AddFlash(msg)
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
	return "`" + com.AddSlashes(col, '`') + "`"
}

func quoteVal(val string, otherChars ...rune) string {
	return "'" + com.AddSlashes(val, otherChars...) + "'"
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

var reFunction1 = regexp.MustCompile(`^([+-]|\\|)$`)
var reFunction2 = regexp.MustCompile(`^[+-] interval$`)
var reSQLValue = regexp.MustCompile(`^(\d+|'[0-9.: -]') [A-Z_]+$`)

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
		if reFunction1.MatchString(function) {
			ret = quoteCol(field.Field) + ` ` + function + ` ` + ret
		} else if reFunction2.MatchString(function) {
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
