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
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/model/base"
)

func NewCloudStorage(ctx echo.Context) *CloudStorage {
	m := &CloudStorage{
		CloudStorage: &dbschema.CloudStorage{},
		base:         base.New(ctx),
	}
	m.SetContext(ctx)
	return m
}

type CloudStorage struct {
	*dbschema.CloudStorage
	base *base.Base
}

func (s *CloudStorage) check() error {
	s.Bucket = strings.TrimSpace(s.Bucket)
	s.Endpoint = strings.TrimSpace(s.Endpoint)
	s.Secret = strings.TrimSpace(s.Secret)
	s.Secret = common.Crypto().Encode(s.Secret)
	return nil
}

func (s *CloudStorage) RawSecret() string {
	return common.Crypto().Decode(s.Secret)
}

func (s *CloudStorage) Get(mw func(db.Result) db.Result, args ...interface{}) error {
	err := s.CloudStorage.Get(mw, args...)
	if err != nil {
		return err
	}
	s.Secret = s.RawSecret()
	return nil
}

func (s *CloudStorage) Add() (pk interface{}, err error) {
	if err = s.check(); err != nil {
		return nil, err
	}
	return s.CloudStorage.Add()
}

func (s *CloudStorage) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = s.check(); err != nil {
		return err
	}
	return s.CloudStorage.Edit(mw, args...)
}
