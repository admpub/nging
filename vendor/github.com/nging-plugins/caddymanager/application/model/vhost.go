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

	"github.com/admpub/null"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/nging-plugins/caddymanager/application/library/cmder"
)

func NewVhost(ctx echo.Context) *Vhost {
	return &Vhost{
		NgingVhost: dbschema.NewNgingVhost(ctx),
	}
}

type VhostAndGroup struct {
	*dbschema.NgingVhost
	Group        *dbschema.NgingVhostGroup `db:"-,relation=id:group_id|gtZero"`
	ServerName   null.String               `db:"serverName" json:",omitempty" xml:",omitempty"`
	ServerEngine null.String               `db:"serverEngine" json:",omitempty" xml:",omitempty"`
}

type Vhost struct {
	*dbschema.NgingVhost
}

func (m *Vhost) RemoveCachedCert() {
	if m.ServerIdent != `default` {
		return
	}
	caddyCfg := cmder.GetCaddyConfig()
	for _, domain := range strings.Split(m.Domain, ` `) {
		domain = strings.TrimSpace(domain)
		if len(domain) == 0 {
			continue
		}
		parts := strings.SplitN(domain, `//`, 2)
		if len(parts) == 2 {
			domain = parts[1]
		} else {
			domain = parts[0]
		}
		if len(domain) == 0 {
			continue
		}
		domain, _ = com.SplitHostPort(domain)
		if len(domain) == 0 {
			continue
		}
		caddyCfg.RemoveCachedCert(domain)
	}
}

func (f *Vhost) check() error {
	ctx := f.Context()
	f.Name = strings.TrimSpace(f.Name)
	if len(f.ServerIdent) == 0 {
		return ctx.NewError(code.InvalidParameter, `请选择引擎配置`).SetZone(`serverIdent`)
	}
	if !com.IsAlphaNumericUnderscoreHyphen(f.ServerIdent) {
		return ctx.NewError(code.InvalidParameter, `引擎配置参数无效`).SetZone(`serverIdent`)
	}
	return nil
}

func (f *Vhost) Add() (interface{}, error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	return f.NgingVhost.Insert()
}

func (f *Vhost) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = f.check(); err != nil {
		return err
	}
	return f.NgingVhost.Update(mw, args...)
}

func (f *Vhost) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	return f.NgingVhost.Delete(mw, args...)
}
