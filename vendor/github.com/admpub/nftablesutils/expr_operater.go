package nftablesutils

import "github.com/google/nftables/expr"

type Operator string

func (o Operator) CmpOp() expr.CmpOp {
	switch o {
	case "!=":
		return expr.CmpOpNeq
	case ">":
		return expr.CmpOpGt
	case ">=":
		return expr.CmpOpGte
	case "<":
		return expr.CmpOpLt
	case "<=":
		return expr.CmpOpLte
	}

	return expr.CmpOpEq
}

func (o Operator) Expr() *expr.Cmp {
	return &expr.Cmp{
		Register: defaultRegister,
		Op:       o.CmpOp(),
	}
}
