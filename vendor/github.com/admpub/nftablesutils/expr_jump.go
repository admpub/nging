package nftablesutils

import "github.com/google/nftables/expr"

// Returns an accept verdict expression
func Accept() *expr.Verdict {
	return ExprAccept()
}

// Returns an drop verdict expression
func Drop() *expr.Verdict {
	return ExprDrop()
}

// Returns an reject expression
func Reject() *expr.Reject {
	return ExprReject(0, 0)
}
