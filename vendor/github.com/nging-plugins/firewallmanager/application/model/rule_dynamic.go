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
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/firewallmanager/application/dbschema"
)

func NewRuleDynamic(ctx echo.Context) *RuleDynamic {
	return &RuleDynamic{
		NgingFirewallRuleDynamic: dbschema.NewNgingFirewallRuleDynamic(ctx),
	}
}

type RuleDynamic struct {
	*dbschema.NgingFirewallRuleDynamic
}

func (r *RuleDynamic) check() error {
	return nil
}

func (r *RuleDynamic) Add() (interface{}, error) {
	if err := r.check(); err != nil {
		return nil, err
	}
	return r.NgingFirewallRuleDynamic.Insert()
}

func (r *RuleDynamic) ExistsAvailable() (bool, error) {
	return r.NgingFirewallRuleDynamic.Exists(nil, `disabled`, `N`)
}

func (r *RuleDynamic) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := r.check(); err != nil {
		return err
	}
	return r.NgingFirewallRuleDynamic.Update(mw, args...)
}

func (r *RuleDynamic) ListPage(cond *db.Compounds, sorts ...interface{}) ([]*dbschema.NgingFirewallRuleDynamic, error) {
	err := r.NgingFirewallRuleDynamic.ListPage(cond, sorts...)
	if err != nil {
		return nil, err
	}
	return r.Objects(), nil
}
