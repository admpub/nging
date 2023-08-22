package mysql

import (
	"strings"

	"github.com/webx-top/db"
)

type Operator string

const (
	OperatorEQ           = `eq`
	OperatorMatch        = `match`
	OperatorSearchSuffix = `seachSuffix`
	OperatorSearchPrefix = `searchPrefix`
	OperatorSearchMiddle = `searchMiddle`
)

type fieldOp struct {
	field    string
	operator Operator
}

func (f fieldOp) isLikeQuery() bool {
	return f.operator == OperatorSearchMiddle || f.operator == OperatorSearchPrefix || f.operator == OperatorSearchSuffix
}

func (f fieldOp) isMatchQuery() bool {
	return f.operator == OperatorMatch
}

func (f fieldOp) isEqualQuery() bool {
	return f.operator == OperatorEQ
}

func (f fieldOp) buildCondMatch(values []string, matchValues *map[string][]string) bool {
	if f.operator == OperatorMatch {
		if _, ok := (*matchValues)[f.field]; !ok {
			(*matchValues)[f.field] = values
		} else {
			(*matchValues)[f.field] = append((*matchValues)[f.field], values...)
		}
		return true
	}
	return false
}

func buildSafelyMatchValues(v string, matchAll bool) []string {
	values := strings.Split(v, ` `)
	vals := make([]string, 0, len(values))
	for _, val := range values {
		val = strings.TrimSpace(val)
		if len(val) == 0 {
			continue
		}
		val = CleanFulltextOperator(val)
		if len(val) == 0 {
			continue
		}
		if matchAll {
			vals = append(vals, `+`+val)
		} else {
			vals = append(vals, val)
		}
	}
	return vals
}

func (f fieldOp) buildCondOther(likeValues []string, originalValues []string, cond ...*db.Compounds) *db.Compounds {
	var c *db.Compounds
	if len(cond) > 0 && cond[0] != nil {
		c = cond[0]
	} else {
		c = db.NewCompounds()
	}
	switch f.operator {
	case OperatorEQ:
		for _, val := range originalValues {
			c.AddKV(f.field, val)
		}
	case OperatorSearchPrefix:
		for _, val := range likeValues {
			c.AddKV(f.field, db.Like(val+`%`))
		}
	case OperatorSearchSuffix:
		for _, val := range likeValues {
			c.AddKV(f.field, db.Like(`%`+val))
		}
	default:
		for _, val := range likeValues {
			c.AddKV(f.field, db.Like(`%`+val+`%`))
		}
	}
	return c
}

func parseFieldOp(fields []string) []fieldOp {
	fieldConds := make([]fieldOp, len(fields))
	for i, f := range fields {
		if len(f) <= 1 {
			fieldConds[i] = fieldOp{field: f, operator: OperatorSearchMiddle}
			continue
		}
		switch f[0] {
		case '=':
			fieldConds[i] = fieldOp{field: f[1:], operator: OperatorEQ}
		case '~':
			fieldConds[i] = fieldOp{field: f[1:], operator: OperatorMatch}
		case '%':
			fieldConds[i] = fieldOp{field: f[1:], operator: OperatorSearchSuffix}
		default:
			if strings.HasSuffix(f, `%`) {
				f = f[0 : len(f)-1]
				fieldConds[i] = fieldOp{field: f, operator: OperatorSearchPrefix}
			} else {
				fieldConds[i] = fieldOp{field: f, operator: OperatorSearchMiddle}
			}
		}
	}
	return fieldConds
}
