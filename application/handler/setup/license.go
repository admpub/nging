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

package setup

import (
	"errors"

	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/library/config"
	"github.com/admpub/nging/v3/application/library/license"
	"github.com/webx-top/echo"
)

// License 获取商业授权
func License(c echo.Context) error {
	err := license.Check(c)
	/*
		if err != nil {
			err = license.Generate(nil)
			if err == nil {
				err = license.Check(c)
			}
		}
	//*/
	if err == nil {
		nextURL := c.Query(`next`)
		if len(nextURL) == 0 {
			nextURL = handler.URLFor(`/`)
		}
		config.Version.Licensed = true
		return c.Redirect(nextURL)
	}
	errMap := map[string]string{
		`Could not read license`:        c.T(`读取授权文件失败`),
		`Could not read private key`:    c.T(`读取私钥失败`),
		`Could not read public key`:     c.T(`读取公钥失败`),
		`Could not read machine number`: c.T(`获取机器码失败`),
		`Invalid private key`:           c.T(`私钥无效`),
		`Invalid public key`:            c.T(`公钥无效`),
		`Invalid License file`:          c.T(`授权文件无效`),
		`Unlicensed Version`:            c.T(`未授权版本`),
		`Invalid MachineID`:             c.T(`机器码无效`),
		`Invalid LicenseID`:             c.T(`授权ID无效`),
		`License expired`:               c.T(`授权已过期`),
		`License does not exist`:        c.T(`授权文件不存在`),
	}
	if errStr, ok := errMap[err.Error()]; ok {
		err = errors.New(errStr)
	}
	//需要重新获取授权文件
	if err == license.ErrLicenseNotFound {
		err = license.DownloadOnce(c)
		c.Set(`downloaded`, err == nil)
	} else {
		c.Set(`downloaded`, false)
	}

	c.Set(`licenseFile`, license.FilePath())
	c.Set(`productURL`, license.ProductDetailURL())
	c.Set(`fileName`, license.FileName())
	return c.Render(`license`, err)
}
