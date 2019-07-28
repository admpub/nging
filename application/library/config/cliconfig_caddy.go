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
	"io"
	"os"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/caddy"
	"github.com/webx-top/com"
)

func (c *CLIConfig) CaddyStopHistory() (err error) {
	if DefaultConfig == nil {
		return nil
	}
	return com.CloseProcessFromPidFile(DefaultConfig.Caddy.PidFile)
}

func (c *CLIConfig) CaddyStart(writer ...io.Writer) (err error) {
	err = c.CaddyStopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	params := []string{os.Args[0], `--config`, c.Conf, `--type`, `webserver`}
	if caddy.EnableReload {
		c.caddyCh = com.NewCmdChanReader()
		c.cmds["caddy"] = com.RunCmdWithReaderWriter(params, c.caddyCh, writer...)
	} else {
		c.cmds["caddy"] = com.RunCmdWithWriter(params, writer...)
	}
	return nil
}

func (c *CLIConfig) CaddyStop() error {
	defer func() {
		if c.caddyCh != nil {
			c.caddyCh.Close()
			c.caddyCh = nil
		}
	}()
	return c.CmdStop("caddy")
}

func (c *CLIConfig) CaddyReload() error {
	if c.caddyCh == nil || com.IsWindows { //windows不支持重载，需要重启
		return c.CaddyRestart()
	}
	c.caddyCh.Send(com.BreakLine)
	return nil
	//c.CmdSendSignal("caddy", os.Interrupt)
}

func (c *CLIConfig) CaddyRestart(writer ...io.Writer) error {
	err := c.CaddyStop()
	if err != nil {
		log.Error(err)
	}
	return c.CaddyStart(writer...)
}
