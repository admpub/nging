Package go-download
===================

Package go-download provides a library for interruptable, resumable download acceleration with automatic Accept-Ranges support

It Features:
- [x] Customizable concurrency and/or chunk size. default is 10 goroutines
- [x] Proxy of download eg. to display a progress bar

## Installation
```shell
go get -u github.com/admpub/go-download/v2
```
or if your looking for the standalone client
```shell
go get -u github.com/admpub/go-download/v2/cmd/goget
```

## Examples

More examples [here](https://github.com/admpub/go-download/tree/master/_examples)

```go
package main

import (
	"log"

	download "github.com/admpub/go-download/v2"
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
