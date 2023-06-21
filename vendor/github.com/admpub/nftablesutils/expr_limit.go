package nftablesutils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

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

// ParseLimits parse expr.Limit
// rateStr := `1+/p/s`
// rateStr := `1+/bytes/second`
func ParseLimits(rateStr string, burst uint32) (*expr.Limit, error) {
	e := &expr.Limit{
		Type:  expr.LimitTypePktBytes,
		Rate:  0,
		Over:  false,
		Unit:  expr.LimitTimeSecond,
		Burst: burst,
	}
	var err error
	parts := strings.SplitN(rateStr, `/`, 3)
	switch len(parts) {
	case 3:
		parts[2] = strings.TrimSpace(parts[2])
		if len(parts[2]) > 0 {
			switch parts[2][0] {
			case 's': // second
				e.Unit = expr.LimitTimeSecond
			case 'm': // minute
				e.Unit = expr.LimitTimeMinute
			case 'h': // hour
				e.Unit = expr.LimitTimeHour
			case 'd': // day
				e.Unit = expr.LimitTimeDay
			case 'w': // week
				e.Unit = expr.LimitTimeWeek
			}
		}
		fallthrough
	case 2:
		parts[1] = strings.TrimSpace(parts[1])
		if len(parts[1]) > 0 {
			switch parts[1][0] {
			case 'p': // pkts
				e.Type = expr.LimitTypePkts
			case 'b': // bytes
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
		e.Rate, err = strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			err = fmt.Errorf(`failed to ParseUint(%q) from %q: %w`, parts[0], rateStr, err)
		}
	}
	return e, err
}

func SetConnLimits(connLimit uint32, rateStr string, burst uint32) (
	[]expr.Any, error) {
	exprLimit, err := ParseLimits(rateStr, burst)
	if err != nil {
		return nil, err
	}
	exprs := make([]expr.Any, 0, 2)
	if connLimit > 0 {
		exprs = append(exprs, &expr.Connlimit{
			Count: connLimit,
			Flags: 1,
		})
	}
	exprs = append(exprs, exprLimit)
	return exprs, err
}

func SetDynamicLimitDropSet(set *nftables.Set, connLimit uint32, rateStr string, burst uint32) (
	[]expr.Any, error) {
	if !set.Dynamic {
		return nil, errors.New(`must set *nftables.Set.Dynamic=true`)
	}
	if !set.HasTimeout {
		return nil, errors.New(`must set *nftables.Set.HasTimeout=true`)
	}
	if set.Timeout == 0 {
		return nil, errors.New(`*nftables.Set.Timeout must be set to greater than 0`)
	}
	exprs, err := SetConnLimits(connLimit, rateStr, burst)
	if err != nil {
		return nil, err
	}
	return []expr.Any{
		&expr.Dynset{
			SrcRegKey: defaultRegister,
			SetName:   set.Name,
			Operation: uint32(unix.NFT_DYNSET_OP_ADD),
			Exprs:     exprs,
		},
		&expr.Verdict{
			Kind: expr.VerdictDrop,
		},
	}, err
}
