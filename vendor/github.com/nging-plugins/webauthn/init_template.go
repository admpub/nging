//go:build embedNgingPluginTemplate

package webauthn

import (
	"embed"
)

//go:embed template
var TemplateFS embed.FS
