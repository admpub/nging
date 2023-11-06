/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// Adder interface
type Adder interface {
	Set(interface{}, ...interface{})
	Add() (interface{}, error)
}

// BatchAdd 批量添加(常用于批量添加分类)
// BatchAdd(ctx, `ident,>name`, adder, before, `=`)
// 可以通过在字段名称前面添加“>”前缀来指定表单字段名称，如果不指定则默认使用第一个作为表单字段名
func BatchAdd(ctx echo.Context, field string, adder Adder, before func(*string) error, seperators ...string) (added []string, err error) {
	errs := []string{}
	field = strings.TrimSpace(field)
	var structFields []string
	var seperator string
	if len(seperators) > 0 {
		seperator = seperators[0]
	}
	var formField string
	var firstField string
	for _, v := range strings.Split(field, ",") {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		if strings.HasPrefix(v, `>`) {
			formField = strings.TrimPrefix(v, `>`)
			v = formField
		}
		if len(firstField) == 0 {
			firstField = v
		}
		structFields = append(structFields, com.Title(v))
	}
	if len(structFields) < 1 {
		return
	}
	if len(formField) == 0 {
		formField = firstField
	}
	values := strings.Split(ctx.Formx(formField).String(), "\n")
	for _, value := range values {
		value = strings.TrimSpace(value)
		if len(value) == 0 {
			continue
		}
		if before != nil {
			if err = before(&value); err != nil {
				errs = append(errs, err.Error())
				continue
			}
		}
		if len(seperator) > 0 {
			arr := strings.SplitN(value, seperator, len(structFields))
			for k, v := range arr {
				if k < len(structFields) {
					adder.Set(structFields[k], v)
				} else {
					break
				}
			}
		} else {
			adder.Set(structFields[0], value)
		}
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
	return
}
