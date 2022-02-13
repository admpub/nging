//go:build bindata
// +build bindata

package module

import (
	//"github.com/admpub/nging/v4/application/library/bindata"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/route"
	"github.com/admpub/nging/v4/application/registry/dashboard"
	"github.com/admpub/nging/v4/application/registry/navigate"
)

func SetTemplate(pa ntemplate.PathAliases, key string, templatePath string) {
}

func SetAssets(so *middleware.StaticOptions, assetsPath string) {
}
