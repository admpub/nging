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
package render

import (
	"github.com/webx-top/echo/logger"
	. "github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/render/standard"
)

var engines = make(map[string]func(string) Driver)

func New(key string, tmplDir string, args ...logger.Logger) Driver {
	if fn, ok := engines[key]; ok {
		return fn(tmplDir)
	}
	return standard.New(tmplDir, args...)
}

func Reg(key string, val func(string) Driver) {
	engines[key] = val
}

func Del(key string) {
	if _, ok := engines[key]; ok {
		delete(engines, key)
	}
}
