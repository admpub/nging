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

package model

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
)

func NewKv(ctx echo.Context) *Kv {
	m := &Kv{
		NgingKv: &dbschema.NgingKv{},
		base:    base.New(ctx),
	}
	m.SetContext(ctx)
	return m
}

type Kv struct {
	*dbschema.NgingKv
	base *base.Base
}

func (s *Kv) check() error {
	s.Key = strings.TrimSpace(s.Key)
	if len(s.Key) == 0 {
		return s.base.E(`键不能为空`)
	}
	s.Value = strings.TrimSpace(s.Value)
	if len(s.Value) == 0 {
		return s.base.E(`值不能为空`)
	}
	s.Type = strings.TrimSpace(s.Type)
	if len(s.Type) == 0 {
		return s.base.E(`类型不能为空`)
	}
	return nil
}

func (s *Kv) Get(mw func(db.Result) db.Result, args ...interface{}) error {
	err := s.NgingKv.Get(mw, args...)
	if err != nil {
		return err
	}
	return nil
}

func (s *Kv) Add() (pk interface{}, err error) {
	if err = s.check(); err != nil {
		return nil, err
	}
	return s.NgingKv.Add()
}

func (s *Kv) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = s.check(); err != nil {
		return err
	}
	return s.NgingKv.Edit(mw, args...)
}

func (s *Kv) SetSingleField(id int, field string, value string) error {
	set := echo.H{}
	switch field {
	case "value", "key", "sort", "child_key_type":
		set[field] = value
	default:
		return s.base.E(`不支持修改字段: %v`, field)
	}
	return s.SetFields(nil, set, `id`, id)
}
