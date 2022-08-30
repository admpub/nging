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
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/ipsimplefilter"

	"github.com/nging-plugins/frpmanager/application/dbschema"
)

var (
	ErrNoneDirectory      = errors.New(`No accessible directory`)
	ErrBannedUser         = errors.New(`The current account has been disabled`)
	ErrIPAddressIsBlocked = errors.New(`IP is blocked`)
)

func NewFrpUser(ctx echo.Context) *FrpUser {
	return &FrpUser{
		NgingFrpUser: dbschema.NewNgingFrpUser(ctx),
	}
}

type FrpUser struct {
	*dbschema.NgingFrpUser
}

func (f *FrpUser) Exists(serverID uint, username string, excludeIDs ...uint64) (bool, error) {
	cond := db.NewCompounds()
	cond.Add(db.Cond{`server_id`: serverID})
	cond.Add(db.Cond{`username`: username})
	if len(excludeIDs) > 0 {
		cond.Add(db.Cond{`id`: db.NotEq(excludeIDs[0])})
	}
	return f.NgingFrpUser.Exists(nil, cond.And())
}

func (f *FrpUser) CheckPasswd(serverID uint, username string, password string) error {
	err := f.NgingFrpUser.Get(nil, db.And(
		db.Cond{`server_id`: serverID},
		db.Cond{`username`: username},
	))
	if err != nil {
		if err == db.ErrNoMoreRows && serverID != 0 {
			err = f.NgingFrpUser.Get(nil, db.And(
				db.Cond{`server_id`: 0},
				db.Cond{`username`: username},
			))
		}
		if err != nil {
			return err
		}
	}
	salt := common.CookieConfig().BlockKey
	if f.Password != com.MakePassword(password, salt) {
		return f.Context().NewError(code.Failure, `密码不正确`)
	}
	if f.Banned == `Y` {
		return ErrBannedUser
	}
	ipAddr := f.Context().RealIP()
	ip := net.ParseIP(ipAddr)
	if len(f.IpWhitelist) > 0 {
		for _, row := range strings.Split(f.IpWhitelist, "\n") {
			row = strings.TrimSpace(row)
			if len(row) == 0 {
				continue
			}
			r := ipsimplefilter.New(row)
			if !r.Contains(ip) {
				return fmt.Errorf("%w: %v", ErrIPAddressIsBlocked, ipAddr)
			}
		}
	}
	if len(f.IpBlacklist) > 0 {
		for _, row := range strings.Split(f.IpBlacklist, "\n") {
			row = strings.TrimSpace(row)
			if len(row) == 0 {
				continue
			}
			r := ipsimplefilter.New(row)
			if r.Contains(ip) {
				return fmt.Errorf("%w: %v", ErrIPAddressIsBlocked, ipAddr)
			}
		}
	}
	now := time.Now().Unix()
	if f.Start > 0 && int64(f.Start) > now {
		return f.Context().NewError(code.DataProcessing, `账号尚未生效`)
	}
	if f.End > 0 && int64(f.End) < now {
		return f.Context().NewError(code.DataHasExpired, `账号已经过期`)
	}
	return err
}

func (f *FrpUser) check() error {
	if len(f.Username) == 0 {
		return f.Context().NewError(code.InvalidParameter, `用户名不能为空`).SetZone(`username`)
	}
	var exists bool
	var err error
	if f.Id > 0 {
		exists, err = f.Exists(f.ServerId, f.Username, f.Id)
	} else {
		exists, err = f.Exists(f.ServerId, f.Username)
	}
	if err != nil {
		return err
	}
	if exists {
		return f.Context().NewError(code.DataAlreadyExists, `用户名已经存在`).SetZone(`username`)
	}
	return err
}

func (f *FrpUser) Add() (pk interface{}, err error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	salt := common.CookieConfig().BlockKey
	f.Password = com.MakePassword(f.Password, salt)
	return f.NgingFrpUser.Insert()
}

func (f *FrpUser) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	old := dbschema.NewNgingFrpUser(f.Context())
	err := old.Get(func(r db.Result) db.Result {
		return r.Select(`password`)
	}, `id`, f.Id)
	if err != nil {
		return err
	}
	if len(f.Password) == 0 {
		f.Password = old.Password
	} else {
		salt := common.CookieConfig().BlockKey
		f.Password = com.MakePassword(f.Password, salt)
	}
	return f.NgingFrpUser.Update(mw, args...)
}
