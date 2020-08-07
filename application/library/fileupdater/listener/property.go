package listener

import (
	"github.com/admpub/nging/application/library/fileupdater"
)

var (
	// GenUpdater 生成Updater
	GenUpdater      = fileupdater.GenUpdater
	NewProperty     = fileupdater.NewProperty
	NewPropertyWith = fileupdater.NewPropertyWith
)

// Property 附加属性
type Property = fileupdater.Property
