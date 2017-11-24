package collector

import (
	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
)

var _ = dbschema.Base{}
var _ = log.SetLevel

/**
案例1：列表页面(支持分页) -> 获取内容页面链接 -> 采集页面内容(支持分页)
案例2：列表页面(支持分页) -> 获取封面页面链接 -> 采集封面页面信息(封面标题、图片、内容页面链接等)(支持分页) -> 采集链接页面内容(支持分页)
*/