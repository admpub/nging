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
package tagfast

import (
	"reflect"
	"strconv"
	"sync"
)

var (
	lock   = new(sync.RWMutex)
	caches = make(map[string]map[string]*tagFast) //{"pkg.struct":{"field":Faster}}
)

func Tag(t reflect.Type, f reflect.StructField, tagName string) (value string, faster Faster) {
	faster = tag(t, f)
	if faster == nil {
		return
	}
	value = faster.Get(tagName)
	return
}

func tag(t reflect.Type, f reflect.StructField) Faster {
	if len(f.Tag) == 0 {
		return nil
	}
	lock.RLock()
	name := t.PkgPath() + "." + t.Name()
	var fast *tagFast
	if cc, ok := caches[name]; ok {
		if tf, ok := cc[f.Name]; ok {
			fast = tf
		} else {
			caches[name][f.Name] = nil
		}
	} else {
		caches[name] = make(map[string]*tagFast)
	}
	if fast == nil {
		fast = &tagFast{tag: f.Tag}
		caches[name][f.Name] = fast
	}
	lock.RUnlock()
	return fast
}

func Parsed(t reflect.Type, f reflect.StructField, tagName string, fns ...func() interface{}) interface{} {
	faster := tag(t, f)
	if faster == nil {
		return nil
	}
	return faster.Parsed(tagName, fns...)
}

func GetParsed(t reflect.Type, f reflect.StructField, tagName string, fns ...func(string) interface{}) interface{} {
	faster := tag(t, f)
	if faster == nil {
		return nil
	}
	return faster.GetParsed(tagName, fns...)
}

func Value(t reflect.Type, f reflect.StructField, tagName string) (value string) {
	value, _ = Tag(t, f, tagName)
	return
}

func Caches() map[string]map[string]*tagFast {
	return caches
}

type Faster interface {
	Get(key string) string
	Parsed(key string, fns ...func() interface{}) interface{}
	GetParsed(key string, fns ...func(string) interface{}) interface{}
	SetParsed(key string, value interface{}) bool
}

type tagFast struct {
	tag    reflect.StructTag      //example: `tagA:"valA" tagB:"valB" tagC:"a,b,c"`
	cached map[string]string      //example: {"tagA":"valA","tagB":"valB","tagC":"a,b,c"}
	parsed map[string]interface{} //example: {"tagC":["a","b","c"]}
}

func (a *tagFast) Get(key string) string {
	if a.cached == nil {
		a.cached = ParseStructTag(string(a.tag))
	}
	lock.RLock()
	defer lock.RUnlock()
	if v, ok := a.cached[key]; ok {
		return v
	}
	return ""
}

func (a *tagFast) Parsed(key string, fns ...func() interface{}) interface{} {
	if a.parsed == nil {
		a.parsed = make(map[string]interface{})
	}
	lock.RLock()
	if v, ok := a.parsed[key]; ok {
		lock.RUnlock()
		return v
	}
	lock.RUnlock()
	if len(fns) > 0 {
		fn := fns[0]
		if fn != nil {
			v := fn()
			a.SetParsed(key, v)
			return v
		}
	}
	return nil
}

func (a *tagFast) GetParsed(key string, fns ...func(string) interface{}) interface{} {
	if a.parsed == nil {
		a.parsed = make(map[string]interface{})
	}
	lock.RLock()
	if v, ok := a.parsed[key]; ok {
		lock.RUnlock()
		return v
	}
	lock.RUnlock()
	if len(fns) > 0 {
		fn := fns[0]
		if fn != nil {
			v := fn(a.Get(key))
			a.SetParsed(key, v)
			return v
		}
	}
	return nil
}

func (a *tagFast) SetParsed(key string, value interface{}) bool {
	if a.parsed == nil {
		a.parsed = make(map[string]interface{})
	}
	lock.Lock()
	defer lock.Unlock()
	a.parsed[key] = value
	return true
}

func ParseStructTag(tag string) map[string]string {
	lock.Lock()
	defer lock.Unlock()
	var tagsArray map[string]string = make(map[string]string)
	for tag != "" {
		// skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// scan to colon.
		// a space or a quote is a syntax error
		i = 0
		for i < len(tag) && tag[i] != ' ' && tag[i] != ':' && tag[i] != '"' {
			i++
		}
		if i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// scan quoted string to find value
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, _ := strconv.Unquote(qvalue)
		tagsArray[name] = value
	}
	return tagsArray
}
