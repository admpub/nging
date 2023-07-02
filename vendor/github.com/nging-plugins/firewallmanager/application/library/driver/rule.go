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

package driver

import (
	"strconv"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type Rule struct {
	ID        uint   `json:"id,omitempty" xml:"id,omitempty"`             // 静态规则 ID
	CustomID  string `json:"customID,omitempty" xml:"customID,omitempty"` // 自定义 ID 字符串， ID 为 0 时有效
	Number    uint   `json:"num,omitempty" xml:"num,omitempty"`           // 防火墙的规则编号。iptables 为 position 值；nftables 为 handle 值
	Type      string `json:"type" xml:"type"`                             // 表 filter / nat / etc.
	Name      string `json:"name" xml:"name"`                             // 名称
	Direction string `json:"direction" xml:"direction"`                   // 链 INPUT / OUTPUT / etc.
	Action    string `json:"action" xml:"action"`                         // ACCEPT / DROP / etc.
	Protocol  string `json:"protocol" xml:"protocol"`                     // tcp / udp / etc.

	// interface 网口
	Interface string `json:"interface" xml:"interface"` // 网络入口网络接口
	Outerface string `json:"outerface" xml:"outerface"` // 网络出口网络接口

	// state
	State string `json:"state" xml:"state"`

	// IP or Port
	RemoteIP   string `json:"remoteIP" xml:"remoteIP"`
	LocalIP    string `json:"localIP" xml:"localIP"`
	NatIP      string `json:"natIP" xml:"natIP"`
	RemotePort string `json:"remotePort" xml:"remotePort"` // 支持指定范围
	LocalPort  string `json:"localPort" xml:"localPort"`   // 支持指定范围
	NatPort    string `json:"natPort" xml:"natPort"`       // 支持指定范围
	IPVersion  string `json:"ipVersion"  xml:"ipVersion"`  // 4 or 6

	// Limit
	ConnLimit   string `json:"connLimit"  xml:"connLimit"`     // 每个IP最大连接数
	RateLimit   string `json:"rateLimit"  xml:"rateLimit"`     // 频率限制规则（格式：200/p/s）
	RateBurst   uint   `json:"rateBurst"  xml:"rateBurst"`     // 频率最大峰值
	RateExpires uint   `json:"rateExpires"  xml:"rateExpires"` // 过期时间（秒）
	Extra       echo.H `json:"extra,omitempty"  xml:"extra,omitempty"`
}

func (r *Rule) IDBytes() []byte {
	if r.ID == 0 {
		return []byte(r.CustomID)
	}
	s := strconv.FormatUint(uint64(r.ID), 10)
	return []byte(s)
}

func (r *Rule) IDString() string {
	if r.ID == 0 {
		return r.CustomID
	}
	s := strconv.FormatUint(uint64(r.ID), 10)
	return s
}

func (r *Rule) GenLimitSetName() string {
	if r.ID == 0 {
		return com.SnakeCase(r.IDString())
	}
	return LimitSetNamePrefix + r.IDString()
}
