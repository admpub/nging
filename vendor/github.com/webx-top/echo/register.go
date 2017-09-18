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

package echo

type RouteRegister interface {
	Any(path string, h interface{}, middleware ...interface{})
	Route(methods string, path string, h interface{}, middleware ...interface{})
	Match(methods []string, path string, h interface{}, middleware ...interface{})
	Connect(path string, h interface{}, m ...interface{})
	Delete(path string, h interface{}, m ...interface{})
	Get(path string, h interface{}, m ...interface{})
	Head(path string, h interface{}, m ...interface{})
	Options(path string, h interface{}, m ...interface{})
	Patch(path string, h interface{}, m ...interface{})
	Post(path string, h interface{}, m ...interface{})
	Put(path string, h interface{}, m ...interface{})
	Trace(path string, h interface{}, m ...interface{})
	Static(prefix, root string)
	File(path, file string)
}

type ContextRegister interface {
	SetContext(Context)
}

type MiddlewareRegister interface {
	Use(middleware ...interface{})
	Pre(middleware ...interface{})
}

type URLBuilder interface {
	URL(interface{}, ...interface{}) string
}

type ICore interface {
	RouteRegister
	MiddlewareRegister
	URLBuilder
}

type Closer interface {
	Close() error
}

type Prefixer interface {
	Prefix() string
}
