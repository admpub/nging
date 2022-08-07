package cmder

import (
	"github.com/admpub/log"
	"github.com/webx-top/db"

	"github.com/admpub/nging/v4/application/library/config"
	"github.com/nging-plugins/frpmanager/pkg/dbschema"
	"github.com/nging-plugins/frpmanager/pkg/library/utils"
)

func NewBase() *Base {
	return &Base{
		CLIConfig: config.FromCLI(),
	}
}

type Base struct {
	CLIConfig *config.CLIConfig
}

func (c *Base) getConfig() *config.Config {
	if config.FromFile() == nil {
		c.CLIConfig.ParseConfig()
	}
	return config.FromFile()
}

func (c *Base) RebuildConfigFile(data interface{}) error {
	return c.rebuildConfigFile(data, false)
}

func (c *Base) MustRebuildConfigFile(data interface{}) error {
	return c.rebuildConfigFile(data, true)
}

func (c *Base) rebuildConfigFile(data interface{}, must bool) (err error) {
	switch v := data.(type) {
	case string:
		switch v {
		case `frpserver`:
			md := &dbschema.NgingFrpServer{}
			_, err = md.ListByOffset(nil, nil, 0, -1, `disabled`, `N`)
			if err != nil {
				return
			}
			for _, row := range md.Objects() {
				err = c.rebuildConfigFile(row, must)
				if err != nil {
					return
				}
			}
			return
		case `frpclient`:
			md := &dbschema.NgingFrpClient{}
			_, err = md.ListByOffset(nil, nil, 0, -1, `disabled`, `N`)
			if err != nil {
				return
			}
			for _, row := range md.Objects() {
				err = c.rebuildConfigFile(row, must)
				if err != nil {
					return
				}
			}
			return
		default:
			return
		}
	default:
		err = utils.SaveConfigFile(data)
	}
	if err != nil {
		if db.ErrNoMoreRows == err {
			if must {
				err = ErrNoAvailibaleConfigFound
			} else {
				err = nil
				return
			}
		}
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
	return
}

func (c *Base) ConfigFile(id uint, isServer bool) string {
	return utils.ConfigFile(id, isServer)
}

func (c *Base) PidFile(id string, isServer bool) string {
	return utils.PidFile(id, isServer)
}

func (c *Base) Reload() error {
	return nil
}
