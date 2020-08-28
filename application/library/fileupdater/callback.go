package fileupdater

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
)

func GenCallbackDefault(fieldName string, fieldValues ...FieldValue) CallbackFunc {
	return func(m factory.Model) (tableID string, content string, property *Property) {
		row := m.AsRow()
		tableID = row.String(`id`, `-1`)
		content = row.String(fieldName)
		property = NewPropertyWith(
			m,
			db.Cond{`id`: row.Get(`id`, `-1`)},
			fieldValues...,
		)
		return
	}
}

func GenCallbackWithCond(cond db.Compound, fieldValues ...FieldValue) CallbackFunc {
	return func(m factory.Model) (tableID string, content string, property *Property) {
		property = NewPropertyWith(
			m,
			cond,
			fieldValues...,
		)
		return
	}
}
