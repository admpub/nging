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
	"syscall"

	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

var forkSelfCmd = &cobra.Command{
	Use:     "forkself",
	Short:   "Fork self",
	Example: "forkself",
	RunE:    forkSelfRunE,
}

func forkSelfRunE(cmd *cobra.Command, args []string) error {
	executable, err := osext.Executable()
	if err != nil {
		return err
	}
	_, err = os.StartProcess(executable, append([]string{executable}, args...), &os.ProcAttr{
		Dir:   filepath.Dir(executable),
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Sys:   &syscall.SysProcAttr{},
	})
	return err
}

func init() {
	rootCmd.AddCommand(forkSelfCmd)
}
