package mysql

import (
	"database/sql"

	"strings"

	"github.com/admpub/nging/application/library/common"
	"github.com/webx-top/db/lib/factory"
)

func (m *mySQL) kvVal(sqlStr string) ([]map[string]string, error) {
	r := []map[string]string{}
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return r, err
	}
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
