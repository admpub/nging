package table


type TableInfoStorer interface {
	SetTableID(uint64) TableInfoStorer
	SetTableName(string) TableInfoStorer
	SetFieldName(string) TableInfoStorer
}
