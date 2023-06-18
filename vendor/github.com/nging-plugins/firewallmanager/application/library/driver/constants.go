package driver

// Protocol to differentiate between IPv4 and IPv6
type Protocol byte

const (
	ProtocolIPv4 Protocol = iota
	ProtocolIPv6
)
