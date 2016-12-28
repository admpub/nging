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

	"github.com/webx-top/com"
)

type KV struct {
	Value string
	Text  string
}

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
			"_Global_": map[string]string{
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
			p.privileges["_Global_"][priv.Privilege.String] = priv.Comment.String
			continue
		}
		for _, context := range strings.Split(priv.Context.String, `,`) {
			context = strings.Replace(context, ` `, `_`, -1)
			if _, ok := p.privileges[context]; !ok {
				p.privileges[context] = map[string]string{}
			}
			p.privileges[context][priv.Privilege.String] = priv.Comment.String
		}
	}
	//com.Dump(p.privileges)
	if _, ok := p.privileges["Server_Admin"]; !ok {
		p.privileges["Server_Admin"] = map[string]string{}
	}
	if vs, ok := p.privileges["File_access_on_server"]; ok {
		for k, v := range vs {
			p.privileges["Server_Admin"][k] = v
		}
	}
	if _, ok := p.privileges["Server_Admin"]["Usage"]; ok {
		delete(p.privileges["Server_Admin"], "Usage")
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

type Operation struct {
	Revoke  []string
	Grant   []string
	Columns string
	On      string
	User    string
	Scope   string //all|database|table|column|proxy
}

type Grant struct {
	Scope    string //all|database|table|column|proxy
	Value    string //*.*|db.*|db.table|db.table(col1,col2)
	Database string
	Table    string
	Columns  string            //col1,col2
	Settings map[string]string //["CREATE"]="1|0"
	*Operation
}

func (op *Operation) Apply(m *mySQL) error {
	//解除权限
	if len(op.Revoke) > 0 {
		on := `ON ` + op.On + ` FROM ` + op.User
		hasAll := op.HasAllPrivileges(&op.Revoke, true)
		hasOpt := op.HasGrantOption(&op.Revoke, true)
		if hasAll {
			if op.Scope == `proxy` {
				r := &Result{}
				r.SQL = `REVOKE PROXY ` + on
				r.Exec(m.newParam())
				m.AddResults(r)
				return r.err
			}
			r := &Result{}
			r.SQL = `REVOKE ALL PRIVILEGES ` + on
			r.Exec(m.newParam())
			m.AddResults(r)
			if r.err != nil {
				return r.err
			}
			if hasOpt {
				r := &Result{}
				r.SQL = `REVOKE GRANT OPTION ` + on
				r.Exec(m.newParam())
				m.AddResults(r)
				return r.err
			}
		}
		if hasOpt {
			if op.Scope == `proxy` {
				r := &Result{}
				r.SQL = `REVOKE PROXY ` + on
				r.Exec(m.newParam())
				m.AddResults(r)
				return r.err
			}
			r := &Result{}
			r.SQL = `REVOKE GRANT OPTION ` + on
			r.Exec(m.newParam())
			m.AddResults(r)
			if r.err != nil {
				return r.err
			}
		}
		if len(op.Revoke) > 0 {
			r := &Result{}
			c := strings.Join(op.Revoke, op.Columns+`, `) + op.Columns
			r.SQL = `REVOKE ` + reGrantOptionValue.ReplaceAllString(c, `$1`) + ` ` + on
			r.Exec(m.newParam())
			m.AddResults(r)
			if r.err != nil {
				return r.err
			}
		}
	}
	//增加授权
	if len(op.Grant) > 0 {
		on := `ON ` + op.On + ` TO ` + op.User
		if op.Scope == `proxy` {
			r := &Result{}
			r.SQL = `GRANT PROXY ` + on
			if op.HasGrantOption(&op.Grant, false) {
				r.SQL += ` WITH GRANT OPTION`
			}
			r.Exec(m.newParam())
			m.AddResults(r)
			return r.err
		}
		hasAll := op.HasAllPrivileges(&op.Grant, true)
		hasOpt := op.HasGrantOption(&op.Grant, true)
		if hasAll && hasOpt {
			r := &Result{}
			r.SQL = `GRANT ALL PRIVILEGES ` + on + ` WITH GRANT OPTION`
			r.Exec(m.newParam())
			m.AddResults(r)
			return r.err
		}
		if hasAll {
			r := &Result{}
			r.SQL = `GRANT ALL PRIVILEGES ` + on
			r.Exec(m.newParam())
			m.AddResults(r)
			return r.err
		}
		if hasOpt {
			r := &Result{}
			r.SQL = `GRANT GRANT OPTION ` + on
			r.Exec(m.newParam())
			m.AddResults(r)
			if r.err != nil {
				return r.err
			}
		}
		if len(op.Grant) > 0 {
			r := &Result{}
			c := strings.Join(op.Grant, op.Columns+`, `) + op.Columns
			r.SQL = `GRANT ` + reGrantOptionValue.ReplaceAllString(c, `$1`) + ` ` + on
			r.Exec(m.newParam())
			m.AddResults(r)
			if r.err != nil {
				return r.err
			}
		}
	}
	return nil
}

func (op *Operation) HasAllPrivileges(values *[]string, deleteIt bool) bool {
	for index, name := range *values {
		if name == `ALL PRIVILEGES` {
			if deleteIt {
				if index+1 < len(*values) {
					*values = append((*values)[0:index], (*values)[index+1:]...)
				} else {
					*values = (*values)[0:index]
				}
			}
			return true
		}
	}
	return false
}

func (op *Operation) HasGrantOption(values *[]string, deleteIt bool) bool {
	for index, name := range *values {
		if name == `GRANT OPTION` {
			if deleteIt {
				if index+1 < len(*values) {
					*values = append((*values)[0:index], (*values)[index+1:]...)
				} else {
					*values = (*values)[0:index]
				}
			}
			return true
		}
	}
	return false
}

func (g *Grant) IsValid(group string, values map[string]*Mapx) bool {
	if group == `_Global_` {
		return true
	}
	switch g.Scope {
	case `all`:
		return group == `Server_Admin` || group == `Procedures`
	case `database`:
		return group == `Databases` || group == `Procedures`
	case `table`:
		return group == `Tables`
	case `column`:
		return group == `Columns`
	case `proxy`:
		return false
	default:
		return false
	}
}

func (g *Grant) String() string {
	switch g.Scope {
	case `proxy`:
		r := strings.SplitN(g.Value, `@`, 2)
		if len(r) != 2 {
			return ``
		}
		r[0] = strings.Trim(r[0], `'`)
		r[1] = strings.Trim(r[1], `'`)
		return quoteVal(r[0]) + `@` + quoteVal(r[1])
	case `all`:
		return `*.*`
	case `database`:
		g.Database = reNotWord.ReplaceAllString(g.Database, ``)
		if len(g.Database) == 0 {
			return ``
		}
		return "`" + g.Database + "`.*"
	case `table`:
		g.Database = reNotWord.ReplaceAllString(g.Database, ``)
		if len(g.Database) == 0 {
			return ``
		}
		g.Table = reNotWord.ReplaceAllString(g.Table, ``)
		if len(g.Table) == 0 {
			return ``
		}
		return "`" + g.Database + "`.`" + g.Table + "`"
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
		if len(g.Columns) > 0 {
			return "`" + g.Database + "`.`" + g.Table + "` (" + g.Columns + ")"
		}
		return "`" + g.Database + "`.`" + g.Table + "`"
	}
	return ``
}

type SupportedEngine struct {
	Engine       sql.NullString //CSV|InnoDB|MyISAM|MEMORY...
	Support      sql.NullString //YES|DEFAULT
	Comment      sql.NullString
	Transactions sql.NullString //NO|YES
	XA           sql.NullString //NO|YES
	Savepoints   sql.NullString //NO|YES
}

type FieldInfo struct {
	Field      sql.NullString
	Type       sql.NullString
	Collation  sql.NullString
	Null       sql.NullString
	Key        sql.NullString
	Default    sql.NullString
	Extra      sql.NullString
	Privileges sql.NullString
	Comment    sql.NullString
}

type IndexInfo struct {
	Table         sql.NullString
	Non_unique    sql.NullString
	Key_name      sql.NullString
	Seq_in_index  sql.NullString
	Column_name   sql.NullString
	Collation     sql.NullString
	Cardinality   sql.NullString
	Sub_part      sql.NullString
	Packed        sql.NullString
	Null          sql.NullString
	Index_type    sql.NullString
	Comment       sql.NullString
	Index_comment sql.NullString
}

type Indexes struct {
	Type    string
	Columns []string
	Lengths []string
	Descs   []string
}

type Field struct {
	Field          string
	Full_type      string
	Type           string
	Length         string
	Unsigned       string
	Default        sql.NullString
	Null           bool
	Auto_increment bool
	On_update      string
	Collation      string
	Privileges     map[string]int
	Comment        string
	Primary        bool
}
