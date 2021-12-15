//go:build embedNgingPluginTemplate

package ddnsmanager

import (
	"embed"
)

//go:embed template
var TemplateFS embed.FS
