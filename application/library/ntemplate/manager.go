package ntemplate

import (
	"github.com/webx-top/echo/middleware/render/driver"
)

func NewManager(mgr driver.Manager, pa PathAliases) driver.Manager {
	return &manager{
		Manager: mgr,
		pa:      pa,
	}
}

type manager struct {
	driver.Manager
	pa PathAliases
}

func (m *manager) AddCallback(rootDir string, callback func(name, typ, event string)) {
	originalCb := callback
	callback = func(name, typ, event string) {
		name = m.pa.RestorePrefix(name)
		originalCb(name, typ, event)
	}
	m.Manager.AddCallback(rootDir, callback)
}
