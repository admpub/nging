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

package common

import (
	"strings"

	"github.com/webx-top/echo"
)

// Adder interface
type Adder interface {
	Set(interface{}, ...interface{})
	Add() (interface{}, error)
}

// BatchAdd 批量添加(常用于批量添加分类)
func BatchAdd(ctx echo.Context, field string, adder Adder, before func(string) error) (added []string, err error) {
	values := strings.Split(ctx.Formx(field).String(), "\n")
	added = []string{}
	errs := []string{}
	statucField := strings.Title(field)
	for _, value := range values {
		value = strings.TrimSpace(value)
		if len(value) == 0 {
			continue
		}
		if before != nil {
			if err = before(value); err != nil {
				errs = append(errs, err.Error())
				continue
			}
		}
		adder.Set(statucField, value)
		_, err = adder.Add()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		added = append(added, value)
	}
	if len(errs) > 0 {
		err = ctx.E(strings.Join(errs, "\n"))
	}
	return added, err
}
