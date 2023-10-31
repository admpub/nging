package enginevent

import (
	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/nging-plugins/caddymanager/application/library/form"
)

type OnVhostConfigSaved interface {
	OnVhostConfigSaved(id uint, formData *form.Values) error
}

type OnVhostConfigSaving interface {
	OnVhostConfigSaving(id uint, formData *form.Values) error
}

func FireVhostConfigSaving(cfg engine.Configer, id uint, formData *form.Values) error {
	if sv, ok := cfg.(OnVhostConfigSaving); ok {
		return sv.OnVhostConfigSaving(id, formData)
	}
	return nil
}

func FireVhostConfigSaved(cfg engine.Configer, id uint, formData *form.Values) error {
	if sv, ok := cfg.(OnVhostConfigSaved); ok {
		return sv.OnVhostConfigSaved(id, formData)
	}
	return nil
}
