//go:build __misc__

package template

import (
	"embed"
)

//go:embed backend
var template embed.FS
