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

import "github.com/webx-top/echo/param"

var (
	_ ICore = &Echo{}
	_ ICore = &Group{}
)

type RouteRegister interface {
	MiddlewareRegister
	Group(prefix string, middleware ...interface{}) *Group
	Any(path string, h interface{}, middleware ...interface{}) IRouter
	Route(methods string, path string, h interface{}, middleware ...interface{}) IRouter
	Match(methods []string, path string, h interface{}, middleware ...interface{}) IRouter
	Connect(path string, h interface{}, m ...interface{}) IRouter
	Delete(path string, h interface{}, m ...interface{}) IRouter
	Get(path string, h interface{}, m ...interface{}) IRouter
	Head(path string, h interface{}, m ...interface{}) IRouter
	Options(path string, h interface{}, m ...interface{}) IRouter
	Patch(path string, h interface{}, m ...interface{}) IRouter
	Post(path string, h interface{}, m ...interface{}) IRouter
	Put(path string, h interface{}, m ...interface{}) IRouter
	Trace(path string, h interface{}, m ...interface{}) IRouter
	Static(prefix, root string)
	File(path, file string)
	Prefix() string
}

type ContextRegister interface {
	SetContext(Context)
}

type RendererRegister interface {
	SetRenderer(Renderer)
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
	URLBuilder
	RendererRegister
	Prefixer
}

type IRouter interface {
	SetName(string) IRouter
	GetName() string
	SetMeta(param.Store) IRouter
	SetMetaKV(string, interface{}) IRouter
	GetMeta() param.Store
}

type Closer interface {
	Close() error
}

type Prefixer interface {
	Prefix() string
}
