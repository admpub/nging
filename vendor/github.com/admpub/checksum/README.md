checksum
==
[![PkgGoDev](https://pkg.go.dev/badge/github.com/codingsince1985/checksum)](https://pkg.go.dev/github.com/codingsince1985/checksum)
[![Go Report Card](https://goreportcard.com/badge/codingsince1985/checksum)](https://goreportcard.com/report/codingsince1985/checksum)
[![codecov](https://codecov.io/gh/codingsince1985/checksum/branch/master/graph/badge.svg?token=BO1wAQwc85)](https://codecov.io/gh/codingsince1985/checksum)

Compute message digest, like MD5, SHA256, SHA1, CRC or BLAKE2s-256 in Golang for potentially large files.

Usage
--
```go
package main

import (
	"fmt"
	"github.com/admpub/checksum"
)

func main() {
	file := "/home/jerry/Downloads/ubuntu-20.04.2.0-desktop-amd64.iso"

	md5, _ := checksum.MD5sum(file)
	fmt.Println(md5)

	sha256, _ := checksum.SHA256sum(file)
	fmt.Println(sha256)

	sha1, _ := checksum.SHA1sum(file)
	fmt.Println(sha1)

	crc32, _ := checksum.CRC32(file)
	fmt.Println(crc32)

	blake2s256, _ := checksum.Blake2s256(file)
	fmt.Println(blake2s256)
}
```
License
==
checksum is distributed under the terms of the MIT license. See LICENSE for details.
