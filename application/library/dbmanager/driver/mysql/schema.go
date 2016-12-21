package mysql

import (
	"database/sql"
	"strings"

	"github.com/webx-top/com"
)

type ProcessList struct {
	Id       sql.NullInt64
	User     sql.NullString
	Host     sql.NullString
	Db       sql.NullString
	Command  sql.NullString
	Time     sql.NullInt64
	State    sql.NullString
	Info     sql.NullString
	Progress sql.NullFloat64
}

type TableStatus struct {
	Name            sql.NullString
	Engine          sql.NullString
	Version         sql.NullString
	Row_format      sql.NullString
	Rows            sql.NullInt64
	Avg_row_length  sql.NullInt64
	Data_length     sql.NullInt64
	Max_data_length sql.NullInt64
	Index_length    sql.NullInt64
	Data_free       sql.NullInt64
	Auto_increment  sql.NullInt64
	Create_time     sql.NullString
	Update_time     sql.NullString
	Check_time      sql.NullString
	Collation       sql.NullString
	Checksum        sql.NullString
	Create_options  sql.NullString
	Comment         sql.NullString
}

func (t *TableStatus) IsView() bool {
	return t.Engine.Valid == false
}

func (t *TableStatus) FKSupport(currentVersion string) bool {
	switch t.Engine.String {
	case `InnoDB`, `IBMDB2I`, `NDB`:
		if com.VersionCompare(currentVersion, `5.6`) >= 0 {
			return true
		}
	}
	return false
}

func (t *TableStatus) Size() int64 {
	return t.Data_length.Int64 + t.Index_length.Int64
}

type Collation struct {
	Collation sql.NullString
	Charset   sql.NullString `json:"-"`
	Id        sql.NullInt64  `json:"-"`
	Default   sql.NullString `json:"-"`
	Compiled  sql.NullString `json:"-"`
	Sortlen   sql.NullInt64  `json:"-"`
}

type Collations struct {
	Collations map[string][]Collation
	Defaults   map[string]int
}

func NewCollations() *Collations {
	return &Collations{
		Collations: make(map[string][]Collation),
		Defaults:   make(map[string]int),
	}
}

type Privilege struct {
	Privilege sql.NullString
	Context   sql.NullString
	Comment   sql.NullString
}

func NewPrivileges() *Privileges {
	return &Privileges{
		Privileges: []*Privilege{},
		privileges: map[string]map[string]string{
			"": map[string]string{
				"All privileges": "",
			},
		},
	}
}

type Privileges struct {
	Privileges []*Privilege
	privileges map[string]map[string]string
}

func (p *Privileges) Parse() {
	for _, priv := range p.Privileges {
		if priv.Privilege.String == `Grant option` {
			p.privileges[""][priv.Privilege.String] = priv.Comment.String
			continue
		}
		for _, context := range strings.Split(priv.Context.String, `,`) {
			if _, ok := p.privileges[context]; !ok {
				p.privileges[context] = map[string]string{}
			}
			p.privileges[context][priv.Privilege.String] = priv.Comment.String
		}
	}
	if _, ok := p.privileges["Server Admin"]; !ok {
		p.privileges["Server Admin"] = map[string]string{}
	}
	if vs, ok := p.privileges["File access on server"]; ok {
		for k, v := range vs {
			p.privileges["Server Admin"][k] = v
		}
	}
	if _, ok := p.privileges["Server Admin"]["Usage"]; ok {
		delete(p.privileges["Server Admin"], "Usage")
	}
	if _, ok := p.privileges["Databases"]; !ok {
		p.privileges["Databases"] = map[string]string{}
	}
	if vs, ok := p.privileges["Procedures"]; ok {
		if v, ok := vs["Create routine"]; ok {
			p.privileges["Databases"]["Create routine"] = v
			delete(p.privileges["Procedures"], "Create routine")
		}
	}
	if _, ok := p.privileges["Tables"]; ok {
		p.privileges["Columns"] = map[string]string{}
		for _, val := range []string{`Select`, `Insert`, `Update`, `References`} {
			if v, y := p.privileges["Tables"][val]; y {
				p.privileges["Columns"][val] = v
			}
		}
		for k := range p.privileges["Tables"] {
			if _, ok := p.privileges["Databases"][k]; !ok {
				delete(p.privileges["Databases"], k)
			}
		}
	}
}

type Mapx struct {
	Map map[string]*Mapx
	Val []string
}

// ParseFormName user[name][test]
func ParseFormName(s string) []string {
	var res []string
	hasLeft := false
	hasRight := true
	var val []rune
	for _, r := range s {
		if r == '[' {
			if hasRight {
				res = append(res, string(val))
				val = []rune{}
			}
			hasLeft = true
			hasRight = false
			continue
		}
		if r == ']' {
			if hasLeft {
				hasRight = true
			}
			continue
		}
		val = append(val, r)
	}
	if len(val) > 0 {
		res = append(res, string(val))
	}
	return res
}

func NewMapx(data map[string][]string) *Mapx {
	m := &Mapx{}
	return m.Parse(data)
}

func (m *Mapx) Parse(data map[string][]string) *Mapx {
	m.Map = map[string]*Mapx{}
	for name, values := range data {
		names := ParseFormName(name)
		end := len(names) - 1
		v := m
		for idx, key := range names {
			if _, ok := v.Map[key]; !ok {
				if idx == end {
					v.Map[key] = &Mapx{Val: values}
					continue
				}
				v.Map[key] = &Mapx{
					Map: map[string]*Mapx{},
				}
				v = v.Map[key]
				continue
			}

			if idx == end {
				v.Map[key] = &Mapx{Val: values}
			} else {
				v = v.Map[key]
			}
		}
	}
	return m
}

func (m *Mapx) Value(names ...string) string {
	v := m.Values(names...)
	if v != nil {
		if len(v) > 0 {
			return v[0]
		}
	}
	return ``
}

func (m *Mapx) Values(names ...string) []string {
	if len(names) == 0 {
		if m.Val == nil {
			return []string{}
		}
		return m.Val
	}
	v := m.Get(names...)
	if v != nil {
		return v.Val
	}
	return []string{}
}

func (m *Mapx) Get(names ...string) *Mapx {
	v := m
	end := len(names) - 1
	for idx, key := range names {
		if _, ok := v.Map[key]; !ok {
			return nil
		}
		v = v.Map[key]

		if idx == end {
			return v
		}
	}
	return nil
}

type Grant struct {
	Scope    string //all|database|table|column
	Value    string //*.*|db.*|db.table|db.table(col1,col2)
	Database string
	Table    string
	Columns  string //col1,col2
}

func (g *Grant) String() string {
	switch g.Scope {
	case `all`:
		return `*.*`
	case `database`:
		g.Database = reNotWord.ReplaceAllString(g.Database, ``)
		if len(g.Database) == 0 {
			return ``
		}
		return g.Database + `.*`
	case `table`:
		g.Database = reNotWord.ReplaceAllString(g.Database, ``)
		if len(g.Database) == 0 {
			return ``
		}
		g.Table = reNotWord.ReplaceAllString(g.Table, ``)
		if len(g.Table) == 0 {
			return ``
		}
		return g.Database + `.` + g.Table
	case `column`:
		g.Database = reNotWord.ReplaceAllString(g.Database, ``)
		if len(g.Database) == 0 {
			return ``
		}
		g.Table = reNotWord.ReplaceAllString(g.Table, ``)
		if len(g.Table) == 0 {
			return ``
		}
		columns := strings.Split(g.Columns, `,`)
		g.Columns = ``
		var sep string
		for _, column := range columns {
			column = reNotWord.ReplaceAllString(column, ``)
			if len(column) == 0 {
				continue
			}
			g.Columns += sep + column
			sep = `,`
		}
		return g.Database + `.` + g.Table + `(` + g.Columns + `)`
	}
	return ``
}
