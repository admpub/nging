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
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
)

const KvRootType = `root`

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
	var (
		exists bool
		err    error
	)
	if s.Id > 0 { // edit
		exists, err = s.Exists(nil, db.And(
			db.Cond{`key`: s.Key},
			db.Cond{`type`: s.Type},
			db.Cond{`id`: db.NotEq(s.Id)},
		))
	} else {
		exists, err = s.Exists(nil, db.And(
			db.Cond{`key`: s.Key},
			db.Cond{`type`: s.Type},
		))
	}
	if err != nil {
		return err
	}
	if exists {
		return s.base.E(`键"%v"已经存在`, s.Key)
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
	s.NgingKv.Updated = uint(time.Now().Unix())
	return s.NgingKv.Add()
}

func (s *Kv) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = s.check(); err != nil {
		return err
	}
	return s.NgingKv.Edit(mw, args...)
}

func (s *Kv) IsRootType(typ string) bool {
	return typ == KvRootType
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

func (s *Kv) KvTypeList(excludeIDs ...uint) []*dbschema.NgingKv {
	cond := db.NewCompounds()
	cond.AddKV(`type`, KvRootType)
	if len(excludeIDs) > 0 && excludeIDs[0] > 0 {
		cond.AddKV(`id`, db.NotEq(excludeIDs[0]))
	}
	_, err := s.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`)
	}, 0, -1, cond.And())
	if err == nil {
		return s.Objects()
	}
	return nil
}
