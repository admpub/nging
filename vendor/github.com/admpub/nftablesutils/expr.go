package nftablesutils

import (
	"net"
	"strconv"
	"strings"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

// Returns a meta expression
func ExprMeta(meta expr.MetaKey, reg uint32) *expr.Meta {
	return &expr.Meta{
		Key:      meta,
		Register: reg,
	}
}

// Returns a counter expression
func ExprCounter() *expr.Counter {
	return &expr.Counter{}
}

// ExprIIFName wrapper
func ExprIIFName() *expr.Meta {
	// [ meta load iifname => reg 1 ]
	return ExprMeta(expr.MetaKeyIIFNAME, defaultRegister)
}

// ExprOIFName wrapper
func ExprOIFName() *expr.Meta {
	// [ meta load oifname => reg 1 ]
	return ExprMeta(expr.MetaKeyOIFNAME, defaultRegister)
}

// ExprCmpEqIFName wrapper
func ExprCmpEqIFName(name string) *expr.Cmp {
	// [ cmp eq reg 1 0x00006f6c 0x00000000 0x00000000 0x00000000 ]
	return &expr.Cmp{
		Op:       expr.CmpOpEq,
		Register: defaultRegister,
		Data:     ifname(name),
	}
}

// ExprCmpNeqIFName wrapper
func ExprCmpNeqIFName(name string) *expr.Cmp {
	// [ cmp neq reg 1 0x00006f6c 0x00000000 0x00000000 0x00000000 ]
	return &expr.Cmp{
		Op:       expr.CmpOpNeq,
		Register: defaultRegister,
		Data:     ifname(name),
	}
}

func ifname(n string) []byte {
	b := make([]byte, 16)
	copy(b, []byte(n+"\x00"))
	return b
}

// ExprPayloadNetHeader wrapper
func ExprPayloadNetHeader(reg, offset, l uint32) *expr.Payload {
	// [ payload load 4b @ network header + 12 => reg 1 ]
	return &expr.Payload{
		DestRegister: reg,
		Base:         expr.PayloadBaseNetworkHeader,
		Offset:       offset,
		Len:          l,
	}
}

// ExprPayloadTransportHeader wrapper
func ExprPayloadTransportHeader(reg, offset, l uint32) *expr.Payload {
	// [ payload load 1b @ transport header + 0 => reg 1 ]
	return &expr.Payload{
		DestRegister: reg,
		Base:         expr.PayloadBaseTransportHeader,
		Offset:       offset,
		Len:          l,
	}
}

// ExprBitwise wrapper
func ExprBitwise(dReg, sReg, l uint32, mask, xor []byte) *expr.Bitwise {
	// [ bitwise reg 1 = (reg=1 & 0x000000ff ) ^ 0x00000000 ]
	return &expr.Bitwise{
		DestRegister:   dReg,
		SourceRegister: sReg,
		Len:            l,
		Mask:           mask,
		Xor:            xor,
	}
}

// ExprCmpEq wrapper
func ExprCmpEq(reg uint32, data []byte) *expr.Cmp {
	// [ cmp eq reg 1 0x0000007f ]
	return &expr.Cmp{
		Op:       expr.CmpOpEq,
		Register: reg,
		Data:     data,
	}
}

// ExprCmpNeq wrapper
func ExprCmpNeq(reg uint32, data []byte) *expr.Cmp {
	// [ cmp eq reg 1 0x0000007f ]
	return &expr.Cmp{
		Op:       expr.CmpOpNeq,
		Register: reg,
		Data:     data,
	}
}

// ExprLookupSetFromSet wrapper
func ExprLookupSetFromSet(set *nftables.Set, reg uint32) *expr.Lookup {
	return ExprLookupSet(reg, set.Name, set.ID)
}

// ExprLookupSet wrapper
func ExprLookupSet(reg uint32, name string, id uint32) *expr.Lookup {
	// [ lookup reg 1 set adminipset ]
	return &expr.Lookup{
		SourceRegister: defaultRegister,
		SetName:        name,
		SetID:          id,
	}
}

// ExprCtState wrapper
func ExprCtState(reg uint32) *expr.Ct {
	// [ ct load state => reg 1 ]
	return &expr.Ct{
		// Key:      unix.NFT_CT_STATE,
		Key:      expr.CtKeySTATE,
		Register: reg,
	}
}

// ExprImmediate wrapper
func ExprImmediate(ip net.IP) *expr.Immediate {
	// [ immediate reg 1 0x0158a8c0 ]
	return &expr.Immediate{
		Register: defaultRegister,
		Data:     ip,
	}
}

// ExprAccept wrapper
func ExprAccept() *expr.Verdict {
	// [ immediate reg 0 accept ]
	return &expr.Verdict{
		Kind: expr.VerdictAccept,
	}
}

// ExprDrop wrapper
func ExprDrop() *expr.Verdict {
	// [ immediate reg 0 accept ]
	return &expr.Verdict{
		Kind: expr.VerdictDrop,
	}
}

// ExprReject wrapper
func ExprReject(t uint32, c uint8) *expr.Reject {
	// [ reject type 0 code 3 ]
	return &expr.Reject{Type: t, Code: c}
}

// ExprConnLimit wrapper
func ExprConnLimit(count uint32, flags uint32) *expr.Connlimit {
	return &expr.Connlimit{
		Count: count,
		Flags: flags,
	}
}

// ExprLimit wrapper
func ExprLimit(t expr.LimitType, rate uint64, over bool, unit expr.LimitTime, burst uint32) *expr.Limit {
	return &expr.Limit{
		Type:  t,
		Rate:  rate,
		Over:  over,
		Unit:  unit,
		Burst: burst,
	}
}

// ExprLimits ExprLimit wrapper
// rateStr := `1+/p/s`
func ExprLimits(rateStr string, burst uint32) *expr.Limit {
	// 1+/p/s
	e := &expr.Limit{
		Type:  expr.LimitTypePktBytes,
		Rate:  0,
		Over:  false,
		Unit:  expr.LimitTimeSecond,
		Burst: burst,
	}

	parts := strings.SplitN(rateStr, `/`, 3)
	switch len(parts) {
	case 3:
		parts[2] = strings.TrimSpace(parts[2])
		if len(parts[2]) > 0 {
			switch parts[2][0] {
			case 's':
				e.Unit = expr.LimitTimeSecond
			case 'm':
				e.Unit = expr.LimitTimeMinute
			case 'h':
				e.Unit = expr.LimitTimeHour
			case 'd':
				e.Unit = expr.LimitTimeDay
			case 'w':
				e.Unit = expr.LimitTimeWeek
			}
		}
		fallthrough
	case 2:
		parts[1] = strings.TrimSpace(parts[1])
		if len(parts[1]) > 0 {
			switch parts[1][0] {
			case 'p':
				e.Type = expr.LimitTypePkts
			case 'b':
				e.Type = expr.LimitTypePktBytes
			}
		}
		fallthrough
	case 1:
		parts[0] = strings.TrimSpace(parts[0])
		e.Over = strings.HasSuffix(parts[0], `+`)
		if e.Over {
			parts[0] = strings.TrimSuffix(parts[0], `+`)
		}
		e.Rate, _ = strconv.ParseUint(parts[0], 10, 64)
	}
	return e
}
