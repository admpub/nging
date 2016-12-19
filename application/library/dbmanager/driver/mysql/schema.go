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
