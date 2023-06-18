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

package enums

import (
	"github.com/webx-top/echo"
)

var Types = echo.NewKVData().
	Add(TableFilter, `过滤器`).
	Add(TableNAT, `网络地址转换器`)
	//Add(TableMangle, `Mangle`).
	//Add(TableRaw, `Raw`)

var Directions = echo.NewKVData().
	Add(ChainInput, `入站`).
	Add(ChainOutput, `出站`).
	Add(ChainForward, `转发`).
	Add(ChainPreRouting, `入站前`).
	Add(ChainPostRouting, `出站后`)

var IPProtocols = echo.NewKVData().
	Add(`4`, `IPv4`).
	Add(`6`, `IPv6`)

var NetProtocols = echo.NewKVData().
	Add(ProtocolTCP, `TCP`).
	Add(ProtocolUDP, `UDP`).
	Add(ProtocolICMP, `ICMP`).
	Add(ProtocolAll, `ALL`)

var Actions = echo.NewKVData().
	Add(TargetAccept, `✅ 接受`).
	Add(TargetDrop, `🚮 丢弃`).
	Add(TargetReject, `🚫 拒绝`).
	Add(TargetLog, `📝 记录日志`)
