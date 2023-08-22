package mysql

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
)

var fulltextOperatorReplacer = strings.NewReplacer(
	`'`, ``,
	`+`, ``,
	`-`, ``,
	`*`, ``,
	`"`, ``,
	`\`, ``,
)

func CleanFulltextOperator(v string) string {
	if com.StrIsAlphaNumeric(v) {
		return v
	}

	return fulltextOperatorReplacer.Replace(v)
}

func Match(value string, booleanMode bool, keys ...string) db.Compound {
	value = CleanFulltextOperator(value)
	return match(value, booleanMode, keys...)
}

func match(safelyMatchValue string, booleanMode bool, keys ...string) db.Compound {
	fields := make([]string, 0, len(keys))
	for _, key := range keys {
		if len(key) == 0 {
			continue
		}
		key = strings.ReplaceAll(key, "`", "``")
		fields = append(fields, "`"+key+"`")
	}
	var mode string
	if booleanMode {
		mode = ` IN BOOLEAN MODE`
	}
	return db.Raw("MATCH(" + strings.Join(fields, ",") + ") AGAINST ('" + safelyMatchValue + "'" + mode + ")")
}
