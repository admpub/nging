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

	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/service"
	"github.com/spf13/cobra"
)

// 将Nging安装为系统服务的工具

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
	return service.Run(args[0])
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
