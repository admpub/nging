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

package config

import (
	"github.com/webx-top/echo/param"
)

type Config struct {
	Verbose      bool
	Backend      string
	SaveFilePath string
	NgingRule    *NgingRule
}

// NgingRule Nging 自身的防火墙规则
type NgingRule struct {
	IPWhitelist string   // IP白名单（如果不设置表示不限制）
	RpsLimit    uint     // 频率限制规则（[NgingRpsLimit]+/p/s）
	RateBurst   uint     // 频率最大峰值
	RateExpires uint     // 限制时间（秒）
	OtherPort   []uint16 // 一般不需要设置。如果 Nging 还使用了其它端口则在此设置
}

func (a *NgingRule) OtherPortStrs(seperator ...string) []string {
	r := make([]string, len(a.OtherPort))
	for i, v := range a.OtherPort {
		r[i] = param.AsString(v)
	}
	return r
}
