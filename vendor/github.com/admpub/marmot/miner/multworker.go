// 
// 	Copyright 2017 by marmot author: gdccmcm14@live.com.
// 	Licensed under the Apache License, Version 2.0 (the "License");
// 	you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
// 		http://www.apache.org/licenses/LICENSE-2.0
// 	Unless required by applicable law or agreed to in writing, software
// 	distributed under the License is distributed on an "AS IS" BASIS,
// 	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// 	See the License for the specific language governing permissions and
// 	limitations under the License
//

package miner

import (
	"sync"
)

// Pool for many Worker, every Worker can only serial execution
var Pool = &_Workers{ws: make(map[string]*Worker)}

type _Workers struct {
	mux sync.RWMutex // simple lock
	ws  map[string]*Worker
}

func (pool *_Workers) Get(name string) (b *Worker, ok bool) {
	pool.mux.RLock()
	b, ok = pool.ws[name]
	pool.mux.RUnlock()
	return
}

func (pool *_Workers) Set(name string, b *Worker) {
	pool.mux.Lock()
	pool.ws[name] = b
	pool.mux.Unlock()
	return
}

func (pool *_Workers) Delete(name string) {
	pool.mux.Lock()
	delete(pool.ws, name)
	pool.mux.Unlock()
	return
}
