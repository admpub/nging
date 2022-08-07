package buildinfo

import (
	"path/filepath"

	"github.com/admpub/nging/v4/application/initialize/backend"
	"github.com/admpub/nging/v4/application/library/module"
	"github.com/webx-top/echo/middleware/render/driver"
)

func SetNgingDir(ngingDir string) {
	backend.AssetsDir = filepath.Join(ngingDir, backend.DefaultAssetsDir)
	backend.TemplateDir = filepath.Join(ngingDir, backend.DefaultTemplateDir)
}

func SetNgingPluginsDir(ngingPluginsDir string) {
	module.NgingPluginDir = ngingPluginsDir
}

func WatchTemplateDir(templateDirs ...string) {
	rendererDo := backend.RendererDo
	backend.RendererDo = func(renderer driver.Driver) {
		rendererDo(renderer)
		for _, templateDir := range templateDirs {
			renderer.Manager().AddWatchDir(templateDir)
		}
	}
}
