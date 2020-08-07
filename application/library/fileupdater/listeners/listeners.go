/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package listeners

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/admpub/nging/application/library/fileupdater"
	"github.com/admpub/nging/application/library/fileupdater/listener"
)

func New() fileupdater.Listener {
	return &Listeners{}
}

type Listeners map[string]func(m factory.Model) (tableID string, content string, property *listener.Property)

func (a *Listeners) Listen(tableName string, embedded bool, seperatorAndSameFields ...string) {
	var sameFields []string
	var seperator string
	if len(seperatorAndSameFields) > 0 {
		seperator = seperatorAndSameFields[0]
	}
	if len(seperatorAndSameFields) > 1 {
		sameFields = seperatorAndSameFields[1:]
	}
	for _field, _listener := range *a {
		listener.New(_listener, embedded, seperator).SetTable(tableName, _field, sameFields...).ListenDefault()
		delete(*a, _field)
	}
}

func (a *Listeners) ListenByField(fieldNames string, tableName string, embedded bool, seperatorAndSameFields ...string) {
	if len(fieldNames) == 0 {
		return
	}
	var sameFields []string
	var seperator string
	if len(seperatorAndSameFields) > 0 {
		seperator = seperatorAndSameFields[0]
	}
	if len(seperatorAndSameFields) > 1 {
		sameFields = seperatorAndSameFields[1:]
	}
	for _, fieldName := range strings.Split(fieldNames, `,`) {
		if _listener, ok := (*a)[fieldName]; ok {
			listener.New(_listener, embedded, seperator).SetTable(tableName, fieldName, sameFields...).ListenDefault()
			delete(*a, fieldName)
		}
	}
}

func GenDefaultCallback(fieldName string) func(m factory.Model) (tableID string, content string, property *listener.Property) {
	return func(m factory.Model) (tableID string, content string, property *listener.Property) {
		row := m.AsRow()
		tableID = row.String(`id`, `-1`)
		content = row.String(fieldName)
		property = listener.NewPropertyWith(m, db.Cond{`id`: row.Get(`id`, `-1`)})
		return
	}
}

func (a *Listeners) Add(fieldName string, callback func(m factory.Model) (tableID string, content string, property *listener.Property)) fileupdater.Listener {
	if callback == nil {
		callback = GenDefaultCallback(fieldName)
	}
	(*a)[fieldName] = callback
	return a
}

func (a *Listeners) AddDefaultCallback(fieldNames ...string) fileupdater.Listener {
	for _, fieldName := range fieldNames {
		(*a)[fieldName] = GenDefaultCallback(fieldName)
	}
	return a
}

func (a *Listeners) Delete(fieldNames ...string) fileupdater.Listener {
	for _, fieldName := range fieldNames {
		if _, ok := (*a)[fieldName]; ok {
			delete(*a, fieldName)
		}
	}
	return a
}
