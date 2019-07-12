package mysql

import (
	"github.com/webx-top/db/lib/sqlbuilder"
)

func NewBuilder(sess interface{}, prefixx ...string) sqlbuilder.SQLBuilder {
	if sess != nil {
		return sqlbuilder.WithSession(sess, template, prefixx...)
	}
	return sqlbuilder.WithTemplate(template, prefixx...)
}
