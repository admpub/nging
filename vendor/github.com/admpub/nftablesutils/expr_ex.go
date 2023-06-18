package nftablesutils

import (
	"fmt"

	"github.com/google/nftables/binaryutil"
	"github.com/google/nftables/expr"
)

type ExprDirection string

const (
	ExprDirectionSource      ExprDirection = `source`
	ExprDirectionDestination ExprDirection = `destination`
)

var (
	zeroXor  = binaryutil.NativeEndian.PutUint32(0)
	zeroXor6 = append(binaryutil.NativeEndian.PutUint64(0), binaryutil.NativeEndian.PutUint64(0)...)
)

// GetPayloadDirectives get expression directives based on ip version and direction
func GetPayloadDirectives(direction ExprDirection, isIPv4 bool, isIPv6 bool) (uint32, uint32, []byte) {
	switch {
	case direction == ExprDirectionSource && isIPv4:
		return IPv4SrcOffset, IPv4AddrLen, zeroXor
	case direction == ExprDirectionDestination && isIPv4:
		return IPv4DstOffset, IPv4AddrLen, zeroXor
	case direction == ExprDirectionSource && isIPv6:
		return IPv6SrcOffset, IPv6AddrLen, zeroXor6
	case direction == ExprDirectionDestination && isIPv6:
		return IPv6DstOffset, IPv6AddrLen, zeroXor6
	default:
		panic("no matched payload directive")
	}
}

func LoadCtByKeyWithRegister(ctKey expr.CtKey, reg uint32) (*expr.Ct, error) {
	// Current upper and lower bound for valid CtKey values
	if ctKey < expr.CtKeySTATE || ctKey > expr.CtKeyEVENTMASK {
		return &expr.Ct{}, fmt.Errorf("invalid CtKey given")
	}

	return &expr.Ct{
		Register: reg,
		Key:      ctKey,
	}, nil
}

func LoadCtByKey(ctKey expr.CtKey) (*expr.Ct, error) {
	return LoadCtByKeyWithRegister(ctKey, defaultRegister)
}
