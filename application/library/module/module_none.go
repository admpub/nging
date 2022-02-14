//go:build !bindata
// +build !bindata

package module

import (
	"github.com/admpub/nging/v4/application/library/bindata"
)

func (m *Module) applyTemplateAndAssets() {
	m.setTemplate(bindata.PathAliases)
	m.setAssets(bindata.StaticOptions)
}

func SetBackendTemplate(key string, templatePath string) {
	SetTemplate(bindata.PathAliases, key, templatePath)
}

func SetBackendAssets(assetsPath string) {
	SetAssets(bindata.StaticOptions, assetsPath)
}
