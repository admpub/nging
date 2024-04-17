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

package bootconfig

import (
	"net/http"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/tplfunc"
)

var (
	Bindata          bool
	StaticMW         interface{}
	HTTPCacheMaxAge  = com.GetenvDuration(`NGING_HTTPCACHE_MAXAGE`, time.Hour*24*7)
	BackendTmplMgr   driver.Manager
	FrontendTmplMgr  driver.Manager
	LangFSFunc       func(dir string) http.FileSystem
	Develop          bool
	SupportManager   bool
	MustLicensed     bool //是否必须被许可才能运行(如为true,则未许可的情况下会强制退出程序,否则不会退出程序) Must be licensed before starting
	FaviconHandler   func(echo.Context) error
	FaviconPath      = "public/assets/backend/images/favicon-xs.ico"
	SoftwareName     = `Nging` // 软件英文名(也作为软件二进制可执行文件的名称)
	SoftwareTitle    = `Nging` // 软件标题(显示在网页上的软件名称)
	OfficialHomepage = `https://github.com/admpub/nging`

	// Short 简述
	Short = `Nging is a web and network service management system`
	// Long 长述
	Long string
	// Welcome 欢迎语
	Welcome             = "Thank you for choosing nging %s, I hope you enjoy using it.\nToday is %s."
	AutoUpgradeDBStruct = true
)

func IsWeb() bool {
	return IsServerType(`web`)
}

func IsServerType(typ string) bool {
	return echo.String(`serverType`) == typ
}

func SetServerType(typ string) {
	echo.Set(`serverType`, typ)
}

func init() {
	tplfunc.TplFuncMap[`SoftwareName`] = func() string { return SoftwareName }
	tplfunc.TplFuncMap[`SoftwareTitle`] = func() string { return SoftwareTitle }
	tplfunc.TplFuncMap[`MustLicensed`] = func() bool { return MustLicensed }
	tplfunc.TplFuncMap[`Develop`] = func() bool { return Develop }
	tplfunc.TplFuncMap[`AutoUpgradeDBStruct`] = func() bool { return AutoUpgradeDBStruct }
	tplfunc.TplFuncMap[`FaviconPath`] = func() string { return FaviconPath }
	tplfunc.TplFuncMap[`ShortDescription`] = func() string { return Short }
	tplfunc.TplFuncMap[`LongDescription`] = func() string { return Long }
	tplfunc.TplFuncMap[`Welcome`] = func() string { return Welcome }
	tplfunc.TplFuncMap[`OfficialHomepage`] = func() string { return OfficialHomepage }
}
