package cmder

import (
	"strings"

	"github.com/admpub/nging/v4/application/initialize/backend"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
)

func onServerConfigChange(file string) error {
	id := config.DefaultCLIConfig.GenerateIDFromConfigFileName(file, true)
	if len(id) == 0 {
		return common.ErrIgnoreConfigChange
	}
	if !config.DefaultCLIConfig.IsRunning(`frpserver.` + id) {
		return common.ErrIgnoreConfigChange
	}
	if cm, err := GetServer(); err == nil {
		return cm.RestartBy(id)
	}
	return nil
}

func onClientConfigChange(file string) error {
	id := config.DefaultCLIConfig.GenerateIDFromConfigFileName(file, true)
	if len(id) == 0 {
		return common.ErrIgnoreConfigChange
	}
	if !config.DefaultCLIConfig.IsRunning(`frpclient.` + id) {
		return common.ErrIgnoreConfigChange
	}
	if cm, err := GetClient(); err == nil {
		return cm.RestartBy(id)
	}
	return nil
}

func init() {
	backend.OnConfigChange(func(filePath string) (err error) {
		if strings.Contains(filePath, `/frp/server/`) {
			err = onServerConfigChange(filePath)
		}
		return
	})
	backend.OnConfigChange(func(filePath string) (err error) {
		if strings.Contains(filePath, `/frp/client/`) {
			err = onClientConfigChange(filePath)
		}
		return
	})
}
