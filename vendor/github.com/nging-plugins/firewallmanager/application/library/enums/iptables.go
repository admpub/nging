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

// 本文件常量来自 iptables

const (
	// 传输协议
	ProtocolTCP  = `tcp`
	ProtocolUDP  = `udp`
	ProtocolICMP = `icmp`
	ProtocolAll  = `all`
)

var ProtocolList = []string{ProtocolAll, ProtocolTCP, ProtocolUDP, ProtocolICMP}

const (
	// 规则表之间的顺序
	// raw → mangle → nat → filter
	// 规则表表
	TableFilter = `filter` // 过滤数据包。三个链：INPUT、FORWARD、OUTPUT
	TableNAT    = `nat`    // 用于网络地址转换（IP、端口）。 三个链：PREROUTING、POSTROUTING、OUTPUT
	TableMangle = `mangle` // 修改数据包的服务类型、TTL、并且可以配置路由实现QOS。五个链：PREROUTING、POSTROUTING、INPUT、OUTPUT、FORWARD
	TableRaw    = `raw`    // 决定数据包是否被状态跟踪机制处理。两个链：OUTPUT、PREROUTING
)

var TableList = []string{TableFilter, TableNAT, TableMangle, TableRaw}

var TablesChains = map[string][]string{
	TableRaw:    {ChainOutput, ChainPreRouting},
	TableMangle: {ChainPreRouting, ChainInput, ChainOutput, ChainForward, ChainPostRouting},
	TableNAT:    {ChainPreRouting /*ChainOutput,*/, ChainPostRouting},
	TableFilter: {ChainInput, ChainOutput, ChainForward},
}

const (
	// 规则链之间的顺序
	// ● 入站: PREROUTING → INPUT
	// ● 出站: OUTPUT → POSTROUTING
	// ● 转发: PREROUTING → FORWARD → POSTROUTIN
	// 规则链
	ChainInput       = `INPUT`       // 进来的数据包应用此规则链中的策略
	ChainOutput      = `OUTPUT`      // 外出的数据包应用此规则链中的策略
	ChainForward     = `FORWARD`     // 转发数据包时应用此规则链中的策略
	ChainPreRouting  = `PREROUTING`  // 对数据包作路由选择前应用此链中的规则（所有的数据包进来的时侯都先由这个链处理）
	ChainPostRouting = `POSTROUTING` // 对数据包作路由选择后应用此链中的规则（所有的数据包出来的时侯都先由这个链处理）
)

var ChainList = []string{ChainPreRouting, ChainInput, ChainOutput, ChainForward, ChainPostRouting}

var InputIfaceChainList = []string{ChainPreRouting, ChainInput, ChainForward}    // PREROUTING、INPUT、FORWARD
var OutputIfaceChainList = []string{ChainOutput, ChainForward, ChainPostRouting} // FORWARD、OUTPUT、POSTROUTING

const (
	StateNew         = `NEW`         // 新连接
	StateEstablished = `ESTABLISHED` // 后续对话连接
	StateRelated     = `RELATED`     // 关联到其他连接的连接
	StateInvalid     = `INVALID`     // 无效连接(没有任何状态)
	StateUntracked   = `UNTRACKED`   // 无法找到相关的连接
)

var StateList = []string{StateNew, StateEstablished, StateRelated, StateInvalid, StateUntracked}

const (
	// 防火墙处理数据包的四种方式
	TargetAccept = `ACCEPT` // 允许数据包通过
	TargetDrop   = `DROP`   // 直接丢弃数据包，不给任何回应信息
	TargetReject = `REJECT` // 拒绝数据包通过，必要时会给数据发送端一个响应的信息
	TargetLog    = `LOG`    // 在 /var/log/messages 文件中记录日志信息，然后将数据包传递给下一条规则
)

var TargetList = []string{TargetAccept, TargetDrop, TargetReject, TargetLog}

const (
	RejectWithICMPPortUnreachable  = `icmp-port-unreachable` // default
	RejectWithICMPNetUnreachable   = `icmp-net-unreachable`
	RejectWithICMPHostUnreachable  = `icmp-host-unreachable`
	RejectWithICMPProtoUnreachable = `icmp-proto-unreachable`
	RejectWithICMPNetProhibited    = `icmp-net-prohibited`
	RejectWithICMPHostProhibited   = `icmp-host-prohibited`
	RejectWithICMPAdminProhibited  = `icmp-admin-prohibited`
)

var RejectWithList = []string{
	RejectWithICMPPortUnreachable, RejectWithICMPNetUnreachable,
	RejectWithICMPHostUnreachable, RejectWithICMPProtoUnreachable,
	RejectWithICMPNetProhibited, RejectWithICMPHostProhibited,
	RejectWithICMPAdminProhibited,
}

const (
	TCPFlagALL = `ALL` // = SYN,ACK,FIN,RST,URG,PSH
	TCPFlagSYN = `SYN`
	TCPFlagACK = `ACK`
	TCPFlagFIN = `FIN`
	TCPFlagRST = `RST`
	TCPFlagURG = `URG`
	TCPFlagPSH = `PSH`
)

var (
	DefaultTCPFlagsWithACK = []string{`ALL`, TCPFlagSYN + `,` + TCPFlagACK}
	DefaultTCPFlags        = []string{`ALL`, TCPFlagSYN}
	DefaultTCPFlagsSimple  = []string{`--syn`} // = DefaultTCPFlags
)
