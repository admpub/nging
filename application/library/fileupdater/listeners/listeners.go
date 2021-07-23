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
	"github.com/admpub/nging/v3/application/library/fileupdater"
	"github.com/admpub/nging/v3/application/library/fileupdater/listener"
)

func New(table ...string) fileupdater.Listener {
	l := &Listeners{
		options: map[string]*fileupdater.Options{},
	}
	if len(table) > 0 {
		l.SetTableName(table[0])
	}
	return l
}

type Listeners struct {
	options   map[string]*fileupdater.Options
	tableName string
}

func (a *Listeners) SetTableName(tableName string) fileupdater.Listener {
	a.tableName = tableName
	return a
}

func (a *Listeners) Listen() {
	for _field, _option := range a.options {
		listener.NewWithOptions(_option).ListenDefault()
		delete(a.options, _field)
	}
}

func (a *Listeners) BuildOptions(fieldName string, setters ...fileupdater.OptionSetter) *fileupdater.Options {
	options := &fileupdater.Options{
		FieldName: fieldName,
		TableName: a.tableName,
	}
	for _, setter := range setters {
		setter(options)
	}
	if options.Callback == nil {
		options.Callback = fileupdater.GenCallbackDefault(fieldName, options.FieldValue)
	}
	if len(options.TableName) == 0 {
		panic(`liseners.options: TableName is empty`)
	}
	return options
}

func (a *Listeners) ListenByOptions(fieldName string, setters ...fileupdater.OptionSetter) {
	options := a.BuildOptions(fieldName, setters...)
	listener.NewWithOptions(options).ListenDefault()
}

func (a *Listeners) Add(fieldName string, setters ...fileupdater.OptionSetter) fileupdater.Listener {
	options := a.BuildOptions(fieldName, setters...)
	a.options[fieldName] = options
	return a
}

func (a *Listeners) Delete(fieldNames ...string) fileupdater.Listener {
	for _, fieldName := range fieldNames {
		if _, ok := a.options[fieldName]; ok {
			delete(a.options, fieldName)
		}
	}
	return a
}
