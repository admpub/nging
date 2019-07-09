package mysql

import (
	"github.com/webx-top/db/lib/sqlbuilder"
)

func NewBuilder(sesses ...interface{}) sqlbuilder.SQLBuilder {
	if len(sesses) > 0 && sesses[0] != nil {
		return sqlbuilder.WithSession(sesses[0], template)
	}
	return sqlbuilder.WithTemplate(template)
}
