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

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/dbschema"
)

func NewKv(ctx echo.Context) *Kv {
	m := &Kv{
		NgingKv: dbschema.NewNgingKv(ctx),
	}
	return m
}

type Kv struct {
	*dbschema.NgingKv
}

func (s *Kv) check() error {
	s.Key = strings.TrimSpace(s.Key)
	if len(s.Key) == 0 {
		return s.Context().E(`键不能为空`)
	}
	s.Value = strings.TrimSpace(s.Value)
	if len(s.Value) == 0 {
		return s.Context().E(`值不能为空`)
	}
	s.Type = strings.TrimSpace(s.Type)
	if len(s.Type) == 0 {
		return s.Context().E(`类型不能为空`)
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
		return s.Context().E(`键"%v"已经存在`, s.Key)
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

func (s *Kv) AutoCreateKey(key string, value ...string) error {
	m := dbschema.NewNgingKv(s.Context())
	m.Key = key
	m.Type = AutoCreatedType
	m.ChildKeyType = KvDefaultDataType
	m.Value = ``
	if len(value) > 0 {
		m.Value = value[0]
	}
	m.Updated = uint(time.Now().Unix())
	if _, err := m.Add(); err != nil {
		return err
	}
	m.Reset()
	err := m.Get(nil, `key`, AutoCreatedType)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		m.Key = AutoCreatedType
		m.Type = KvRootType
		m.ChildKeyType = KvDefaultDataType
		m.Value = `自动创建`
		_, err = m.Add()
	}
	return err
}

func (s *Kv) GetValue(key string, defaultValue ...string) (string, error) {
	err := s.NgingKv.Get(func(r db.Result) db.Result {
		return r.Select(`value`)
	}, db.And(
		db.Cond{`key`: key},
		db.Cond{`type`: db.NotEq(KvRootType)},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			if err = s.AutoCreateKey(key, defaultValue...); err != nil {
				s.Context().Logger().Error(err)
			}
		}
		if len(defaultValue) > 0 {
			return defaultValue[0], err
		}
		return ``, err
	}
	if len(defaultValue) > 0 && len(s.Value) == 0 {
		return defaultValue[0], err
	}
	return s.Value, err
}

func (s *Kv) GetTypeValues(typ string, defaultValue ...string) ([]string, error) {
	_, err := s.NgingKv.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`value`)
	}, 0, -1, db.Cond{`type`: typ})
	if err != nil {
		return defaultValue, err
	}
	rows := s.Objects()
	if len(rows) == 0 {
		return defaultValue, err
	}
	values := make([]string, len(rows))
	for index, row := range rows {
		values[index] = row.Value
	}
	return values, err
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

func (s *Kv) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	var rows []*dbschema.NgingKv
	s.NgingKv.ListByOffset(&rows, nil, 0, -1, args...)
	var types []string
	for _, row := range rows {
		if row.Type != KvRootType {
			continue
		}
		if com.InSlice(row.Key, types) {
			continue
		}
		types = append(types, row.Key)
	}
	if len(types) > 0 {
		err = s.NgingKv.Delete(nil, db.Cond{`type`: db.In(types)})
		if err != nil {
			return
		}
	}
	return s.NgingKv.Delete(mw, args...)
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
		return s.Context().E(`不支持修改字段: %v`, field)
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

func (s *Kv) ListByType(typ string, excludeIDs ...uint) []*dbschema.NgingKv {
	cond := db.NewCompounds()
	cond.AddKV(`type`, typ)
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

func (s *Kv) GetFromTypeList(typeList []*dbschema.NgingKv, key string) string {
	if key == KvRootType {
		return KvRootType
	}
	for _, row := range typeList {
		if row.Key == key {
			return row.Value
		}
	}
	return key
}

func (s *Kv) ListToMap(typeList []*dbschema.NgingKv) map[string]*dbschema.NgingKv {
	r := map[string]*dbschema.NgingKv{}
	for _, row := range typeList {
		r[row.Key] = row
	}
	return r
}
