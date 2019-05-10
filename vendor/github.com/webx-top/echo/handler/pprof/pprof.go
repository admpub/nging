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
package pprof

import (
	"net/http/pprof"

	"github.com/webx-top/echo"
)

var typeNames = []string{`heap`, `goroutine`, `block`, `threadcreate`, `allocs`, `mutex`}

// Wrap adds several routes from package `net/http/pprof` to *gin.Engine object
func Wrap(router *echo.Echo) {

	router.Get("/debug/pprof/", pprof.Index)
	for _, name := range typeNames {
		router.Get("/debug/pprof/"+name, pprof.Handler(name))
	}
	router.Get("/debug/pprof/cmdline", pprof.Cmdline)
	router.Get("/debug/pprof/profile", pprof.Profile)
	router.Get("/debug/pprof/symbol", pprof.Symbol)
	router.Get("/debug/pprof/trace", pprof.Trace)
}

// Wrapper make sure we are backward compatible
var Wrapper = Wrap
