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
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/registry/settings"
	"github.com/admpub/nging/application/model/base"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewConfig(ctx echo.Context) *Config {
	return &Config{
		Config: &dbschema.Config{},
		base:   base.New(ctx),
	}
}

type Config struct {
	*dbschema.Config
	base *base.Base
}

func (f *Config) Upsert() (pk interface{}, err error) {
	m := &dbschema.Config{}
	condition := db.And(
		db.Cond{`key`: f.Key},
		db.Cond{`group`: f.Group},
	)
	n, err := m.Count(nil, condition)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
	}
	if n == 0 {
		return f.Config.Add()
	}
	err = f.Config.Edit(nil, condition)
	return
}

func (f *Config) ValueByKey(group, key string) string {
	condition := db.And(
		db.Cond{`key`: key},
		db.Cond{`group`: group},
	)
	err := f.Get(nil, condition)
	if err != nil {
		return err.Error()
	}
	return f.Value
}

func (f *Config) Add() (pk interface{}, err error) {
	return f.Config.Add()
}

func (f *Config) EditByPK(mw func(db.Result) db.Result, group string, key string) error {
	condition := db.And(
		db.Cond{`key`: key},
		db.Cond{`group`: group},
	)
	return f.Config.Edit(mw, condition)
}

func (f *Config) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.Config.Edit(mw, args...)
}

func (f *Config) ListByGroup(group string) (func() int64, error) {
	return f.Config.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`)
	}, 0, -1, `group`, group)
}

func (f *Config) ListMapByGroup(group string) (echo.H, error) {
	_, err := f.ListByGroup(group)
	if err != nil {
		return nil, err
	}
	cfg := echo.H{}
	decoder := settings.GetDecoder(group)
	for _, v := range f.Objects() {
		cfg, err = settings.DecodeConfig(v, cfg, decoder)
		if err != nil {
			return cfg, err
		}
	}
	return cfg, err
}

func (f *Config) ListAllMapByGroup() (echo.H, error) {
	_, err := f.Config.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`)
	}, 0, -1)
	if err != nil {
		return nil, err
	}
	cfg := echo.H{}
	for _, v := range f.Objects() {
		if _, _y := cfg[v.Group]; !_y {
			cfg[v.Group] = echo.H{}
		}
		cfg.Store(v.Group).Set(v.Key, v)
	}
	return cfg, err
}
