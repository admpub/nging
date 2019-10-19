package table

type TableInfoStorer interface {
	SetTableID(string) TableInfoStorer
	SetTableName(string) TableInfoStorer
	SetFieldName(string) TableInfoStorer
	TableID() string
	TableName() string
	FieldName() string
}

var _ TableInfoStorer = &TableInfo{}

type TableInfo struct {
	tableName string
	tableID   string
	fieldName string
}

func (t *TableInfo) SetTableID(tableID string) TableInfoStorer {
	t.tableID = tableID
	return t
}
func (t *TableInfo) SetTableName(tableName string) TableInfoStorer {
	t.tableName = tableName
	return t
}
func (t *TableInfo) SetFieldName(fieldName string) TableInfoStorer {
	t.fieldName = fieldName
	return t
}
func (t *TableInfo) TableID() string {
	return t.tableID
}
func (t *TableInfo) TableName() string {
	return t.tableName
}
func (t *TableInfo) FieldName() string {
	return t.fieldName
}
