[![GitHub release](https://img.shields.io/github/v/release/bodgit/sevenzip)](https://github.com/bodgit/sevenzip/releases)
[![Build Status](https://img.shields.io/github/workflow/status/bodgit/sevenzip/build)](https://github.com/bodgit/sevenzip/actions?query=workflow%3Abuild)
[![Coverage Status](https://coveralls.io/repos/github/bodgit/sevenzip/badge.svg?branch=master)](https://coveralls.io/github/bodgit/sevenzip?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bodgit/sevenzip)](https://goreportcard.com/report/github.com/bodgit/sevenzip)
[![GoDoc](https://godoc.org/github.com/bodgit/sevenzip?status.svg)](https://godoc.org/github.com/bodgit/sevenzip)
![Go version](https://img.shields.io/badge/Go-1.18-brightgreen.svg)
![Go version](https://img.shields.io/badge/Go-1.17-brightgreen.svg)

sevenzip
========

A reader for 7-zip archives inspired by `archive/zip`.

Current status:

* Pure Go, no external libraries or binaries needed.
* Handles uncompressed headers, (`7za a -mhc=off test.7z ...`).
* Handles compressed headers, (`7za a -mhc=on test.7z ...`).
* Handles password-protected versions of both of the above (`7za a -mhc=on|off -mhe=on -ppassword test.7z ...`).
* Handles archives split into multiple volumes, (`7za a -v100m test.7z ...`).
* Validates CRC values as it parses the file.
* Supports BCJ2, Brotli, Bzip2, Copy, Deflate, Delta, LZ4, LZMA, LZMA2 and Zstandard methods.
* Implements the `fs.FS` interface so you can treat an opened 7-zip archive like a filesystem.

More examples of 7-zip archives are needed to test all of the different combinations/algorithms possible.
