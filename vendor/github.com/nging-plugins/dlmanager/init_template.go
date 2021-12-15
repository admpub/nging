//go:build embedNgingPluginTemplate

package dlmanager

import (
	"embed"
)

//go:embed template
var TemplateFS embed.FS
