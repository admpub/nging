//go:build embedNgingPluginTemplate

package frpmanager

import (
	"embed"
)

//go:embed template
var TemplateFS embed.FS
