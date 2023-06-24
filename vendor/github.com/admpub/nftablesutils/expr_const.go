package nftablesutils

import "net"

// Transport protocol lengths and offsets
const (
	SrcPortOffset = 0
	DstPortOffset = 2
	PortLen       = 2
)

// IPv4 lengths and offsets
const (
	IPv4SrcOffset = 12
	IPv4DstOffset = 16
	IPv4AddrLen   = net.IPv4len
)

// IPv6 lengths and offsets
const (
	IPv6SrcOffset = 8
	IPv6DstOffset = 24
	IPv6AddrLen   = net.IPv6len
)

const (
	ConnTrackStateLen = 4
)

const (
	ProtoTCPOffset = 9
	ProtoTCPLen    = 1
)

const (
	ProtoUDPOffset = 9
	ProtoUDPLen    = 1
)

const (
	ProtoICMPOffset = 9
	ProtoICMPLen    = 1
)

const (
	ProtoICMPv6Offset = 6
	ProtoICMPv6Len    = 1
)

// Default register and default xt_bpf version
const (
	defaultRegister = 1
	bpfRevision     = 1
)
