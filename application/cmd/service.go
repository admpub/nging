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

package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/cmd/bootconfig"
	"github.com/admpub/nging/v5/application/handler/setup"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/service"
)

// 将Nging安装为系统服务的工具

// ServiceOptions 服务选项
var ServiceOptions = &service.Options{
	Name:          ``,
	DisplayName:   ``,
	Description:   ``,
	Options:       map[string]interface{}{},
	MaxRetries:    10,
	RetryInterval: 60,
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
	if len(args) < 1 {
		return cmd.Usage()
	}
	if len(ServiceOptions.Name) == 0 {
		ServiceOptions.Name = bootconfig.SoftwareName
	}
	if len(ServiceOptions.DisplayName) == 0 {
		ServiceOptions.DisplayName = strings.Title(ServiceOptions.Name)
	}
	if len(ServiceOptions.Name) == 0 {
		ServiceOptions.Description = ServiceOptions.DisplayName + ` Service`
	}
	if config.FromFile() != nil && config.FromFile().Extend != nil {
		systemService := config.FromFile().Extend.Children(`systemService`)
		maxRetries := systemService.Int(`maxRetries`)
		if maxRetries > 0 {
			ServiceOptions.MaxRetries = maxRetries
		}
		retryInterval := systemService.Int(`retryInterval`)
		if retryInterval > 0 {
			ServiceOptions.RetryInterval = retryInterval
		}
		systemServiceOptions := systemService.Children(`options`)
		for k, v := range systemServiceOptions {
			ServiceOptions.Options[k] = v
		}
	}

	if config.IsInstalled() && bootconfig.AutoUpgradeDBStruct && !conf.Sys.DisableAutoUpgradeDB {
		if err := setup.Upgrade(); err != nil && os.ErrNotExist != err {
			log.Error(`upgrade.sql: `, err)
		}
	}

	return service.Run(ServiceOptions, args[0])
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
