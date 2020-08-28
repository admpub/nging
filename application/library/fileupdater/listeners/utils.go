package listeners

import (
	"strings"
)

// parseSameFieldInfo 解析相似字段信息
// example parseSameFieldInfo(`image_original(image;image2)`)
func parseSameFieldInfo(fieldName string) (field string, sameFields []string) {
	if strings.HasSuffix(fieldName, `)`) {
		same := fieldName[0 : len(fieldName)-1]
		arr := strings.SplitN(same, `(`, 2)
		if len(arr) == 2 {
			fieldName = arr[0]
			same = arr[1]
			if len(same) > 0 {
				sameFields = strings.Split(same, ";")
			}
		}
	}
	field = fieldName
	return
}
