//go:build bindata
// +build bindata

package cloud

import (
	"net/http"

	"github.com/coscms/webcore/library/bindata"
)

func init() {
	RegisterFileSource(`assetfs`, `读取本程序内嵌的静态资源文件(例如：assetfs:public/assets/backend 或 assetfs:public/assets/frontend)`, func() http.FileSystem {
		return bindata.StaticAssetFS
	})
}
