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

package compose

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func New(username ...string) *Compose {
	var uname string
	if len(username) > 0 {
		uname = username[0]
	}
	return &Compose{Username: uname}
}

type Compose struct {
	ConfigFile    string
	ConfigContent string
	WorkDir       string
	ProjectName   string
	Username      string
	noticer       notice.Noticer
}

func (c *Compose) SetUsername(username string) *Compose {
	c.Username = username
	return c
}

func (c *Compose) SetName(projectName string) *Compose {
	c.ProjectName = projectName
	return c
}

func (c *Compose) Noticer(eCtx echo.Context) notice.Noticer {
	if c.noticer == nil {
		c.noticer = notice.New(eCtx, `dockerStack`, c.Username, context.Background())
	}
	return c.noticer
}

func (c *Compose) SetWorkDir(dir string) *Compose {
	c.WorkDir = dir
	return c
}

func (c *Compose) SetConfigFile(file string) *Compose {
	c.ConfigFile = file
	return c
}

func (c *Compose) SetConfigContent(content string) *Compose {
	c.ConfigContent = content
	return c
}

func (c *Compose) commonArgs() []string {
	if len(c.ConfigFile) == 0 {
		c.ConfigFile = ConfigPath(c.ProjectName)
		if len(c.ConfigContent) > 0 {
			os.WriteFile(c.ConfigFile, com.Str2bytes(c.ConfigContent), os.ModePerm)
		}
	}
	args := []string{`-f`, c.ConfigFile, `--compatibility`}
	return args
}

func (c *Compose) exec(ctx echo.Context, args []string) error {
	command, args := utils.DockerCompose(args)
	outStr, errStr, err := utils.RunCommand(ctx, command, args, c.Noticer(ctx), func(cmd *exec.Cmd) {
		cmd.Dir = c.WorkDir
	})
	if err != nil {
		err = fmt.Errorf(`%w: %s`, err, errStr)
	} else {
		log.Info(outStr)
	}
	return err
}

func (c *Compose) Up(ctx echo.Context, daemon bool) error {
	args := c.commonArgs()
	//args = append(args, `--env-file`,`.env`)
	if len(c.WorkDir) > 0 {
		args = append(args, `--project-directory`, c.WorkDir)
	}
	if len(c.ProjectName) > 0 {
		args = append(args, `--project-name`, c.ProjectName)
	}
	args = append(args, `up`, `--build`)
	if daemon {
		args = append(args, `-d`)
	}
	return c.exec(ctx, args)
}

func (c *Compose) Reload(ctx echo.Context) error {
	return c.Up(ctx, true)
}

// List containers
func (c *Compose) ListContainers(ctx echo.Context, opts ...echo.H) ([]ContainerItem, error) {
	args := c.commonArgs()
	args = append(args, `ps`, `--format`, `json`, `--all`)
	var opt echo.H
	if len(opts) > 0 {
		opt = opts[0]
	}
	services := opt.String(`services`)
	if len(services) > 0 {
		args = append(args, `--services`, services)
	}
	status := opt.String(`status`)
	if len(status) > 0 {
		args = append(args, `--status`, status) // Filter services by status. Values: [paused | restarting | removing | running | dead | created | exited]
	}
	command, args := utils.DockerCompose(args)
	outStr, errStr, err := com.ExecCmdWithContext(ctx, command, args...)
	if err != nil {
		return nil, fmt.Errorf(`%w: %s`, err, errStr)
	}
	outStr = strings.TrimSpace(outStr)
	if len(outStr) == 0 {
		return nil, err
	}
	rows := strings.Split(outStr, com.StrLF)
	list := make([]ContainerItem, 0, len(rows))
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if len(row) == 0 {
			continue
		}
		if !strings.HasPrefix(row, `{`) {
			continue
		}
		item := ContainerItem{}
		err = json.Unmarshal(com.Str2bytes(row), &item)
		if err != nil {
			log.Error(err)
			continue
		}
		list = append(list, item)
	}
	return list, err
}

func (c *Compose) Down(ctx echo.Context) error {
	args := c.commonArgs()
	args = append(args, `down`)
	return c.exec(ctx, args)
}

func (c *Compose) StartService(ctx echo.Context, service string) error {
	args := c.commonArgs()
	args = append(args, `start`, service)
	return c.exec(ctx, args)
}

func (c *Compose) StopService(ctx echo.Context, service string) error {
	args := c.commonArgs()
	args = append(args, `stop`, service)
	return c.exec(ctx, args)
}

func (c *Compose) RestartService(ctx echo.Context, service string) error {
	args := c.commonArgs()
	args = append(args, `restart`, service)
	return c.exec(ctx, args)
}

func (c *Compose) ScaleService(ctx echo.Context, service string, replicas uint) error {
	args := c.commonArgs()
	args = append(args, `scale`, service+`=`+param.AsString(replicas))
	return c.exec(ctx, args)
}

func (c *Compose) RunCommand(ctx echo.Context, service string, command string, runArgs []string) error {
	args := c.commonArgs()
	args = append(args, `run`, service, command)
	args = append(args, runArgs...)
	return c.exec(ctx, args)
}
