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

// 也可以以服务的方式启动nging
// 服务支持的操作有：
// nging service install  	-- 安装服务
// nging service uninstall  -- 卸载服务
// nging service start 		-- 启动服务
// nging service stop 		-- 停止服务
// nging service restart 	-- 重启服务
package main

import (
	"time"

	_ "github.com/admpub/bindata/v3"
	"github.com/admpub/log"
	_ "github.com/admpub/nging/v5/application"
	_ "github.com/admpub/nging/v5/application/ico"

	"github.com/webx-top/com"

	//register

	"github.com/coscms/webcore"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/buildinfo"
	"github.com/coscms/webcore/library/module"
	"github.com/coscms/webcore/version"

	// module
	"github.com/admpub/nging/v5/application/handler/cloud"
	"github.com/admpub/nging/v5/application/handler/task"

	"github.com/nging-plugins/caddymanager"
	"github.com/nging-plugins/collector"
	"github.com/nging-plugins/dbmanager"
	"github.com/nging-plugins/ddnsmanager"
	"github.com/nging-plugins/dlmanager"
	"github.com/nging-plugins/frpmanager"
	"github.com/nging-plugins/ftpmanager"
	"github.com/nging-plugins/servermanager"
	"github.com/nging-plugins/sshmanager"
	"github.com/nging-plugins/webauthn"
)

var (
	BUILD_TIME string
	BUILD_OS   string
	BUILD_ARCH string
	CLOUD_GOX  string
	COMMIT     string
	LABEL      = `dev` //beta/alpha/stable
	VERSION    = `5.2.8`
	PACKAGE    = `free`

	schemaVer = version.DBSCHEMA //数据表结构版本
)

func main() {
	log.SetEmoji(com.IsMac)
	defer log.Close()
	bootconfig.ContainerImage = `admpub/nging`
	buildinfo.New(
		buildinfo.Time(BUILD_TIME),
		buildinfo.OS(BUILD_OS),
		buildinfo.Arch(BUILD_ARCH),
		buildinfo.CloudGox(CLOUD_GOX),
		buildinfo.Commit(COMMIT),
		buildinfo.Label(LABEL),
		buildinfo.Version(VERSION),
		buildinfo.Package(PACKAGE),
		buildinfo.SchemaVer(schemaVer),
	).Apply()
	if com.FileExists(`config/install.sql`) {
		com.Rename(`config/install.sql`, `config/install.sql.`+time.Now().Format(`20060102150405.000`))
	}
	webcore.Start(modules...)
}

var modules = []module.IModule{
	&caddymanager.Module,
	&servermanager.Module,
	&ftpmanager.Module,
	&collector.Module,
	&task.Module,
	&dlmanager.Module,
	&cloud.Module,
	&dbmanager.Module,
	&frpmanager.Module,
	&sshmanager.Module,
	&ddnsmanager.Module,
	&webauthn.Module,
}
