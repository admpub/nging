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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

//TODO: 移出去

//FRP Client

func (c *CLIConfig) FRPClientStopHistory(ids ...string) (err error) {
	if len(ids) > 0 {
		for _, id := range ids {
			pidPath := c.FRPPidFile(id, false)
			err = com.CloseProcessFromPidFile(pidPath)
			if err != nil {
				log.Error(err.Error() + `: ` + pidPath)
			}
		}
		return nil
	}
	pidFilePath := filepath.Join(echo.Wd(), `data/pid/frp/client`)
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

func (c *CLIConfig) FRPClientStartID(id uint, writer ...io.Writer) (err error) {
	configFile := c.FRPConfigFile(id, false)
	params := []string{os.Args[0], `--config`, c.Conf, `--subconfig`, configFile, `--type`, `frpclient`}
	cmd, iErr, _, rErr := com.RunCmdWithWriterx(params, time.Millisecond*500, writer...)
	key := fmt.Sprintf("frpclient.%d", id)
	c.cmds[key] = cmd
	if iErr != nil {
		err = fmt.Errorf(iErr.Error()+`: %s`, rErr.Buffer().String())
		return
	}
	return
}

func (c *CLIConfig) FRPClientStart(writer ...io.Writer) (err error) {
	err = c.FRPClientStopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	md := &dbschema.NgingFrpClient{}
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
		err := c.FRPClientStartID(row.Id, writer...)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (c *CLIConfig) FRPClientStop() error {
	return c.CmdGroupStop("frpclient")
}

func (c *CLIConfig) FRPClientRestart(writer ...io.Writer) error {
	err := c.FRPClientStop()
	if err != nil {
		return err
	}
	return c.FRPClientStart(writer...)
}

func (c *CLIConfig) FRPClientRestartID(id string, writer ...io.Writer) error {
	err := c.FRPClientStopID(id)
	if err != nil {
		return err
	}
	idv, _ := strconv.ParseUint(id, 10, 32)
	return c.FRPClientStartID(uint(idv), writer...)
}

func (c *CLIConfig) FRPClientStopID(id string) error {
	err := c.CmdStop("frpclient." + id)
	if err != nil && !strings.Contains(err.Error(), `finished`) {
		return err
	}
	pidPath := c.FRPPidFile(id, false)
	err = com.CloseProcessFromPidFile(pidPath)
	if err != nil && !strings.Contains(err.Error(), `finished`) {
		log.Error(err.Error() + `: ` + pidPath)
	}
	return nil
}
