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
	"github.com/webx-top/com"
)

//TODO: 移出去

func (c *CLIConfig) FTPStopHistory() (err error) {
	if DefaultConfig == nil {
		return nil
	}
	return com.CloseProcessFromPidFile(DefaultConfig.FTP.PidFile)
}

func (c *CLIConfig) FTPStart(writer ...io.Writer) (err error) {
	err = c.FTPStopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	params := []string{os.Args[0], `--config`, c.Conf, `--type`, `ftpserver`}
	c.cmds["ftp"] = com.RunCmdWithWriter(params, writer...)
	return nil
}

func (c *CLIConfig) FTPStop() error {
	return c.CmdStop("ftp")
}

func (c *CLIConfig) FTPRestart(writer ...io.Writer) error {
	err := c.FTPStop()
	if err != nil {
		return err
	}
	return c.FTPStart(writer...)
}
