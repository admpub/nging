package nftablesutils

import (
	"github.com/google/nftables"
	"golang.org/x/sys/unix"
)

var (
	typeProtoICMP                 = []byte{unix.IPPROTO_ICMP}
	typeProtoICMPV6               = []byte{unix.IPPROTO_ICMPV6}
	typeICMPTypeEchoRequest       = []byte{unix.ECHO}
	typeICMPv6TypeEchoRequest     = []byte{128}
	typeProtoUDP                  = []byte{unix.IPPROTO_UDP}
	typeProtoTCP                  = []byte{unix.IPPROTO_TCP}
	typeConntrackStateNew         = []byte{0x08, 0x00, 0x00, 0x00}
	typeConntrackStateEstablished = []byte{0x02, 0x00, 0x00, 0x00}
	typeConntrackStateRelated     = []byte{0x04, 0x00, 0x00, 0x00}
)

// TypeProtoICMP bytes.
func TypeProtoICMP() []byte {
	return typeProtoICMP
}

// TypeProtoICMPV6 bytes.
func TypeProtoICMPV6() []byte {
	return typeProtoICMPV6
}

// TypeICMPTypeEchoRequest bytes.
func TypeICMPTypeEchoRequest() []byte {
	return typeICMPTypeEchoRequest
}

// TypeICMPv6TypeEchoRequest bytes.
func TypeICMPv6TypeEchoRequest() []byte {
	return typeICMPv6TypeEchoRequest
}

// TypeProtoUDP bytes.
func TypeProtoUDP() []byte {
	return typeProtoUDP
}

// TypeProtoTCP bytes.
func TypeProtoTCP() []byte {
	return typeProtoTCP
}

// TypeConntrackStateNew bytes.
func TypeConntrackStateNew() []byte {
	return typeConntrackStateNew
}

// TypeConntrackStateEstablished bytes.
func TypeConntrackStateEstablished() []byte {
	return typeConntrackStateEstablished
}

// TypeConntrackStateRelated bytes.
func TypeConntrackStateRelated() []byte {
	return typeConntrackStateRelated
}

// ConntrackStateDatatype object.
func TypeConntrackStateDatatype() nftables.SetDatatype {
	return nftables.TypeCTState
}
