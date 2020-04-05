package table

import "strings"

// GetTableInfo 从上传所属内容类型获取表相关信息
func GetTableInfo(uploadType string) (tableName string, fieldName string, defaults []string) {
	r := strings.SplitN(uploadType, `.`, 2)
	switch len(r) {
	case 2:
		fieldName = r[1]
		fallthrough
	case 1:
		tableName = r[0]
		if r[0] != uploadType {
			defaults = append(defaults, tableName)
		}
	default:
		tableName = uploadType
	}
	return
}
