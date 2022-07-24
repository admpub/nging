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

	"github.com/spf13/cobra"
	"github.com/webx-top/echo/defaults"

	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/handler/setup"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/config/subconfig/sdb"
)

// 静默安装
var initDBConfig = &sdb.DB{
	Type:     `mysql`, // mysql / sqlite
	User:     `root`,
	Database: `nging`,
	Host:     `127.0.0.1:3306`,
}

var initInstallConfig = &struct {
	Charset    string
	AdminUser  string
	AdminPass  string
	AdminEmail string
}{
	Charset:   sdb.MySQLDefaultCharset,
	AdminUser: `admin`,
}

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Silent install",
	Example: filepath.Base(os.Args[0]) + " init",
	RunE:    initRunE,
}

func initRunE(cmd *cobra.Command, args []string) error {
	conf, err := config.InitConfig()
	if err != nil {
		return err
	}
	conf.AsDefault()
	ctx := defaults.NewMockContext()
	ctx.Request().Form().Set(`type`, initDBConfig.Type)
	ctx.Request().Form().Set(`user`, initDBConfig.User)
	ctx.Request().Form().Set(`host`, initDBConfig.Host)
	ctx.Request().Form().Set(`password`, initDBConfig.Password)
	ctx.Request().Form().Set(`database`, initDBConfig.Database)
	ctx.Request().Form().Set(`prefix`, initDBConfig.Prefix)
	ctx.Request().Form().Set(`charset`, initInstallConfig.Charset)
	ctx.Request().Form().Set(`adminUser`, initInstallConfig.AdminUser)
	ctx.Request().Form().Set(`adminPass`, initInstallConfig.AdminPass)
	ctx.Request().Form().Set(`adminEmail`, initInstallConfig.AdminEmail)
	err = setup.Setup(ctx)
	if err == nil {
		log.Okay(`Congratulations, this program has been installed successfully`)
	}
	return err
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&initDBConfig.Type, "type", initDBConfig.Type, "database type")
	initCmd.Flags().StringVar(&initDBConfig.User, "user", initDBConfig.User, "database user")
	initCmd.Flags().StringVar(&initDBConfig.Host, "host", initDBConfig.Host, "database host")
	initCmd.Flags().StringVar(&initDBConfig.Password, "password", initDBConfig.Password, "database password")
	initCmd.Flags().StringVar(&initDBConfig.Database, "database", initDBConfig.Database, "database name")
	initCmd.Flags().StringVar(&initDBConfig.Prefix, "prefix", initDBConfig.Prefix, "database table prefix")
	initCmd.Flags().StringVar(&initInstallConfig.Charset, "charset", initInstallConfig.Charset, "database table charset")
	initCmd.Flags().StringVar(&initInstallConfig.AdminUser, "adminUser", initInstallConfig.AdminUser, "administrator name")
	initCmd.Flags().StringVar(&initInstallConfig.AdminPass, "adminPass", initInstallConfig.AdminPass, "administrator password")
	initCmd.Flags().StringVar(&initInstallConfig.AdminEmail, "adminEmail", initInstallConfig.AdminEmail, "administrator e-mail")
}
