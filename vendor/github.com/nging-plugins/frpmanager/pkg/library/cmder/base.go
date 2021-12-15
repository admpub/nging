package cmder

import (
	"io/ioutil"
	"os"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/webx-top/db"

	"github.com/admpub/nging/v4/application/library/config"
	"github.com/nging-plugins/frpmanager/pkg/dbschema"
	"github.com/nging-plugins/frpmanager/pkg/library/frp"
	"github.com/nging-plugins/frpmanager/pkg/library/utils"
)

func NewBase() *Base {
	return &Base{
		CLIConfig: config.DefaultCLIConfig,
	}
}

type Base struct {
	CLIConfig *config.CLIConfig
}

func (c *Base) getConfig() *config.Config {
	if config.DefaultConfig == nil {
		c.CLIConfig.ParseConfig()
	}
	return config.DefaultConfig
}

func (c *Base) RebuildConfigFile(data interface{}, configFiles ...string) error {
	return c.rebuildConfigFile(data, false, configFiles...)
}

func (c *Base) MustRebuildConfigFile(data interface{}, configFiles ...string) error {
	return c.rebuildConfigFile(data, true, configFiles...)
}

func (c *Base) rebuildConfigFile(data interface{}, must bool, configFiles ...string) (err error) {
	var m interface{}
	var configFile string
	if len(configFiles) > 0 {
		configFile = configFiles[0]
	}
	switch v := data.(type) {
	case *dbschema.NgingFrpClient:
		if len(configFile) == 0 {
			configFile = c.ConfigFile(v.Id, false)
		}
		cfg := frp.NewClientConfig()
		cfg.NgingFrpClient = v
		cfg.Extra, err = frp.Table2Config(v)
		if err != nil {
			return
		}
		cfg.NgingFrpClient.Extra = ``
		m = cfg
	case *dbschema.NgingFrpServer:
		if len(configFile) == 0 {
			configFile = c.ConfigFile(v.Id, true)
		}
		m = v
	case string:
		switch v {
		case `frpserver`:
			md := &dbschema.NgingFrpServer{}
			_, err = md.ListByOffset(nil, nil, 0, -1, `disabled`, `N`)
			if err != nil {
				return
			}
			for _, row := range md.Objects() {
				err = c.rebuildConfigFile(row, must, configFiles...)
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
				err = c.rebuildConfigFile(row, must, configFiles...)
				if err != nil {
					return
				}
			}
			return
		default:
			return
		}
	default:
		return
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
	var b []byte
	b, err = confl.Marshal(m)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = ioutil.WriteFile(configFile, b, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		return
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
