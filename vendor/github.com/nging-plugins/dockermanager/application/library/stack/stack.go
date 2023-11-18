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

package stack

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func New(username ...string) *Stack {
	var uname string
	if len(username) > 0 {
		uname = username[0]
	}
	return &Stack{Username: uname}
}

type Stack struct {
	ConfigFile    string
	ConfigContent string
	WorkDir       string
	RegistryAuth  string
	StackName     string
	Username      string
	noticer       notice.Noticer
}

func (c *Stack) SetName(stackName string) *Stack {
	c.StackName = stackName
	return c
}

func (c *Stack) Noticer(eCtx echo.Context) notice.Noticer {
	if c.noticer == nil {
		c.noticer = notice.New(eCtx, `dockerStack`, c.Username, context.Background())
	}
	return c.noticer
}

func (c *Stack) SetUsername(username string) *Stack {
	c.Username = username
	return c
}

func (c *Stack) SetRegistryAuth(registryAuth string) *Stack {
	c.RegistryAuth = registryAuth
	return c
}

func (c *Stack) SetWorkDir(dir string) *Stack {
	c.WorkDir = dir
	return c
}

func (c *Stack) SetConfigFile(file string) *Stack {
	c.ConfigFile = file
	return c
}

func (c *Stack) SetConfigContent(content string) *Stack {
	c.ConfigContent = content
	return c
}

func (c *Stack) commonArgs(op string) []string {
	if op == `deploy` && len(c.ConfigFile) == 0 {
		c.ConfigFile = `-`
	}
	args := []string{`stack`, op}
	if len(c.ConfigFile) > 0 {
		args = append(args, `--compose-file`, c.ConfigFile)
	}
	if len(c.RegistryAuth) > 0 {
		args = append(args, `--with-registry-auth`, c.RegistryAuth)
	}
	return args
}

func (c *Stack) exec(ctx echo.Context, args []string) (outStr string, errStr string, err error) {
	return utils.RunCommand(ctx, `docker`, args, c.Noticer(ctx), func(cmd *exec.Cmd) {
		cmd.Dir = c.WorkDir
		if len(args) > 1 && args[1] == `deploy` && (len(c.ConfigFile) == 0 || c.ConfigFile == `-`) {
			buf := bytes.NewBufferString(c.ConfigContent)
			cmd.Stdin = buf
		}
	})
}

func (c *Stack) Up(ctx echo.Context) error {
	args := c.commonArgs(`deploy`)
	args = append(args, c.StackName)
	outStr, _, err := c.exec(ctx, args)
	if err != nil {
		return err
	}
	log.Info(outStr)
	return err
}

func (c *Stack) Down(ctx echo.Context) error {
	args := c.commonArgs(`rm`)
	args = append(args, c.StackName)
	outStr, _, err := c.exec(ctx, args)
	if err != nil {
		return err
	}
	log.Info(outStr)
	return err
}

func (c *Stack) Reload(ctx echo.Context) error {
	c.Down(ctx)
	return c.Up(ctx)
}

func (c *Stack) ListServices(ctx echo.Context) ([]ServiceItem, error) {
	args := c.commonArgs(`services`)
	args = append(args, `--format`, `json`)
	args = append(args, c.StackName)
	outStr, _, err := c.exec(ctx, args)
	if err != nil {
		return nil, err
	}
	outStr = strings.TrimSpace(outStr)
	rows := strings.Split(outStr, com.StrLF)
	list := make([]ServiceItem, 0, len(rows))
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if len(row) == 0 {
			continue
		}
		if !strings.HasPrefix(row, `{`) {
			continue
		}
		item := ServiceItem{}
		err = json.Unmarshal(com.Str2bytes(row), &item)
		if err != nil {
			log.Error(err)
			continue
		}
		list = append(list, item)
	}
	return list, err
}

func (c *Stack) ListTasks(ctx echo.Context) ([]TaskItem, error) {
	args := c.commonArgs(`ps`)
	args = append(args, `--format`, `json`)
	//args = append(args, `--no-trunc`)
	args = append(args, c.StackName)
	outStr, _, err := c.exec(ctx, args)
	if err != nil {
		return nil, err
	}
	outStr = strings.TrimSpace(outStr)
	rows := strings.Split(outStr, com.StrLF)
	list := make([]TaskItem, 0, len(rows))
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if len(row) == 0 {
			continue
		}
		if !strings.HasPrefix(row, `{`) {
			continue
		}
		item := TaskItem{}
		err = json.Unmarshal(com.Str2bytes(row), &item)
		if err != nil {
			log.Error(err)
			continue
		}
		list = append(list, item)
	}
	return list, err
}
