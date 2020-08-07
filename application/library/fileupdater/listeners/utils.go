package listeners

import (
	"strings"

	"github.com/admpub/nging/application/library/fileupdater"
	"github.com/admpub/nging/application/library/fileupdater/listener"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
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

func GenDefaultCallback(fieldName string) fileupdater.CallbackFunc {
	return func(m factory.Model) (tableID string, content string, property *listener.Property) {
		row := m.AsRow()
		tableID = row.String(`id`, `-1`)
		content = row.String(fieldName)
		property = listener.NewPropertyWith(m, db.Cond{`id`: row.Get(`id`, `-1`)})
		return
	}
}
