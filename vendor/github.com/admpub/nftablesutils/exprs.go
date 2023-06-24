package nftablesutils

import (
	"github.com/google/nftables/expr"
)

type Exprs []expr.Any

func (e Exprs) Add(v ...expr.Any) Exprs {
	e = append(e, v...)
	return e
}

func JoinExprs(exprs ...[]expr.Any) Exprs {
	var sum int
	for _, vals := range exprs {
		sum += len(vals)
	}
	result := make([]expr.Any, 0, sum)
	for _, vals := range exprs {
		result = append(result, vals...)
	}
	return result
}
