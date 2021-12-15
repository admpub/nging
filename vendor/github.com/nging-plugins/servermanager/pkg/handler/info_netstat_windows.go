//go:build windows
// +build windows

// Package server Copyright 2016 Wenhui Shen <www.webx.top>
/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package handler

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/shirou/gopsutil/v3/net"
)

const (
	TCPTableBasicListener = iota
	TCPTableBasicConnections
	TCPTableBasicAll
	TCPTableOwnerPIDListener
	TCPTableOwnerPIDConnections
	TCPTableOwnerPIDAll
	TCPTableOwnerModuleListener
	TCPTableOwnerModuleConnections
	TCPTableOwnerModuleAll
)

const (
	UDPTableBasicListener = iota
	UDPTableOwnerPIDAll
	UDPTableOwnerModuleAll
)

type State int

const (
	MIB_TCP_STATE_CLOSED State = 1 + iota
	MIB_TCP_STATE_LISTEN
	MIB_TCP_STATE_SYN_SENT
	MIB_TCP_STATE_SYN_RCVD
	MIB_TCP_STATE_ESTAB
	MIB_TCP_STATE_FIN_WAIT1
	MIB_TCP_STATE_FIN_WAIT2
	MIB_TCP_STATE_CLOSE_WAIT
	MIB_TCP_STATE_CLOSING
	MIB_TCP_STATE_LAST_ACK
	MIB_TCP_STATE_TIME_WAIT
	MIB_TCP_STATE_DELETE_TCB
)

var (
	modIphlpapi             = syscall.NewLazyDLL("iphlpapi.dll")
	procGetExtendedTcpTable = modIphlpapi.NewProc("GetExtendedTcpTable")
	procGetExtendedUdpTable = modIphlpapi.NewProc("GetExtendedUdpTable")
	AF_INET                 = 2
	MaxCount                = 200
)

type MIB_TCPROW_OWNER_PID struct {
	dwState      uint32
	dwLocalAddr  uint32
	dwLocalPort  uint32
	dwRemoteAddr uint32
	dwRemotePort uint32
	dwOwningPid  uint32
}

type MIB_TCPTABLE_OWNER_PID struct {
	dwNumEntries uint32
	table        [200]MIB_TCPROW_OWNER_PID
}

type MIB_UDPROW_OWNER_PID struct {
	dwLocalAddr uint32
	dwLocalPort uint32
	dwOwningPid uint32
}

type MIB_UDPTABLE_OWNER_PID struct {
	dwNumEntries uint32
	table        [200]MIB_UDPROW_OWNER_PID
}

func NetStatTCP() (<-chan net.ConnectionStat, error) {
	b := make([]byte, MaxCount)
	size := uint32(len(b))
	ret, _, _ := procGetExtendedTcpTable.Call(
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		uintptr(AF_INET),
		TCPTableOwnerPIDAll,
		0)

	if ret == uintptr(syscall.ERROR_INSUFFICIENT_BUFFER) {
		b = make([]byte, size)
		ret, _, _ = procGetExtendedTcpTable.Call(
			uintptr(unsafe.Pointer(&b[0])),
			uintptr(unsafe.Pointer(&size)),
			0,
			uintptr(AF_INET),
			TCPTableOwnerPIDAll,
			0)
	}
	if ret != 0 {
		return nil, syscall.GetLastError()
	}
	ch := make(chan net.ConnectionStat)
	go func() {
		table := (*MIB_TCPTABLE_OWNER_PID)(unsafe.Pointer(&b[0]))
		for i := 0; i < int(table.dwNumEntries) && i < 200; i++ {
			row := net.ConnectionStat{}
			row.Status = getState(State(table.table[i].dwState))
			row.Laddr.IP = getIpAddress(table.table[i].dwLocalAddr)
			row.Laddr.Port = uint32(getPortNumber(table.table[i].dwLocalPort))
			row.Raddr.IP = getIpAddress(table.table[i].dwRemoteAddr)
			row.Raddr.Port = uint32(getPortNumber(table.table[i].dwRemotePort))
			row.Pid = int32(table.table[i].dwOwningPid)
			ch <- row
		}
		close(ch)
	}()
	return ch, nil
}

func NetStatUDP() (<-chan net.ConnectionStat, error) {
	b := make([]byte, MaxCount)
	size := uint32(len(b))
	ret, _, _ := procGetExtendedUdpTable.Call(
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		uintptr(AF_INET),
		UDPTableOwnerPIDAll,
		0)

	if ret == uintptr(syscall.ERROR_INSUFFICIENT_BUFFER) {
		b = make([]byte, size)
		ret, _, _ = procGetExtendedUdpTable.Call(
			uintptr(unsafe.Pointer(&b[0])),
			uintptr(unsafe.Pointer(&size)),
			0,
			uintptr(AF_INET),
			UDPTableOwnerPIDAll,
			0)
	}
	if ret != 0 {
		return nil, syscall.GetLastError()
	}
	ch := make(chan net.ConnectionStat)
	go func() {
		table := (*MIB_UDPTABLE_OWNER_PID)(unsafe.Pointer(&b[0]))
		for i := 0; i < int(table.dwNumEntries) && i < 200; i++ {
			row := net.ConnectionStat{}
			row.Laddr.IP = getIpAddress(table.table[i].dwLocalAddr)
			row.Laddr.Port = uint32(getPortNumber(table.table[i].dwLocalPort))
			row.Pid = int32(table.table[i].dwOwningPid)
			ch <- row
		}
		close(ch)
	}()
	return ch, nil
}

func getPortNumber(port uint32) int {
	return int(port)/256 + (int(port)%256)*256
}

func getIpAddress(ip uint32) string {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, ip)
	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

func getState(state State) string {
	m := map[State]string{
		MIB_TCP_STATE_CLOSED:     "CLOSED",
		MIB_TCP_STATE_LISTEN:     "LISTEN",
		MIB_TCP_STATE_SYN_SENT:   "SYN_SEND",
		MIB_TCP_STATE_SYN_RCVD:   "SYN_RECV",
		MIB_TCP_STATE_ESTAB:      "ESTABLISHED",
		MIB_TCP_STATE_FIN_WAIT1:  "FIN_WAIT_1",
		MIB_TCP_STATE_FIN_WAIT2:  "FIN_WAIT_2",
		MIB_TCP_STATE_CLOSE_WAIT: "CLOSE_WAIT",
		MIB_TCP_STATE_CLOSING:    "CLOSING",
		MIB_TCP_STATE_LAST_ACK:   "LAST_ACK",
		MIB_TCP_STATE_TIME_WAIT:  "TIME_WAIT",
		MIB_TCP_STATE_DELETE_TCB: "DELETE_TBC",
	}
	return m[state]
}
