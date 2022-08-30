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

package frp

import (
	"github.com/webx-top/echo"

	"github.com/nging-plugins/frpmanager/application/dbschema"
)

func NewProxyConfig() *ProxyConfg {
	return &ProxyConfg{
		Proxy:   echo.H{},
		Visitor: echo.H{},
	}
}

type ProxyConfg struct {
	Proxy   echo.H
	Visitor echo.H
}

func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		Extra:          echo.H{},
		NgingFrpClient: dbschema.NewNgingFrpClient(nil),
	}
}

type ClientConfig struct {
	Extra echo.H
	*dbschema.NgingFrpClient
}
