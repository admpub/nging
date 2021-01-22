package internal

import (
	"fmt"
	"strings"
)

type AlterType int

const (
	AlterTypeNo     AlterType = 0
	AlterTypeCreate           = 1
	AlterTypeDrop             = 2
	AlterTypeAlter            = 3
)

func (at AlterType) String() string {
	switch at {
	case AlterTypeNo:
		return "not_change"
	case AlterTypeCreate:
		return "create"
	case AlterTypeDrop:
		return "drop"
	case AlterTypeAlter:
		return "alter"
	default:
		return "unknow"
	}

}

func NewAlterData(tableName string) *TableAlterData {
	return &TableAlterData{
		Table: tableName,
		Type:  AlterTypeNo,
	}
}

// TableAlterData 表的变更情况
type TableAlterData struct {
	Table      string
	Type       AlterType
	SQL        string
	SchemaDiff *SchemaDiff
}

func (ta *TableAlterData) String() string {
	relationTables := ta.SchemaDiff.RelationTables()
	fmtStr := `
-- Table : %s
-- Type  : %s
-- RealtionTables : %s
-- SQL   : 
%s
`
	return fmt.Sprintf(fmtStr, ta.Table, ta.Type, strings.Join(relationTables, ","), ta.SQL)
}
