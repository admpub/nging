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
	"strings"

	"github.com/webx-top/echo"
)

func NewRegister(e *echo.Echo, groupNamers ...func(string) string) *Register {
	return &Register{
		Echo:              e,
		Handlers:          []func(echo.RouteRegister){},
		GroupHandlers:     map[string][]func(echo.RouteRegister){},
		GroupMiddlewares:  map[string][]interface{}{},
		PrefixMiddlewares: map[string][]interface{}{},
		GroupNamers:       groupNamers,
	}
}

type Register struct {
	Echo              *echo.Echo
	RootGroup         string
	Handlers          []func(echo.RouteRegister)
	GroupHandlers     map[string][]func(echo.RouteRegister)
	GroupNamers       []func(string) string
	GroupMiddlewares  map[string][]interface{}
	PrefixMiddlewares map[string][]interface{}
}

func (r *Register) AddGroupNamer(namers ...func(string) string) {
	r.GroupNamers = append(r.GroupNamers, namers...)
}

func (r *Register) SetGroupNamer(namers ...func(string) string) {
	r.GroupNamers = namers
}

func (r *Register) Apply() {
	e := r.Echo
	for _, register := range r.Handlers {
		register(e)
	}
	var groupDefaultMiddlewares []interface{}
	middlewares, ok := r.GroupMiddlewares[`*`]
	if ok {
		groupDefaultMiddlewares = append(groupDefaultMiddlewares, middlewares...)
	}
	for group, handlers := range r.GroupHandlers {
		for _, namer := range r.GroupNamers {
			group = namer(group)
		}
		g := e.Group(group)
		if group != r.RootGroup { // 组名为空时，为顶层组
			g.Use(groupDefaultMiddlewares...)
		}
		for prefix, middlewares := range r.PrefixMiddlewares {
			if strings.HasPrefix(group, prefix) {
				g.Use(middlewares...)
			}
		}
		middlewares, ok := r.GroupMiddlewares[group]
		if ok {
			g.Use(middlewares...)
		}
		for _, register := range handlers {
			register(g)
		}
	}
}

func (r *Register) Use(groupName string, middlewares ...interface{}) {
	if groupName != `*` && strings.HasSuffix(groupName, `*`) {
		groupName = strings.TrimRight(groupName, `*`)
		if _, ok := r.PrefixMiddlewares[groupName]; !ok {
			r.PrefixMiddlewares[groupName] = []interface{}{}
		}
		r.PrefixMiddlewares[groupName] = append(r.PrefixMiddlewares[groupName], middlewares...)
		return
	}
	if _, ok := r.GroupMiddlewares[groupName]; !ok {
		r.GroupMiddlewares[groupName] = []interface{}{}
	}
	r.GroupMiddlewares[groupName] = append(r.GroupMiddlewares[groupName], middlewares...)
}

func (r *Register) Register(fn func(echo.RouteRegister)) {
	r.Handlers = append(r.Handlers, fn)
}

func (r *Register) RegisterToGroup(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
	_, ok := r.GroupHandlers[groupName]
	if !ok {
		r.GroupHandlers[groupName] = []func(echo.RouteRegister){}
	}
	if len(middlewares) > 0 {
		_, ok = r.GroupMiddlewares[groupName]
		if !ok {
			r.GroupMiddlewares[groupName] = []interface{}{}
		}
		r.GroupMiddlewares[groupName] = append(r.GroupMiddlewares[groupName], middlewares...)
	}
	r.GroupHandlers[groupName] = append(r.GroupHandlers[groupName], fn)
}
