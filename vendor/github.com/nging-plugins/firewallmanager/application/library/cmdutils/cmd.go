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

package cmdutils

import (
	"context"
	"errors"
	"io"
	"os/exec"

	"github.com/admpub/packer"
)

var ErrCmdForcedExit = errors.New(`cmd forced exit`)

func RunCmd(ctx context.Context, path string, args []string, stdout io.Writer, stdin ...io.Reader) error {
	if args == nil {
		args = []string{}
	}
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stdout = stdout
	cmd.Stderr = packer.Stderr
	if len(stdin) > 0 {
		cmd.Stdin = stdin[0]
	}

	if err := cmd.Run(); err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			return e
		default:
			return err
		}
	}
	return nil
}

func RunCmdWithCallback(ctx context.Context, path string, args []string, cb func(*exec.Cmd) error) error {
	if args == nil {
		args = []string{}
	}
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stderr = packer.Stderr
	if err := cb(cmd); err != nil {
		return err
	}
	//log.Debugf(`Firewall Command: %s`, cmd.String())
	if err := cmd.Run(); err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			return e
		default:
			if !errors.Is(err, ErrCmdForcedExit) {
				return err
			}
		}
	}
	return nil
}
