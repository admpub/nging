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
	"errors"
	"fmt"

	"github.com/admpub/nging/v3/application/library/config"

	"github.com/spf13/cobra"
)

// go run main.go testsuite

var testSuiteCmd = &cobra.Command{
	Use:   "testsuite",
	Short: "",
	RunE:  testSuiteRunE,
}

var (
	testSuiteName  *string
	testSuites     = map[string]func(cmd *cobra.Command, args []string) error{}
	errNoTestSuite = errors.New(`no test suite`)
)

func TestSuiteRegister(name string, fn func(cmd *cobra.Command, args []string) error) {
	testSuites[name] = fn
}

func testSuiteRunE(cmd *cobra.Command, args []string) error {
	err := config.ParseConfig()
	if err != nil {
		panic(err)
	}
	if fn, ok := testSuites[*testSuiteName]; ok {
		return fn(cmd, args)
	}
	return fmt.Errorf(`%w: %s`, errNoTestSuite, *testSuiteName)
}

func init() {
	rootCmd.AddCommand(testSuiteCmd)
	testSuiteName = testSuiteCmd.Flags().String("name", "", "name")
}
