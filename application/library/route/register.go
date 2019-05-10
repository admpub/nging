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

func NewRegister() *Register {
    return &Register{
        Handlers: []func(*echo.Echo){},
        GroupHandlers:map[string][]func(*echo.Group){},
        GroupMiddlewares:map[string][]interface{}{},
    }
}

type Register struct{
    Handlers         []func(*echo.Echo)
	GroupHandlers    map[string][]func(*echo.Group)
	GroupMiddlewares map[string][]interface{}
}

func (r *Register) Apply(e *echo.Echo){
	for _, register := range r.Handlers {
		register(e)
    }
    var groupDefaultMiddlewares []interface{}
    middlewares, ok := r.GroupMiddlewares[`*`]
    if ok {
        groupDefaultMiddlewares = append(groupDefaultMiddlewares, middlewares...)
    }
	for group, handlers := range r.GroupHandlers {
		g := e.Group(group, groupDefaultMiddlewares...)
		middlewares, ok := r.GroupMiddlewares[group]
		if ok {
			g.Use(middlewares...)
		}
		for _, register := range handlers {
			register(g)
		}
	}
}

func (r *Register) Use(groupName string,middlewares ...interface{}){
	if _,ok := r.GroupMiddlewares[groupName]; !ok {
		r.GroupMiddlewares[groupName]=[]interface{}{}
	}
	r.GroupMiddlewares[groupName]=append(r.GroupMiddlewares[groupName],middlewares...)
}

func (r *Register) Register(fn func(*echo.Echo)) {
	r.Handlers = append(r.Handlers, fn)
}

func (r *Register) RegisterToGroup(groupName string, fn func(*echo.Group), middlewares ...interface{}) {
	_, ok := r.GroupHandlers[groupName]
	if !ok {
		r.GroupHandlers[groupName] = []func(*echo.Group){}
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
