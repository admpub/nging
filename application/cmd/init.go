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
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/middleware/language"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/cmd/bootconfig"
	"github.com/admpub/nging/v5/application/handler/setup"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/subconfig/sdb"
	"github.com/admpub/nging/v5/application/request"
	"github.com/admpub/once"
)

// 静默安装
var InitDBConfig = &sdb.DB{
	Type:     `mysql`, // mysql / sqlite
	User:     `root`,
	Database: `nging`,
	Host:     `127.0.0.1:3306`,
}

var InitInstallConfig = &struct {
	Charset        string
	AdminUser      string
	AdminPass      string
	AdminEmail     string
	Language       string // en / zh-cn
	AdminPwdRandom bool
}{
	Charset:   sdb.MySQLDefaultCharset,
	AdminUser: `admin`,
	Language:  `zh-cn`,
}

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Silent install",
	Example: filepath.Base(os.Args[0]) + " init [options]",
	RunE:    initRunE,
}

var translate *language.Translate
var translock once.Once

func initTranslate() {
	translate = BuildTranslator(config.FromFile().Language)
}

func GetTranslator() *language.Translate {
	translock.Do(initTranslate)
	return translate
}

func ResetTranslator() {
	translock.Reset()
}

func BuildTranslator(c language.Config) *language.Translate {
	c.SetFSFunc(bootconfig.LangFSFunc)
	i18n := language.NewI18n(&c)
	tr := &language.Translate{}
	tr.Reset(InitInstallConfig.Language, i18n)
	return tr
}

func NewContext() echo.Context {
	ctx := defaults.NewMockContext().SetAuto(true)

	// 启用多语言支持
	ctx.SetTranslator(GetTranslator())

	ctx.Request().Header().Set(echo.HeaderAccept, echo.MIMETextPlain)
	return ctx
}

func initRunE(cmd *cobra.Command, args []string) error {
	conf, err := config.InitConfig()
	if err != nil {
		return err
	}
	conf.AsDefault()
	ctx := NewContext()
	ctx.Request().SetMethod(echo.POST)
	req := &request.Setup{
		Type:       InitDBConfig.Type,
		User:       InitDBConfig.User,
		Password:   InitDBConfig.Password,
		Host:       InitDBConfig.Host,
		Database:   InitDBConfig.Database,
		Charset:    InitInstallConfig.Charset,
		AdminUser:  InitInstallConfig.AdminUser,
		AdminPass:  InitInstallConfig.AdminPass,
		AdminEmail: InitInstallConfig.AdminEmail,
	}
	if InitInstallConfig.AdminPwdRandom {
		req.AdminPass = com.RandomAlphanumeric(16)
	}
	err = echo.ValidateStruct(ctx, req)
	if err != nil {
		return err
	}
	req.SetFromCLI(true)
	ctx.Internal().Set(`validated`, req)
	err = setup.Setup(ctx)
	if err == nil {
		log.Okay(ctx.T(`Congratulations, this program has been installed successfully`))
		if InitInstallConfig.AdminPwdRandom {
			pwdFile := filepath.Join(echo.Wd(), `nging-password.txt`)
			err := os.WriteFile(pwdFile, []byte(req.AdminPass), 0660)
			if err != nil {
				log.Errorf(ctx.T(`failed to record the generated password to the file %q: %v`), pwdFile, err.Error())
			} else {
				log.Okay(ctx.T(`The generated password has been recorded to the file %q. Please review the contents of this file and remember it and delete it.`, pwdFile))
			}
		}
	}
	return err
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&InitDBConfig.Type, "type", InitDBConfig.Type, "database type")
	initCmd.Flags().StringVar(&InitDBConfig.User, "user", InitDBConfig.User, "database user")
	initCmd.Flags().StringVar(&InitDBConfig.Host, "host", InitDBConfig.Host, "database host")
	initCmd.Flags().StringVar(&InitDBConfig.Password, "password", InitDBConfig.Password, "database password")
	initCmd.Flags().StringVar(&InitDBConfig.Database, "database", InitDBConfig.Database, "database name")
	//initCmd.Flags().StringVar(&InitDBConfig.Prefix, "prefix", InitDBConfig.Prefix, "database table prefix")
	initCmd.Flags().StringVar(&InitInstallConfig.Charset, "charset", InitInstallConfig.Charset, "database table charset")
	initCmd.Flags().StringVar(&InitInstallConfig.AdminUser, "adminUser", InitInstallConfig.AdminUser, "administrator name")
	initCmd.Flags().StringVar(&InitInstallConfig.AdminPass, "adminPass", InitInstallConfig.AdminPass, "administrator password")
	initCmd.Flags().StringVar(&InitInstallConfig.AdminEmail, "adminEmail", InitInstallConfig.AdminEmail, "administrator e-mail")
	initCmd.Flags().BoolVar(&InitInstallConfig.AdminPwdRandom, "adminPwdRandom", InitInstallConfig.AdminPwdRandom, "generate administrator random password")
}
