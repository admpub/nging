package request

import "github.com/docker/docker/api/types/volume"

type VolumeAdd struct {
	Driver     string            // `validate:"required"`
	DriverOpts map[string]string `form_decoder:"splitKVRows" form_encoder:"joinKVRows"`
	Labels     map[string]string `form_decoder:"splitKVRows" form_encoder:"joinKVRows"`
	Name       string
}

func (v *VolumeAdd) Options() volume.CreateOptions {
	return volume.CreateOptions{
		Driver:     v.Driver,
		DriverOpts: v.DriverOpts,
		Labels:     v.Labels,
		Name:       v.Name,
	}
}
