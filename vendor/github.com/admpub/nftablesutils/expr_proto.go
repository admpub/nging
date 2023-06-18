package nftablesutils

import "github.com/google/nftables/expr"

func ProtoTCP(reg uint32) *expr.Payload {
	return ExprPayloadNetHeader(reg, ProtoTCPOffset, ProtoTCPLen)
}

func ProtoUDP(reg uint32) *expr.Payload {
	return ExprPayloadNetHeader(reg, ProtoUDPOffset, ProtoUDPLen)
}

// SetProtoICMP helper.
func SetProtoICMP() Exprs {
	exprs := []expr.Any{
		ExprPayloadNetHeader(defaultRegister, ProtoICMPOffset, ProtoICMPLen),
		ExprCmpEq(defaultRegister, TypeProtoICMP()),
	}

	return exprs
}

func SetProtoICMPv6() Exprs {
	return []expr.Any{
		// payload load 1b @ network header + 6 => reg 1
		ExprPayloadNetHeader(defaultRegister, ProtoICMPv6Offset, ProtoICMPv6Len),
		// cmp eq reg 1 0x0000003a
		ExprCmpEq(defaultRegister, TypeProtoICMPV6()),
	}
}

// SetINetProtoICMP helper.
func SetINetProtoICMP() Exprs {
	exprs := []expr.Any{
		ExprMeta(expr.MetaKeyL4PROTO, defaultRegister),
		ExprCmpEq(defaultRegister, TypeProtoICMP()),
	}

	return exprs
}

func SetINetProtoICMPv6() Exprs {
	return []expr.Any{
		ExprMeta(expr.MetaKeyL4PROTO, defaultRegister),
		ExprCmpEq(defaultRegister, TypeProtoICMPV6()),
	}
}

// SetProtoUDP helper.
func SetProtoUDP() Exprs {
	exprs := []expr.Any{
		ExprMeta(expr.MetaKeyL4PROTO, defaultRegister),
		ExprCmpEq(defaultRegister, TypeProtoUDP()),
	}

	return exprs
}

// SetProtoTCP helper.
func SetProtoTCP() Exprs {
	exprs := []expr.Any{
		ExprMeta(expr.MetaKeyL4PROTO, defaultRegister),
		ExprCmpEq(defaultRegister, TypeProtoTCP()),
	}

	return exprs
}
