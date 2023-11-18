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
	"context"
	"encoding/json"
	"strings"

	"github.com/admpub/log"
	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/webx-top/com"
)

type Item struct {
	Name         string
	Namespace    string
	Orchestrator string
	Services     string
	Running      bool
}

func List(ctx context.Context, filters map[string]string) ([]Item, error) {
	args := []string{`stack`, `ls`, `--format`, `json`}
	for k, v := range filters {
		args = append(args, `--filter`, k+`=`+v)
	}
	outStr, errStr, err := com.ExecCmdWithContext(ctx, utils.DockerPath(), args...)
	if err != nil {
		return nil, err
	}
	_ = errStr
	// {"Name":"myapp","Namespace":"","Orchestrator":"Swarm","Services":"3"}
	outStr = strings.TrimSpace(outStr)
	rows := strings.Split(outStr, com.StrLF)
	list := make([]Item, 0, len(rows))
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if len(row) == 0 {
			continue
		}
		if !strings.HasPrefix(row, `{`) {
			continue
		}
		item := Item{}
		err = json.Unmarshal(com.Str2bytes(row), &item)
		if err != nil {
			log.Error(err)
			continue
		}
		item.Running = true
		list = append(list, item)
	}
	return list, err
}

// {"ID":"l00o8glnbr81","Image":"mysql:5.7","Mode":"replicated","Name":"ecc_mysql57","Ports":"*:33065-\u003e3306/tcp","Replicas":"1/1"}
type ServiceItem struct {
	ID       string
	Image    string
	Mode     string
	Name     string
	Ports    string
	Replicas string
}

// {"CurrentState":"Failed 4 minutes ago","DesiredState":"Shutdown","Error":"\"task: non-zero exit (1)\"","ID":"rjw0994fd78e","Image":"nginx:latest","Name":"ecc_nginx.1","Node":"orbstack","Ports":""}
type TaskItem struct {
	ID           string
	Image        string
	Name         string
	Node         string
	Ports        string
	CurrentState string
	DesiredState string
	Error        string
}
