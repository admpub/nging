/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package service

import (
	"io"

	"github.com/admpub/service"
)

var DefaultMaxRetries = 10
var DefaultRetryInterval = 60 //60s

type Options struct {
	Name          string // Required name of the service. No spaces suggested.
	DisplayName   string // Display name, spaces allowed.
	Description   string // Long description of service.
	Options       map[string]interface{}
	MaxRetries    int // 最大重试次数
	RetryInterval int // 重试间隔（秒）
}

// Config is the runner app config structure.
type Config struct {
	service.Config
	logger service.Logger

	Dir           string
	Exec          string
	Args          []string
	Env           []string
	MaxRetries    int
	RetryInterval int // 重试间隔（秒）

	OnExited       func() error `json:"-"`
	Stderr, Stdout io.Writer    `json:"-"`
}

func (c *Config) CopyFromOptions(options *Options) *Config {
	c.Name = options.Name
	c.DisplayName = options.DisplayName
	c.Description = options.Description
	c.Option = service.KeyValue{}
	for k, v := range c.DefaultOptions() {
		c.Option[k] = v
	}
	if options.Options != nil {
		for k, v := range options.Options {
			c.Option[k] = v
		}
	}
	c.MaxRetries = options.MaxRetries
	c.RetryInterval = options.RetryInterval
	return c
}

func (c *Config) DefaultOptions() service.KeyValue {
	return map[string]interface{}{
		//  * POSIX
		//`PIDFile`: pidFile,
		`Restart`: `always`,

		//  * OS X
		`RunAtLoad`: true,
		`KeepAlive`: true,
		//`UserService`: true, //Install as a current user service.
		//`SessionCreate`: true, //Create a full user session.

		//  * Windows
		`OnFailure`:              `restart`,
		`OnFailureDelayDuration`: `5s`,
		`OnFailureResetPeriod`:   10,
	}
}
