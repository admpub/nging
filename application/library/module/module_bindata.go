//go:build bindata
// +build bindata

package module

import (
	"github.com/admpub/nging/v4/application/library/ntemplate"
	"github.com/webx-top/echo/middleware"
)

func (m *Module) applyTemplateAndAssets() {
}

func SetTemplate(pa ntemplate.PathAliases, key string, templatePath string) {
}

func SetAssets(so *middleware.StaticOptions, assetsPath string) {
}
