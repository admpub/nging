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

package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/cmd/event"
	"github.com/admpub/nging/v3/application/handler/setup"
	"github.com/admpub/nging/v3/application/library/config"
	"github.com/admpub/nging/v3/application/library/service"
)

// 将Nging安装为系统服务的工具

// ServiceOptions 服务选项
var ServiceOptions = &service.Options{
	Name:        ``,
	DisplayName: ``,
	Description: ``,
}

var serviceCmd = &cobra.Command{
	Use:     "service",
	Short:   "Running as a service on major platforms.",
	Example: filepath.Base(os.Args[0]) + " service [install|uninstall|start|restart|stop]",
	RunE:    serviceRunE,
}

func serviceRunE(cmd *cobra.Command, args []string) error {
	conf, err := config.InitConfig()
	config.MustOK(err)
	conf.AsDefault()
	//application.DefaultConfigWatcher(false)
	if len(args) < 1 {
		return cmd.Usage()
	}
	if len(ServiceOptions.Name) == 0 {
		ServiceOptions.Name = event.SoftwareName
	}
	if len(ServiceOptions.DisplayName) == 0 {
		ServiceOptions.DisplayName = ServiceOptions.Name
	}
	if len(ServiceOptions.Name) == 0 {
		ServiceOptions.Description = ServiceOptions.DisplayName + ` Service`
	}

	if config.IsInstalled() {
		if err := setup.Upgrade(); err != nil && os.ErrNotExist != err {
			log.Error(`upgrade.sql: `, err)
		}
	}

	return service.Run(ServiceOptions, args[0])
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
