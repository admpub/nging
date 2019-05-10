// +build lazygit

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
	"fmt"
	"log"
	"os"

	"github.com/admpub/nging/application/library/config"
	"github.com/jesseduffield/lazygit/pkg/app"
	cfg "github.com/jesseduffield/lazygit/pkg/config"
	"github.com/spf13/cobra"
)

// GIT 高效命令行图形工具
// 快捷键参考：
// https://github.com/jesseduffield/lazygit/blob/master/docs/Keybindings.md

var (
	layzyGITConfigFlag    *bool
	layzyGITDebuggingFlag *bool
)

var layzyGITCmd = &cobra.Command{
	Use:  "git",
	RunE: layzyGITRunE,
}

func layzyGITRunE(cmd *cobra.Command, args []string) error {

	if *layzyGITConfigFlag {
		fmt.Printf("%s\n", cfg.GetDefaultConfig())
		os.Exit(0)
	}
	appConfig, err := cfg.NewAppConfig("lazygit", config.VersionNumber(), config.CommitID(), config.BuildTime(), `nging`, layzyGITDebuggingFlag)
	if err != nil {
		log.Fatal(err.Error())
	}

	app, err := app.Setup(appConfig)
	if err != nil {
		app.Log.Error(err.Error())
		log.Fatal(err.Error())
	}

	app.Gui.RunWithSubprocesses()
	return nil
}

func init() {
	rootCmd.AddCommand(layzyGITCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	layzyGITConfigFlag = layzyGITCmd.Flags().BoolP("config", "c", false, "Print the current default config")
	layzyGITDebuggingFlag = layzyGITCmd.Flags().BoolP("debug", "d", false, "a boolean")
}
