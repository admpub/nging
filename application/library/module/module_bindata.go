//go:build bindata
// +build bindata

package module

import (
	"github.com/admpub/nging/v4/application/library/bindata"
	"github.com/admpub/nging/v4/application/library/ntemplate"
	"github.com/webx-top/echo/middleware"
)

func (m *Module) applyTemplateAndAssets() {
}

func SetBackendTemplate(key string, templatePath string) {
	SetTemplate(bindata.PathAliases, key, templatePath)
}

func SetBackendAssets(assetsPath string) {
	SetAssets(bindata.StaticOptions, assetsPath)
}
