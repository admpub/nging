/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
	"github.com/webx-top/echo/logger"
)

func NewRegister(e *echo.Echo, groupNamers ...func(string) string) IRegister {
	return &Register{
		echo:  e,
		group: NewGroup(groupNamers...),
		hosts: make(map[string]*Host),
	}
}

type Register struct {
	echo           *echo.Echo
	prefix         string
	rootGroup      string
	handlers       []func(echo.RouteRegister)
	preMiddlewares []interface{}
	middlewares    []interface{}
	group          *Group
	hosts          map[string]*Host
}

func (r *Register) Echo() *echo.Echo {
	return r.echo
}

func (r *Register) Routes() []*echo.Route {
	return r.echo.Routes()
}

func (r *Register) Logger() logger.Logger {
	return r.echo.Logger()
}

func (r *Register) Prefix() string {
	return r.echo.Prefix()
}

func (r *Register) SetPrefix(prefix string) {
	r.prefix = prefix
}

func (r *Register) MetaHandler(m echo.H, handler interface{}, requests ...interface{}) echo.Handler {
	return r.echo.MetaHandler(m, handler, requests...)
}

func (r *Register) MetaHandlerWithRequest(m echo.H, handler interface{}, requests interface{}, methods ...string) echo.Handler {
	return r.echo.MetaHandlerWithRequest(m, handler, requests, methods...)
}

func (r *Register) HandlerWithRequest(handler interface{}, requests interface{}, methods ...string) echo.Handler {
	return r.echo.MetaHandlerWithRequest(nil, handler, requests, methods...)
}

func (r *Register) AddGroupNamer(namers ...func(string) string) {
	r.group.AddNamer(namers...)
}

func (r *Register) SetGroupNamer(namers ...func(string) string) {
	r.group.SetNamer(namers...)
}

func (r *Register) SetRootGroup(groupName string) {
	r.rootGroup = groupName
}

func (r *Register) RootGroup() string {
	return r.rootGroup
}

func (r *Register) Apply() {
	e := r.echo
	if len(r.prefix) > 0 {
		e.SetPrefix(r.prefix)
	}
	e.Pre(r.preMiddlewares...)
	e.Use(r.middlewares...)
	for _, register := range r.handlers {
		register(e)
	}
	r.group.Apply(e, r.rootGroup)
	for _, host := range r.hosts {
		hst := e.Host(host.Name, host.Middlewares...)
		if len(host.Alias) > 0 {
			hst.SetAlias(host.Alias)
		}
		host.Group.Apply(hst, r.rootGroup)
	}
}

func (r *Register) Pre(middlewares ...interface{}) {
	r.preMiddlewares = append(r.preMiddlewares, middlewares...)
}

func (r *Register) Use(middlewares ...interface{}) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *Register) PreToGroup(groupName string, middlewares ...interface{}) {
	r.group.Pre(groupName, middlewares...)
}

func (r *Register) UseToGroup(groupName string, middlewares ...interface{}) {
	r.group.Use(groupName, middlewares...)
}

func (r *Register) Register(fn func(echo.RouteRegister)) {
	r.handlers = append(r.handlers, fn)
}

func (r *Register) RegisterToGroup(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
	r.group.Register(groupName, fn, middlewares...)
}

func (r *Register) Host(hostName string, middlewares ...interface{}) *Host {
	host, ok := r.hosts[hostName]
	if !ok {
		host = &Host{
			Name:  hostName,
			Group: NewGroup(),
		}
		r.hosts[hostName] = host
	}
	host.Use(middlewares...)
	return host
}
