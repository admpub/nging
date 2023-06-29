package nftablesutils

import "github.com/google/nftables/expr"

func ProtoTCP(reg uint32) *expr.Payload {
	return ExprPayloadNetHeader(reg, ProtoTCPOffset, ProtoTCPLen)
}

func ProtoUDP(reg uint32) *expr.Payload {
	return ExprPayloadNetHeader(reg, ProtoUDPOffset, ProtoUDPLen)
}

// SetProtoICMP helper.
func SetProtoICMP(isEq ...bool) Exprs {
	exprs := []expr.Any{
		ExprPayloadNetHeader(defaultRegister, ProtoICMPOffset, ProtoICMPLen),
		ExprCmp(GetCmpOp(isEq...), TypeProtoICMP()),
	}
	return exprs
}

func SetProtoICMPv6(isEq ...bool) Exprs {
	return []expr.Any{
		// payload load 1b @ network header + 6 => reg 1
		ExprPayloadNetHeader(defaultRegister, ProtoICMPv6Offset, ProtoICMPv6Len),
		// cmp eq reg 1 0x0000003a
		ExprCmp(GetCmpOp(isEq...), TypeProtoICMPV6()),
	}
}

// SetINetProtoICMP helper.
func SetINetProtoICMP(isEq ...bool) Exprs {
	exprs := []expr.Any{
		ExprMeta(expr.MetaKeyL4PROTO, defaultRegister),
		ExprCmp(GetCmpOp(isEq...), TypeProtoICMP()),
	}
	return exprs
}

func SetINetProtoICMPv6(isEq ...bool) Exprs {
	return []expr.Any{
		ExprMeta(expr.MetaKeyL4PROTO, defaultRegister),
		ExprCmp(GetCmpOp(isEq...), TypeProtoICMPV6()),
	}
}

// SetProtoUDP helper.
func SetProtoUDP(isEq ...bool) Exprs {
	exprs := []expr.Any{
		ExprMeta(expr.MetaKeyL4PROTO, defaultRegister),
		ExprCmp(GetCmpOp(isEq...), TypeProtoUDP()),
	}
	return exprs
}

// SetProtoTCP helper.
func SetProtoTCP(isEq ...bool) Exprs {
	exprs := []expr.Any{
		ExprMeta(expr.MetaKeyL4PROTO, defaultRegister),
		ExprCmp(GetCmpOp(isEq...), TypeProtoTCP()),
	}
	return exprs
}
