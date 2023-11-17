//go:build embedNgingPluginTemplate

package dockermanager

import (
	"embed"
)

//go:embed template
var TemplateFS embed.FS

//go:embed public
var AssetsFS embed.FS
