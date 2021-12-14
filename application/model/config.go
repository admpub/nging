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
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/registry/settings"
)

func NewConfig(ctx echo.Context) *Config {
	m := &Config{
		NgingConfig: dbschema.NewNgingConfig(ctx),
	}
	return m
}

type Config struct {
	*dbschema.NgingConfig
}

func (f *Config) Upsert() (pk interface{}, err error) {
	m := dbschema.NewNgingConfig(f.Context())
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
		return f.NgingConfig.Add()
	}
	err = f.NgingConfig.Edit(nil, condition)
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
	return f.NgingConfig.Add()
}

func (f *Config) EditByPK(mw func(db.Result) db.Result, group string, key string) error {
	condition := db.And(
		db.Cond{`key`: key},
		db.Cond{`group`: group},
	)
	return f.NgingConfig.Edit(mw, condition)
}

func (f *Config) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.NgingConfig.Edit(mw, args...)
}

func (f *Config) ListByGroup(group string) (func() int64, error) {
	return f.NgingConfig.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`)
	}, 0, -1, `group`, group)
}

func (f *Config) ListMapByGroup(group string) (echo.H, error) {
	errs := common.NewErrors()
	_, err := f.ListByGroup(group)
	if err != nil {
		errs.Add(err)
		return nil, errs
	}
	cfg := echo.H{}
	decoder := settings.GetDecoder(group)
	for _, v := range f.Objects() {
		cfg, err = settings.DecodeConfig(v, cfg, decoder)
		if err != nil {
			err = fmt.Errorf(`[key:%s] %w`, v.Key, err)
			errs.Add(err)
		}
	}
	return cfg, errs.ToError()
}

func (f *Config) ListAllMapByGroup() (echo.H, error) {
	_, err := f.NgingConfig.ListByOffset(nil, func(r db.Result) db.Result {
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
		cfg.GetStore(v.Group).Set(v.Key, v)
	}
	return cfg, err
}
