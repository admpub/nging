package common

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	mysqlUtil "github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type SelectPageSortValues struct {
	PKName   string
	PKValues []string
}

func (s SelectPageSortValues) IsEmpty() bool {
	return len(s.PKValues) == 0
}

func (s SelectPageSortValues) IsMultiple() bool {
	return len(s.PKValues) > 1
}

func (s SelectPageSortValues) AddToSorts(sorts []interface{}) []interface{} {
	if s.IsMultiple() {
		sorts = append(sorts, db.Raw(s.OrderByString()))
	}
	return sorts
}

var SQLKeyReplacer = strings.NewReplacer("`", "``", ".", "`.`")

func (s SelectPageSortValues) OrderByString() string {
	values := make([]string, len(s.PKValues))
	for index, value := range s.PKValues {
		values[index] = com.AddSlashes(value)
	}
	return "FIELD(`" + SQLKeyReplacer.Replace(s.PKName) + "`,'" + strings.Join(values, `','`) + "')"
}

func SelectPageCond(ctx echo.Context, cond *db.Compounds, pkAndLabelFields ...string) (sv *SelectPageSortValues) {
	pk := `id`
	lb := `name`
	switch len(pkAndLabelFields) {
	case 2:
		if len(pkAndLabelFields[1]) > 0 {
			lb = pkAndLabelFields[1]
		}
		fallthrough
	case 1:
		if len(pkAndLabelFields[0]) > 0 {
			pk = pkAndLabelFields[0]
		}
	}
	searchValue := param.StringSlice(ctx.Formx(`searchValue`).Split(`,`)).Unique().Filter()
	if len(searchValue) > 0 {
		if len(searchValue) > 1 {
			cond.AddKV(pk, db.In(searchValue))
		} else {
			cond.AddKV(pk, searchValue[0])
		}
		sv = &SelectPageSortValues{
			PKName:   pk,
			PKValues: searchValue,
		}
	} else {
		keywords := ctx.FormValues(`q_word[]`)
		q := strings.Join(keywords, ` `)
		if len(q) == 0 {
			q = ctx.Formx(`q`).String()
		}
		if len(q) > 0 {
			cond.From(mysqlUtil.MatchAnyField(lb, q))
		}
	}
	return
}
