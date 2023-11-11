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

package tagfast

import (
	"reflect"
	"sync"
)

type Faster interface {
	Get(key string) string
	Parsed(key string, fns ...func() interface{}) interface{}
	GetParsed(key string, fns ...func(string) interface{}) interface{}
	SetParsed(key string, value interface{}) bool
}

func New(tag reflect.StructTag) Faster {
	return &tagFast{
		tag:    tag,
		cached: ParseStructTag(string(tag)),
		parsed: map[string]interface{}{},
	}
}

type tagFast struct {
	tag    reflect.StructTag      //example: `tagA:"valA" tagB:"valB" tagC:"a,b,c"`
	cached map[string]string      //example: {"tagA":"valA","tagB":"valB","tagC":"a,b,c"}
	parsed map[string]interface{} //example: {"tagC":["a","b","c"]}
	mu     sync.RWMutex
}

func (a *tagFast) Get(key string) string {
	a.mu.RLock()
	v, _ := a.cached[key]
	a.mu.RUnlock()
	return v
}

func (a *tagFast) Parsed(key string, fns ...func() interface{}) interface{} {
	v, ok := a.getParsed(key)
	if ok {
		return v
	}
	if len(fns) > 0 {
		fn := fns[0]
		if fn != nil {
			v = fn()
			a.SetParsed(key, v)
		}
	}
	return v
}

func (a *tagFast) getParsed(key string) (interface{}, bool) {
	a.mu.RLock()
	v, ok := a.parsed[key]
	a.mu.RUnlock()
	return v, ok
}

func (a *tagFast) GetParsed(key string, fns ...func(string) interface{}) interface{} {
	v, ok := a.getParsed(key)
	if ok {
		return v
	}
	if len(fns) > 0 {
		fn := fns[0]
		if fn != nil {
			v = fn(a.Get(key))
			a.SetParsed(key, v)
		}
	}
	return v
}

func (a *tagFast) SetParsed(key string, value interface{}) bool {
	a.mu.Lock()
	a.parsed[key] = value
	a.mu.Unlock()
	return true
}
