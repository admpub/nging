package nftablesutils

import (
	"github.com/google/nftables"
	"github.com/google/nftables/binaryutil"
	"github.com/google/nftables/expr"
)

// Returns a source port payload expression
func SourcePort(reg uint32) *expr.Payload {
	return ExprPayloadTransportHeader(reg, SrcPortOffset, PortLen)
}

// Returns a destination port payload expression
func DestinationPort(reg uint32) *expr.Payload {
	return ExprPayloadTransportHeader(reg, DstPortOffset, PortLen)
}

func GetCmpOp(isEq ...bool) expr.CmpOp {
	var cmpOp expr.CmpOp
	if len(isEq) > 0 && !isEq[0] {
		cmpOp = expr.CmpOpNeq
	} else {
		cmpOp = expr.CmpOpEq
	}
	return cmpOp
}

func IsInvert(isEq ...bool) bool {
	return len(isEq) > 0 && !isEq[0]
}

// SetSPort helper.
func SetSPort(p uint16, isEq ...bool) Exprs {
	exprs := []expr.Any{
		SourcePort(defaultRegister),
		ExprCmpPort(GetCmpOp(isEq...), p),
	}
	return exprs
}

// SetDPort helper.
func SetDPort(p uint16, isEq ...bool) Exprs {
	exprs := []expr.Any{
		DestinationPort(defaultRegister),
		ExprCmpPort(GetCmpOp(isEq...), p),
	}
	return exprs
}

// SetSPortSet helper.
func SetSPortSet(s *nftables.Set, isEq ...bool) Exprs {
	exprs := []expr.Any{
		SourcePort(defaultRegister),
		ExprLookupSet(defaultRegister, s.Name, s.ID, isEq...),
	}
	return exprs
}

// SetDPortSet helper.
func SetDPortSet(s *nftables.Set, isEq ...bool) Exprs {
	exprs := []expr.Any{
		DestinationPort(defaultRegister),
		ExprLookupSet(defaultRegister, s.Name, s.ID, isEq...),
	}
	return exprs
}

// ExprCmpPort returns a new port expression with the given matching operator.
func ExprCmpPort(op expr.CmpOp, port uint16, reg ...uint32) *expr.Cmp {
	return ExprCmp(op, binaryutil.BigEndian.PutUint16(port), reg...)
}

// SetSPortRange returns a new port range expression.
func SetSPortRange(min uint16, max uint16) []expr.Any {
	return []expr.Any{
		SourcePort(defaultRegister),
		ExprCmpPort(expr.CmpOpGte, min),
		ExprCmpPort(expr.CmpOpLte, max),
	}
}

// SetDPortRange returns a new port range expression.
func SetDPortRange(min uint16, max uint16) []expr.Any {
	return []expr.Any{
		DestinationPort(defaultRegister),
		ExprCmpPort(expr.CmpOpGte, min),
		ExprCmpPort(expr.CmpOpLte, max),
	}
}

// GetPortSet helper.
func GetPortSet(t *nftables.Table) *nftables.Set {
	s := &nftables.Set{
		Anonymous: true,
		Constant:  true,
		Table:     t,
		KeyType:   nftables.TypeInetService,
	}
	return s
}

// GetPortElems helper.
func GetPortElems(ports []uint16) []nftables.SetElement {
	elems := make([]nftables.SetElement, len(ports))
	for i, p := range ports {
		elems[i] = nftables.SetElement{Key: binaryutil.BigEndian.PutUint16(p)}
	}
	return elems
}
