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

func Match(value string, keys ...string) db.Compound {
	value = CleanFulltextOperator(value)
	return match(value, keys...)
}

func match(safelyMatchValue string, keys ...string) db.Compound {
	for idx, key := range keys {
		key = strings.ReplaceAll(key, "`", "``")
		keys[idx] = "`" + key + "`"
	}
	return db.Raw("MATCH(" + strings.Join(keys, ",") + ") AGAINST ('" + safelyMatchValue + "')")
}
