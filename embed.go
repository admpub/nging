//go:build __misc__

package nging

import (
	"embed"
)

//go:embed public/assets
var assets embed.FS

//go:embed template
var template embed.FS
