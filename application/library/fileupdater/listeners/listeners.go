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

	"github.com/webx-top/db/lib/factory"

	"github.com/admpub/nging/application/library/fileupdater/listener"
)

func New() *Listeners {
	return &Listeners{}
}

type Listeners map[string]func(m factory.Model) (tableID string, content string, property *listener.Property)

func (a *Listeners) Listen(tableName string, embedded bool, seperatorAndSameFields ...string) *Listeners {
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
	return a
}

func (a *Listeners) ListenByField(fieldNames string, tableName string, embedded bool, seperatorAndSameFields ...string) *Listeners {
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
	return a
}

func (a *Listeners) Add(fieldName string, callback func(m factory.Model) (tableID string, content string, property *listener.Property)) *Listeners {
	(*a)[fieldName] = callback
	return a
}

func (a *Listeners) Delete(fieldNames ...string) *Listeners {
	for _, fieldName := range fieldNames {
		if _, ok := (*a)[fieldName]; ok {
			delete(*a, fieldName)
		}
	}
	return a
}
