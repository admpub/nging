package nftablesutils

import (
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

// SetConntrackStateSet helper.
func SetConntrackStateSet(s *nftables.Set) Exprs {
	exprs := []expr.Any{
		ExprCtState(defaultRegister),
		ExprLookupSet(defaultRegister, s.Name, s.ID),
	}
	return exprs
}

// SetConntrackStateNew helper.
func SetConntrackStateNew() Exprs {
	exprs := []expr.Any{
		ExprCtState(defaultRegister),
		ExprBitwise(defaultRegister, defaultRegister, ConnTrackStateLen,
			TypeConntrackStateNew(),
			[]byte{0x00, 0x00, 0x00, 0x00},
		),
		ExprCmpNeq(defaultRegister, []byte{0x00, 0x00, 0x00, 0x00}),
	}
	return exprs
}

// SetConntrackStateEstablished helper.
func SetConntrackStateEstablished() Exprs {
	exprs := []expr.Any{
		ExprCtState(defaultRegister),
		ExprBitwise(defaultRegister, defaultRegister, ConnTrackStateLen,
			TypeConntrackStateEstablished(),
			[]byte{0x00, 0x00, 0x00, 0x00},
		),
		ExprCmpNeq(defaultRegister, []byte{0x00, 0x00, 0x00, 0x00}),
	}
	return exprs
}

// SetConntrackStateRelated helper.
func SetConntrackStateRelated() Exprs {
	exprs := []expr.Any{
		ExprCtState(defaultRegister),
		ExprBitwise(defaultRegister, defaultRegister, ConnTrackStateLen,
			TypeConntrackStateRelated(),
			[]byte{0x00, 0x00, 0x00, 0x00},
		),
		ExprCmpNeq(defaultRegister, []byte{0x00, 0x00, 0x00, 0x00}),
	}
	return exprs
}

// GetConntrackStateSet helper.
func GetConntrackStateSet(t *nftables.Table) *nftables.Set {
	s := &nftables.Set{
		Anonymous: true,
		Constant:  true,
		Table:     t,
		KeyType:   TypeConntrackStateDatatype(),
	}
	return s
}

const (
	StateNew         = `new`
	StateEstablished = `established`
	StateRelated     = `related`
)

// GetConntrackStateSetElems helper.
func GetConntrackStateSetElems(states []string) []nftables.SetElement {
	elems := make([]nftables.SetElement, 0, len(states))
	for _, s := range states {
		switch s {
		case StateNew:
			elems = append(elems,
				nftables.SetElement{Key: TypeConntrackStateNew()})
		case StateEstablished:
			elems = append(elems,
				nftables.SetElement{Key: TypeConntrackStateEstablished()})
		case StateRelated:
			elems = append(elems,
				nftables.SetElement{Key: TypeConntrackStateRelated()})
		}
	}

	return elems
}
