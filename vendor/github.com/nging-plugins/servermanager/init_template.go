//go:build embedNgingPluginTemplate

package servermanager

import (
	"embed"
)

//go:embed template
var TemplateFS embed.FS
