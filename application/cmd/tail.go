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
	"log"

	ninetail "github.com/admpub/9t"
	"github.com/spf13/cobra"
)

// 多文件Tail动态

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "multi-file tailer (like `tail -f a.log b.log ...`)",
	RunE:  tailRunE,
}

func tailRunE(cmd *cobra.Command, filenames []string) error {
	runner, err := ninetail.Runner(filenames, ninetail.Config{Colorize: true})
	if err != nil {
		log.Fatal(err)
	}
	runner.Run()
	return nil
}

func init() {
	rootCmd.AddCommand(tailCmd)
}
