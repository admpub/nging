package license

import (
	"github.com/webx-top/echo/middleware/tplfunc"
)

func init() {
	tplfunc.TplFuncMap[`HasFeature`] = HasFeature
	tplfunc.TplFuncMap[`HasAnyFeature`] = HasAnyFeature
}
