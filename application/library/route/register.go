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

package route

import (
	"github.com/webx-top/echo"
)

func NewRegister(e *echo.Echo, groupNamers ...func(string) string) *Register {
	return &Register{
		Echo:     e,
		Handlers: []func(echo.RouteRegister){},
		Group:    NewGroup(groupNamers...),
		Hosts:    make(map[string]*Host),
	}
}

type Register struct {
	Echo      *echo.Echo
	RootGroup string
	Handlers  []func(echo.RouteRegister)
	Group     *Group
	Hosts     map[string]*Host
}

func (r *Register) AddGroupNamer(namers ...func(string) string) {
	r.Group.AddNamer(namers...)
}

func (r *Register) SetGroupNamer(namers ...func(string) string) {
	r.Group.SetNamer(namers...)
}

func (r *Register) Apply() {
	e := r.Echo
	for _, register := range r.Handlers {
		register(e)
	}
	r.Group.Apply(e, r.RootGroup)
	for _, host := range r.Hosts {
		hst := e.Host(host.Name, host.Middlewares...)
		if len(host.Alias) > 0 {
			hst.SetAlias(host.Alias)
		}
		host.Group.Apply(hst, r.RootGroup)
	}
}

func (r *Register) Use(groupName string, middlewares ...interface{}) {
	r.Group.Use(groupName, middlewares...)
}

func (r *Register) Register(fn func(echo.RouteRegister)) {
	r.Handlers = append(r.Handlers, fn)
}

func (r *Register) RegisterToGroup(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
	r.Group.Register(groupName, fn, middlewares...)
}

func (r *Register) Host(hostName string, middlewares ...interface{}) *Host {
	host, ok := r.Hosts[hostName]
	if !ok {
		host = &Host{
			Name:  hostName,
			Group: NewGroup(),
		}
		r.Hosts[hostName] = host
	}
	host.Use(middlewares...)
	return host
}
