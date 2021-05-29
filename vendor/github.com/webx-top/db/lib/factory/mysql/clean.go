package mysql

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// Clean 清理掉表中不存在的字段
func Clean(linkID int, dbName string, tableName string, data echo.H, excludeFields ...string) (echo.H, error) {
	columns, _, err := ColumnComment(linkID, dbName, tableName)
	if err != nil {
		return data, err
	}
	return CleanBy(columns, linkID, dbName, tableName, data, excludeFields...)
}

// CleanBy 清理掉表中不存在的字段
func CleanBy(columns map[string]param.StringMap, linkID int, dbName string, tableName string, data echo.H, excludeFields ...string) (echo.H, error) {
	if len(excludeFields) > 0 {
		for _, field := range excludeFields {
			if data.Has(field) {
				delete(data, field)
			}
		}
	}
	for field := range data {
		if _, ok := columns[field]; !ok {
			delete(data, field)
		}
	}
	return data, nil
}
