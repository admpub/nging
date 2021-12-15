//go:build embedNgingPluginTemplate

package sshmanager

import (
	"embed"
)

//go:embed template
var TemplateFS embed.FS
