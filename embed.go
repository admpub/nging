//go:build misc

package nging

import (
	"embed"
)

//go:embed public/assets
var assets embed.FS

//go:embed template
var template embed.FS
