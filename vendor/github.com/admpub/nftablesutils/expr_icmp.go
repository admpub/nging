package nftablesutils

import "github.com/google/nftables/expr"

// SetICMPTypeEchoRequest helper.
func SetICMPTypeEchoRequest() Exprs {
	exprs := []expr.Any{
		ExprPayloadTransportHeader(defaultRegister, 0, 1),
		ExprCmpEq(defaultRegister, TypeICMPTypeEchoRequest()),
	}

	return exprs
}

// SetICMPv6TypeEchoRequest helper.
func SetICMPv6TypeEchoRequest() Exprs {
	exprs := []expr.Any{
		ExprPayloadTransportHeader(defaultRegister, 0, 1),
		ExprCmpEq(defaultRegister, TypeICMPv6TypeEchoRequest()),
	}

	return exprs
}
