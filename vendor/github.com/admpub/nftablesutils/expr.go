package nftablesutils

import (
	"net"

	"github.com/google/nftables"
	"github.com/google/nftables/binaryutil"
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

// ExprCmp wrapper
func ExprCmp(op expr.CmpOp, data []byte, reg ...uint32) *expr.Cmp {
	var register uint32
	if len(reg) > 0 {
		register = reg[0]
	}
	if register <= 0 {
		register = defaultRegister
	}
	return &expr.Cmp{
		Register: register,
		Op:       op,
		Data:     data,
	}
}

// ExprLookupSetFromSet wrapper
func ExprLookupSetFromSet(set *nftables.Set, reg uint32, isEq ...bool) *expr.Lookup {
	return ExprLookupSet(reg, set.Name, set.ID, isEq...)
}

// ExprLookupSet wrapper
func ExprLookupSet(reg uint32, name string, id uint32, isEq ...bool) *expr.Lookup {
	// [ lookup reg 1 set adminipset ]
	return &expr.Lookup{
		SourceRegister: defaultRegister,
		SetName:        name,
		SetID:          id,
		Invert:         IsInvert(isEq...),
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
func ExprImmediate(reg uint32, ip net.IP) *expr.Immediate {
	// [ immediate reg 1 0x0158a8c0 ]
	return &expr.Immediate{
		Register: reg,
		Data:     ip,
	}
}

// ExprImmediateWithPort wrapper
func ExprImmediateWithPort(reg uint32, port uint16) *expr.Immediate {
	// [ immediate reg 1 0x0158a8c0 ]
	return &expr.Immediate{
		Register: reg,
		Data:     binaryutil.BigEndian.PutUint16(port),
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
