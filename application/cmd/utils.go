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
	stdLog "log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/cmd/bootconfig"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/license"
)

func Init() {
	config.Version.Number = echo.String(`VERSION`)
	config.Version.Label = echo.String(`LABEL`)
	config.Version.Package = echo.String(`PACKAGE`)
	license.SetVersion(config.Version.Number + `-` + config.Version.Label)
	license.SetPackage(config.Version.Package)
	if !bootconfig.Bindata {
		bootconfig.Develop = true
	}
	buildTime := echo.String(`BUILD_TIME`)
	if len(buildTime) == 0 {
		gitFile := filepath.Join(echo.Wd(), `.git/index`)
		f, err := os.Stat(gitFile)
		if err == nil {
			buildTime = f.ModTime().Format(`20060102150405`)
		}
		echo.Set(`BUILD_TIME`, buildTime)
	}
	config.Version.BuildTime = buildTime
	config.Version.BuildOS = echo.String(`BUILD_OS`)
	if len(config.Version.BuildOS) == 0 {
		config.Version.BuildOS = runtime.GOOS
	}
	config.Version.BuildArch = echo.String(`BUILD_ARCH`)
	if len(config.Version.BuildArch) == 0 {
		config.Version.BuildArch = runtime.GOARCH
	}
	stdLogWriter := log.Writer(log.LevelInfo)
	middleware.DefaultLogWriter = stdLogWriter
	stdLog.SetOutput(stdLogWriter)
	stdLog.SetFlags(stdLog.Lshortfile)
	bootconfig.MustLicensed = echo.Bool(`MUST_LICENSED`)
	config.Version.CommitID = echo.String(`COMMIT`)
	config.Version.DBSchema = echo.Float64(`SCHEMA_VER`)
	config.Version.Licensed = license.Ok(nil)
	if config.Version.Licensed {
		config.Version.Expired = license.License().Info.Expiration
	}
	rootCmd.Short = bootconfig.Short
	rootCmd.Long = bootconfig.Long
}

func Add(cmds ...*cobra.Command) {
	rootCmd.AddCommand(cmds...)
}
