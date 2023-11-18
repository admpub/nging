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
	"path/filepath"
	"strings"

	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type ComposeItem struct {
	Name        string
	Status      string
	ConfigFiles string
	Running     bool
}

func List(ctx context.Context, filters map[string]string) ([]ComposeItem, error) {
	args := []string{`compose`, `ls`, `--format`, `json`, `--all`}
	for k, v := range filters {
		args = append(args, `--filter`, k+`=`+v)
	}
	outStr, errStr, err := com.ExecCmdWithContext(ctx, utils.DockerPath(), args...)
	if err != nil {
		return nil, fmt.Errorf(`%w: %s`, err, errStr)
	}
	outStr = strings.TrimSpace(outStr)
	if len(outStr) == 0 || !strings.HasPrefix(outStr, `[`) {
		return nil, err
	}
	var list []ComposeItem
	err = json.Unmarshal(com.Str2bytes(outStr), &list)
	for index, item := range list {
		item.Running = true
		list[index] = item
	}
	return list, err
}

type ContainerItem struct {
	Name       string
	Names      string
	Image      string
	Command    string
	Service    string
	CreatedAt  string
	Status     string
	Networks   string
	Ports      string
	Mounts     string
	State      string
	Size       string
	ExitCode   int
	ID         string
	Labels     string
	Publishers []ContainerPublisher
}

type ContainerPublisher struct {
	URL           string
	TargetPort    int
	PublishedPort int
	Protocol      string
}

/*
{
	"Command":"\"docker-entrypoint.sh --default-authentication-plugin=mysql_native_password\"",
	"CreatedAt":"2023-11-18 00:05:14 +0800 CST",
	"ExitCode":0,
	"Health":"",
	"ID":"419cedec8a5f218e18c641f4e1e78ea64da06e3522bcfbe01e3739fa67c1199a",
	"Image":"mysql:5.7",
	"Labels":"com.docker.compose.depends_on=,com.docker.compose.oneoff=False,com.docker.compose.project=docker-ecaio,com.docker.compose.project.working_dir=/Users/hank/work-doc/yicai/docker/docker-ecaio,com.docker.compose.config-hash=fdc2ace019c522cddedd7c2589289e4e837f2d0e2df37355c8aa6daec0c628a8,com.docker.compose.image=sha256:92034fe9a41f4344b97f3fc88a8796248e2cfa9b934be58379f3dbc150d07d9d,com.docker.compose.project.config_files=/Users/hank/work-doc/yicai/docker/docker-ecaio/docker-compose-local.yml,com.docker.compose.service=mysql57,com.docker.compose.version=2.22.0,com.docker.compose.container-number=1",
	"LocalVolumes":"0",
	"Mounts":"/Users/hank/work-doc/yicai/docker/docker-ecaio/mysql/data,/Users/hank/work-doc/yicai/docker/docker-ecaio/mysql/conf",
	"Name":"mysql57",
	"Names":"mysql57",
	"Networks":"docker-ecaio_default",
	"Ports":"33060/tcp, 0.0.0.0:33065-\u003e3306/tcp, :::33065-\u003e3306/tcp",
	"Publishers":[
		{"URL":"0.0.0.0","TargetPort":3306,"PublishedPort":33065,"Protocol":"tcp"},
		{"URL":"::","TargetPort":3306,"PublishedPort":33065,"Protocol":"tcp"},
		{"URL":"","TargetPort":33060,"PublishedPort":0,"Protocol":"tcp"}
	],
	"RunningFor":"12 hours ago",
	"Service":"mysql57",
	"Size":"0B",
	"State":"running",
	"Status":"Up 12 hours"
}
*/

func ConfigPath(name string) string {
	ppath := filepath.Join(echo.Wd(), `data`, `docker`, `compose`, name, `docker-compose.yml`)
	com.MkdirAll(filepath.Dir(ppath), os.ModePerm)
	return ppath
}
