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

package upload

import (
	"io"
)

var clients = make(map[string]func() Client)
var defaults = func() Client {
	return New(nil)
}

func Register(name string, c func() Client) {
	clients[name] = c
}

func Get(name string) Client {
	fn, ok := clients[name]
	if !ok || fn == nil {
		fn = defaults
	}
	return fn()
}

func Has(name string) bool {
	_, ok := clients[name]
	return ok
}

func Delete(name string) {
	delete(clients, name)
}

type Storer interface {
	Put(dstFile string, body io.Reader, size int64) (savePath string, fileURL string, err error)
	Delete(file string) error
}
