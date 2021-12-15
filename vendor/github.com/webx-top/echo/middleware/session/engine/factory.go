/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package engine

import (
	"encoding/base32"
	"strings"

	"github.com/admpub/securecookie"
	"github.com/admpub/sessions"
	"github.com/webx-top/echo"
)

func NewSession(ctx echo.Context) echo.Sessioner {
	options := ctx.SessionOptions()
	store := StoreEngine(options)
	return NewMySession(store, options.Name, ctx)
}

func NewMySession(store sessions.Store, name string, ctx echo.Context) echo.Sessioner {
	return &Session{
		name:    name,
		context: ctx,
		store:   store,
		session: nil,
		written: false,
	}
}

func StoreEngine(options *echo.SessionOptions) (store sessions.Store) {
	if options == nil {
		return nil
	}
	store = Get(options.Engine)
	if store == nil {
		if options.Engine != `cookie` {
			store = Get(`cookie`)
		}
	}
	return
}

func GenerateSessionID(prefix ...string) string {
	var _prefix string
	if len(prefix) > 0 {
		_prefix = prefix[0]
	}
	return _prefix + strings.TrimRight(
		base32.StdEncoding.EncodeToString(
			securecookie.GenerateRandomKey(32),
		),
		"=",
	)
}

var stores = map[string]sessions.Store{}

type Closer interface {
	Close() error
}

func Reg(name string, store sessions.Store) {
	if old, ok := stores[name]; ok {
		if c, ok := old.(Closer); ok {
			c.Close()
		}
	}
	stores[name] = store
}

func Get(name string) sessions.Store {
	if store, ok := stores[name]; ok {
		return store
	}
	return nil
}

func Del(name string) {
	if old, ok := stores[name]; ok {
		if c, ok := old.(Closer); ok {
			c.Close()
		}
		delete(stores, name)
	}
}
