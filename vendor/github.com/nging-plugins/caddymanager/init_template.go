//go:build embedNgingPluginTemplate

package caddymanager

import (
	"embed"
)

//go:embed template
var TemplateFS embed.FS
