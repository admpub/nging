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

type Options struct {
	Name        string // Required name of the service. No spaces suggested.
	DisplayName string // Display name, spaces allowed.
	Description string // Long description of service.
}

// Config is the runner app config structure.
type Config struct {
	service.Config
	logger service.Logger

	Dir  string
	Exec string
	Args []string
	Env  []string

	OnExited       func() error `json:"-"`
	Stderr, Stdout io.Writer    `json:"-"`
}

func (c *Config) CopyFromOptions(options *Options) *Config {
	c.Name = options.Name
	c.DisplayName = options.DisplayName
	c.Description = options.Description
	return c
}
