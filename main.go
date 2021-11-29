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

// 也可以以服务的方式启动nging
// 服务支持的操作有：
// nging service install  	-- 安装服务
// nging service uninstall  -- 卸载服务
// nging service start 		-- 启动服务
// nging service stop 		-- 停止服务
// nging service restart 	-- 重启服务
package main

//go:generate go install github.com/admpub/bindata/v3/go-bindata@latest
//go:generate go-bindata -fs -o bindata_assetfs.go -ignore "\\.(git|svn|DS_Store|less|scss)$" -minify "\\.(js|css)$" -tags bindata public/assets/... template/... config/i18n/...

import (
	_ "github.com/admpub/bindata/v3"
	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/cmd"
	_ "github.com/admpub/nging/v3/application/ico"

	//"github.com/admpub/nging/v3/application/library/loader"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	//register

	_ "github.com/admpub/nging/v3/application"
	_ "github.com/admpub/nging/v3/application/initialize/manager"
	_ "github.com/admpub/nging/v3/application/library/sqlite"

	"github.com/admpub/nging/v3/application/version"
)

var (
	BUILD_TIME string
	BUILD_OS   string
	BUILD_ARCH string
	CLOUD_GOX  string
	COMMIT     string
	LABEL      = `dev` //beta/alpha/stable
	VERSION    = `3.6.5`
	PACKAGE    = `free`

	schemaVer = version.DBSCHEMA //数据表结构版本
)

func main() {
	log.SetEmoji(com.IsMac)
	defer log.Close()
	// if err := loader.LoadPlugins(); err != nil {
	// 	panic(err)
	// }
	echo.Set(`BUILD_TIME`, BUILD_TIME)
	echo.Set(`BUILD_OS`, BUILD_OS)
	echo.Set(`BUILD_ARCH`, BUILD_ARCH)
	echo.Set(`COMMIT`, COMMIT)
	echo.Set(`LABEL`, LABEL)
	echo.Set(`VERSION`, VERSION)
	echo.Set(`PACKAGE`, PACKAGE)
	echo.Set(`SCHEMA_VER`, schemaVer)
	exec()
}

func exec() {
	cmd.Execute()
}
