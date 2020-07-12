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
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/webx-top/echo"

	"github.com/spf13/cobra"
)

// go run main.go testsuite

var testSuiteCmd = &cobra.Command{
	Use:   "testsuite",
	Short: "",
	RunE:  testSuiteRunE,
}

func testSuiteRunE(cmd *cobra.Command, filenames []string) error {
	err := config.ParseConfig()
	if err != nil {
		panic(err)
	}
	row, err := common.SQLQuery().GetRow("SELECT * FROM nging_user WHERE id > 0")
	if err != nil {
		panic(err)
	}
	echo.Dump(row.Timestamp(`created`).Format(`2006-01-02 15:04:05`))
	return nil
}

func init() {
	rootCmd.AddCommand(testSuiteCmd)
}
