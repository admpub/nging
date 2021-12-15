//go:build embedNgingPluginTemplate

package ftpmanager

import (
	"embed"
)

//go:embed template
var TemplateFS embed.FS
