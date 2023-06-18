//go:build embedNgingPluginTemplate

package firewallmanager

import (
	"embed"
)

//go:embed template
var TemplateFS embed.FS
