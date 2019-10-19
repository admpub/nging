package table

type TableInfoStorer interface {
	SetTableID(string) TableInfoStorer
	SetTableName(string) TableInfoStorer
	SetFieldName(string) TableInfoStorer
	TableID() string
	TableName() string
	FieldName() string
}
