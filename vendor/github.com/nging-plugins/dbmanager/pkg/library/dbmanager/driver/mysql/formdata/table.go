package formdata

import (
	"database/sql"
	"strconv"

	"github.com/webx-top/echo/param"
)

type Table struct {
	Name               string
	Engine             string
	Collation          string
	AutoIncrementStart sql.NullInt64 `form_options:"-"`
	Auto_increment     string
	Auto_increment_n   int
	Ai_start_val       string
	Comment            string
	FieldIndexes       []string
	Fields             map[string]*Field
}

func (t *Table) Init() {
	autoIncrementStartValue := t.Ai_start_val
	t.AutoIncrementStart = sql.NullInt64{Valid: len(autoIncrementStartValue) > 0}
	if t.AutoIncrementStart.Valid {
		t.AutoIncrementStart.Int64, _ = strconv.ParseInt(autoIncrementStartValue, 10, 64)
	}
	t.Auto_increment_n = param.AsInt(t.Auto_increment)
}

type Field struct {
	Field         string
	Orig          string // 字段旧名称
	Type          string
	Length        string
	Unsigned      string
	Collation     string
	On_delete     string
	On_update     string
	Null          bool
	Comment       string
	Default       string
	Has_default   bool
	AutoIncrement sql.NullString `form_options:"-"`
}

func (f *Field) Init(t *Table, index string) {
	f.AutoIncrement = sql.NullString{
		Valid: t.Auto_increment == index,
	}
	if f.AutoIncrement.Valid {
		f.AutoIncrement.String = t.Ai_start_val
	}
}
