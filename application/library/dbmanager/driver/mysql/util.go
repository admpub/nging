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
