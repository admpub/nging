//go:build !bindata
// +build !bindata

package module

import (
	"strings"

	"github.com/admpub/nging/v4/application/library/bindata"
	"github.com/admpub/nging/v4/application/library/ntemplate"
	"github.com/webx-top/echo/middleware"
)

func (m *Module) applyTemplateAndAssets() {
	m.setTemplate(bindata.PathAliases)
	m.setAssets(bindata.StaticOptions)
}

func SetTemplate(pa ntemplate.PathAliases, key string, templatePath string) {
	if len(templatePath) == 0 {
		return
	}
	if templatePath[0] != '.' && templatePath[0] != '/' && !strings.HasPrefix(templatePath, `vendor/`) {
		templatePath = NgingPluginDir + `/` + templatePath
	}
	pa.Add(key, templatePath)
}

func SetAssets(so *middleware.StaticOptions, assetsPath string) {
	if len(assetsPath) == 0 {
		return
	}
	if assetsPath[0] != '.' && assetsPath[0] != '/' && !strings.HasPrefix(assetsPath, `vendor/`) {
		assetsPath = NgingPluginDir + `/` + assetsPath
	}
	so.AddFallback(assetsPath)
}
