package sqlite

import (
	"github.com/webx-top/db/lib/sqlbuilder"
)

func NewBuilder(sess ...interface{}) sqlbuilder.SQLBuilder {
	if len(sess) > 0 && sess[0] != nil {
		return sqlbuilder.WithSession(sess[0], template)
	}
	return sqlbuilder.WithTemplate(template)
}
