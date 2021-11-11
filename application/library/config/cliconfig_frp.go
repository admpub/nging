/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/com/encoding/json"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/library/frp"
)

//TODO: 移出去

var (
	EmptyBytes                 = []byte{}
	ErrNoAvailibaleConfigFound = errors.New(`no available configurations found`)
)

//FRP Server

func (c *CLIConfig) FRPStopHistory(ids ...string) (err error) {
	if len(ids) > 0 {
		for _, id := range ids {
			pidPath := c.FRPPidFile(id, true)
			err = com.CloseProcessFromPidFile(pidPath)
			if err != nil {
				log.Error(err.Error() + `: ` + pidPath)
			}
		}
		return nil
	}
	pidFilePath := filepath.Join(echo.Wd(), `data/pid/frp/server`)
	err = filepath.Walk(pidFilePath, func(pidPath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		err = com.CloseProcessFromPidFile(pidPath)
		if err != nil {
			log.Error(err.Error() + `: ` + pidPath)
		}
		return os.Remove(pidPath)
	})
	return
}

func (c *CLIConfig) FRPRebuildConfigFile(data interface{}, configFiles ...string) error {
	return c.rebuildFRPConfigFile(data, false, configFiles...)
}

func (c *CLIConfig) MustFRPRebuildConfigFile(data interface{}, configFiles ...string) error {
	return c.rebuildFRPConfigFile(data, true, configFiles...)
}

func (c *CLIConfig) rebuildFRPConfigFile(data interface{}, must bool, configFiles ...string) (err error) {
	var m interface{}
	var configFile string
	if len(configFiles) > 0 {
		configFile = configFiles[0]
	}
	switch v := data.(type) {
	case *dbschema.NgingFrpClient:
		if len(configFile) == 0 {
			configFile = c.FRPConfigFile(v.Id, false)
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
			configFile = c.FRPConfigFile(v.Id, true)
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
				err = c.rebuildFRPConfigFile(row, must, configFiles...)
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
				err = c.rebuildFRPConfigFile(row, must, configFiles...)
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
	if strings.HasSuffix(configFile, `.json`) {
		b, err = json.MarshalIndent(m, ``, `  `)
	} else {
		b, err = confl.Marshal(m)
	}
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

const FRPConfigExtension = `.json` //`.yaml`

func (c *CLIConfig) FRPConfigFile(id uint, isServer bool) string {
	configFile := `server`
	if !isServer {
		configFile = `client`
	}
	configFile = filepath.Join(echo.Wd(), `config`, `frp`, configFile)
	err := com.MkdirAll(configFile, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	return filepath.Join(configFile, fmt.Sprintf(`%d`, id)+FRPConfigExtension)
}

func (c *CLIConfig) FRPPidFile(id string, isServer bool) string {
	pidFile := `server`
	if !isServer {
		pidFile = `client`
	}
	pidFile = filepath.Join(echo.Wd(), `data/pid/frp`, pidFile)
	err := com.MkdirAll(pidFile, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	return filepath.Join(pidFile, id+`.pid`)
}

func (c *CLIConfig) FRPSaveConfigFile(data interface{}) (err error) {
	var configFile string
	switch v := data.(type) {
	case *dbschema.NgingFrpServer:
		configFile = c.FRPConfigFile(v.Id, true)
		if v.Disabled == `Y` {
			if !com.FileExists(configFile) {
				return nil
			}
			return os.Remove(configFile)
		}
		if len(v.Plugins) > 0 {
			serverConfigExtra := frp.NewServerConfigExtra()
			serverConfigExtra.PluginOptions = frp.ServerPluginOptions(strings.Split(v.Plugins, `,`)...)
			copied := *v
			if len(copied.Extra) > 0 {
				serverConfigExtra.Extra = []byte(copied.Extra)
			}
			copied.Extra = serverConfigExtra.String()
			data = copied
		}
	case *dbschema.NgingFrpClient:
		configFile = c.FRPConfigFile(v.Id, false)
		if v.Disabled == `Y` {
			if !com.FileExists(configFile) {
				return nil
			}
			return os.Remove(configFile)
		}
		cfg := frp.NewClientConfig()
		cfg.NgingFrpClient = v
		cfg.Extra, err = frp.Table2Config(v)
		if err != nil {
			return
		}
		cfg.NgingFrpClient.Extra = ``
		data = cfg
	default:
		return fmt.Errorf(`unsupport save config: %T`, v)
	}
	var b []byte
	if strings.HasSuffix(configFile, `.json`) {
		b, err = json.MarshalIndent(data, ``, `  `)
	} else {
		b, err = confl.Marshal(data)
	}
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, b, os.ModePerm)
}

func (c *CLIConfig) FRPStartID(id uint, writer ...io.Writer) (err error) {
	configFile := c.FRPConfigFile(id, true)
	params := []string{os.Args[0], `--config`, c.Conf, `--subconfig`, configFile, `--type`, `frpserver`}
	cmd, iErr, _, rErr := com.RunCmdWithWriterx(params, time.Millisecond*500, writer...)
	key := fmt.Sprintf("frpserver.%d", id)
	c.cmds[key] = cmd
	if iErr != nil {
		err = fmt.Errorf(iErr.Error()+`: %s`, rErr.Buffer().String())
		return
	}
	return
}

func (c *CLIConfig) FRPStart(writer ...io.Writer) (err error) {
	err = c.FRPStopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	md := &dbschema.NgingFrpServer{}
	cd := db.And(
		db.Cond{`disabled`: `N`},
	)
	_, err = md.ListByOffset(nil, nil, 0, -1, cd)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return
	}
	for _, row := range md.Objects() {
		err := c.FRPStartID(row.Id, writer...)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (c *CLIConfig) FRPStop() error {
	return c.CmdGroupStop("frpserver")
}

func (c *CLIConfig) FRPRestart(writer ...io.Writer) error {
	err := c.FRPStop()
	if err != nil {
		return err
	}
	return c.FRPStart(writer...)
}

func (c *CLIConfig) FRPRestartID(id string, writer ...io.Writer) error {
	err := c.FRPStopID(id)
	if err != nil {
		return err
	}
	idv, _ := strconv.ParseUint(id, 10, 32)
	return c.FRPStartID(uint(idv), writer...)
}

func (c *CLIConfig) FRPStopID(id string) error {
	err := c.CmdStop("frpserver." + id)
	if err != nil && !strings.Contains(err.Error(), `finished`) {
		return err
	}
	pidPath := c.FRPPidFile(id, true)
	err = com.CloseProcessFromPidFile(pidPath)
	if err != nil && !strings.Contains(err.Error(), `finished`) {
		log.Error(err.Error() + `: ` + pidPath)
	}
	return nil
}
