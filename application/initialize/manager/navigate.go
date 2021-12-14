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

package manager

import (
	"github.com/admpub/nging/v4/application/cmd/event"
	//"github.com/admpub/nging/v4/application/handler/caddy"
	"github.com/admpub/nging/v4/application/handler/cloud"
	//"github.com/admpub/nging/v4/application/handler/collector"
	//"github.com/admpub/nging/v4/application/handler/database"
	//"github.com/admpub/nging/v4/application/handler/download"
	//"github.com/admpub/nging/v4/application/handler/frp"
	//"github.com/admpub/nging/v4/application/handler/ftp"
	//"github.com/admpub/nging/v4/application/handler/server"
	"github.com/admpub/nging/v4/application/handler/task"
	//"github.com/admpub/nging/v4/application/handler/term"
	"github.com/admpub/nging/v4/application/registry/navigate"
)

var LeftNavigate = &navigate.List{
	//caddy.LeftNavigate,
	//server.LeftNavigate,
	//ftp.LeftNavigate,
	//collector.LeftNavigate,
	task.LeftNavigate,
	//download.LeftNavigate,
	cloud.LeftNavigate,
	//database.LeftNavigate,
	//frp.LeftNavigate,
	//term.LeftNavigate,
}

var Project = navigate.NewProject(`Nging`, `nging`, `/index`, navigate.LeftNavigate)

func init() {
	event.SupportManager = true
	navigate.LeftNavigate.Add(0, *LeftNavigate...)
	navigate.ProjectAdd(2, Project)
}
