Package go-download
===================
![Project status](https://img.shields.io/badge/version-2.1.0-green.svg)
[![Build Status](https://travis-ci.org/joeybloggs/go-download.svg?branch=master)](https://travis-ci.org/joeybloggs/go-download)
[![Coverage Status](https://coveralls.io/repos/github/joeybloggs/go-download/badge.svg?branch=master)](https://coveralls.io/github/joeybloggs/go-download?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/joeybloggs/go-download)](https://goreportcard.com/report/github.com/joeybloggs/go-download)
[![GoDoc](https://godoc.org/github.com/joeybloggs/go-download?status.svg)](https://godoc.org/github.com/joeybloggs/go-download)
![License](https://img.shields.io/badge/license-BSD%202--clause-blue.svg)

Package go-download provides a library for interruptable, resumable download acceleration with automatic Accept-Ranges support

It Features:
- [x] Customizable concurrency and/or chunk size. default is 10 goroutines
- [x] Proxy of download eg. to display a progress bar

## Installation
```shell
go get -u github.com/joeybloggs/go-download
```
or if your looking for the standalone client
```shell
go get -u github.com/joeybloggs/go-download/cmd/goget
```

## Examples

More examples [here](https://github.com/joeybloggs/go-download/tree/master/_examples)

```go
package main

import (
	"log"

	download "github.com/joeybloggs/go-download"
)

func main() {

	// no options specified so will default to 10 concurrent download by default

	f, err := download.Open("https://storage.googleapis.com/golang/go1.8.1.src.tar.gz", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// f implements io.Reader, write file somewhere or do some other sort of work with it
}
```

## Contributing

Pull requests, bug fixes and issue reports are welcome.

Before proposing a change, please discuss your change by raising an issue.

## License

Distributed under BSD 2-clause license, please see license file in code for more details.
