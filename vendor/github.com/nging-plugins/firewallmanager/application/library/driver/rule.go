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

import "strconv"

type Rule struct {
	ID        uint   `json:"id,omitempty" xml:"id,omitempty"`
	Number    uint64 `json:"num,omitempty" xml:"num,omitempty"`
	Type      string `json:"type" xml:"type"` // filter / nat / etc.
	Name      string `json:"name" xml:"name"`
	Direction string `json:"direction" xml:"direction"` // INPUT / OUTPUT / etc.
	Action    string `json:"action" xml:"action"`       // ACCEPT / DROP / etc.
	Protocol  string `json:"protocol" xml:"protocol"`   // tcp / udp / etc.

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
	ConnLimit uint64 `json:"connLimit"  xml:"connLimit"` // 每个IP最大连接数
	RateLimit string `json:"rateLimit"  xml:"rateLimit"` // 频率限制规则（格式：200/pkt/second）
	RateBurst uint   `json:"rateBurst"  xml:"rateBurst"` // 频率最大峰值
}

func (r *Rule) IDBytes() []byte {
	if r.ID == 0 {
		return nil
	}
	s := strconv.FormatUint(uint64(r.ID), 10)
	return []byte(s)
}
