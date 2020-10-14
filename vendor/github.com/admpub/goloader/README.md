
# Goloader

[![Build Status](https://travis-ci.org/pkujhd/goloader.svg?branch=master)](https://travis-ci.org/pkujhd/goloader)

Goloader can load and run Golang code at runtime.

Forked from **https://github.com/dearplain/goloader**, Take over maintenance because the original author is not in maintenance

## How does it work?

Goloader works like a linker: it relocates the address of symbols in an object file, generates runnable code, and then reuses the runtime function and the type pointer of the loader.

Goloader provides some information to the runtime and gc of Go, which allows it to work correctly with them.

Please note that Goloader is not a scripting engine. It reads the output of Go compiler and makes them runnable. All features of Go are supported, and run just as fast and lightweight as native Go code.

## Comparison with plugin

Goloader reuses the Go runtime, which makes it much smaller. And code loaded by Goloader is unloadable.

Goloader is debuggable, and supports pprof tool(Yes, you can see code loaded by Goloader in pprof).

## Build

**Make sure you're using go >= 1.8.**

First, execute the following command, then do build and test. This is because Goloader relies on the internal package, which is forbidden by the Go compiler.
```
cp -r $GOROOT/src/cmd/internal $GOROOT/src/cmd/objfile
```

## Examples

```
go build github.com/pkujhd/goloader/examples/loader

go tool compile $GOPATH/src/github.com/pkujhd/goloader/examples/schedule/schedule.go
./loader -o schedule.o -run main.main -times 10

go tool compile $GOPATH/src/github.com/pkujhd/goloader/examples/base/base.go
./loader -o base.o -run main.main

go tool compile $GOPATH/src/github.com/pkujhd/goloader/examples/http/http.go
./loader -o http.o -run main.main

go install github.com/pkujhd/goloader/examples/basecontext
go tool compile -I $GOPATH/pkg/`go env GOOS`_`go env GOARCH`/ $GOPATH/src/github.com/pkujhd/goloader/examples/inter/inter.go
./loader -o $GOPATH/pkg/`go env GOOS`_`go env GOARCH`/github.com/pkujhd/goloader/examples/basecontext.a:github.com/pkujhd/goloader/examples/basecontext -o inter.o

#build multiple go files
go tool compile -I $GOPATH/pkg/`go env GOOS`_`go env GOARCH`/ -o test.o test1.go test2.go
./loader -o test.o -run main.main

```

## Warning

This has currently only been tested and developed on:

Golang 1.8-1.15 (x64/x86, darwin, linux, windows)

Golang 1.10-1.15 (arm, linux)

Golang 1.8-1.15 (arm64(LE) linux)
