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
	"github.com/webx-top/echo"
)

// Wrap adds several routes from package `net/http/pprof` to *gin.Engine object
func Wrap(router *echo.Echo) {
	router.Get("/debug/pprof/", IndexHandler())
	router.Get("/debug/pprof/heap", HeapHandler())
	router.Get("/debug/pprof/goroutine", GoroutineHandler())
	router.Get("/debug/pprof/block", BlockHandler())
	router.Get("/debug/pprof/threadcreate", ThreadCreateHandler())
	router.Get("/debug/pprof/cmdline", CmdlineHandler())
	router.Get("/debug/pprof/profile", ProfileHandler())
	router.Get("/debug/pprof/symbol", SymbolHandler())
	router.Get("/debug/pprof/trace", TraceHandler())
}

// Wrapper make sure we are backward compatible
var Wrapper = Wrap

// IndexHandler will pass the call from /debug/pprof to pprof
func IndexHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		Index(ctx)
		return nil
	}
}

// HeapHandler will pass the call from /debug/pprof/heap to pprof
func HeapHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		Handler("heap").ServeHTTP(ctx)
		return nil
	}
}

// GoroutineHandler will pass the call from /debug/pprof/goroutine to pprof
func GoroutineHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		Handler("goroutine").ServeHTTP(ctx)
		return nil
	}
}

// BlockHandler will pass the call from /debug/pprof/block to pprof
func BlockHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		Handler("block").ServeHTTP(ctx)
		return nil
	}
}

// ThreadCreateHandler will pass the call from /debug/pprof/threadcreate to pprof
func ThreadCreateHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		Handler("threadcreate").ServeHTTP(ctx)
		return nil
	}
}

// CmdlineHandler will pass the call from /debug/pprof/cmdline to pprof
func CmdlineHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		Cmdline(ctx)
		return nil
	}
}

// ProfileHandler will pass the call from /debug/pprof/profile to pprof
func ProfileHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		Profile(ctx)
		return nil
	}
}

// SymbolHandler will pass the call from /debug/pprof/symbol to pprof
func SymbolHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		Symbol(ctx)
		return nil
	}
}

// TraceHandler will pass the call from /debug/pprof/trace to pprof
func TraceHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		Trace(ctx)
		return nil
	}
}
