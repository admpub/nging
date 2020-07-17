# goseaweedfs

[![Build Status](https://travis-ci.org/linxGnu/goseaweedfs.svg?branch=master)](https://travis-ci.org/linxGnu/goseaweedfs)
[![Go Report Card](https://goreportcard.com/badge/github.com/admpub/goseaweedfs)](https://goreportcard.com/report/github.com/admpub/goseaweedfs)
[![Coverage Status](https://coveralls.io/repos/github/linxGnu/goseaweedfs/badge.svg?branch=master)](https://coveralls.io/github/linxGnu/goseaweedfs?branch=master)
[![godoc](https://img.shields.io/badge/docs-GoDoc-green.svg)](https://godoc.org/github.com/admpub/goseaweedfs)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/admpub/goseaweedfs/blob/master/LICENSE)

A complete Golang client for [SeaweedFS](https://github.com/chrislusf/seaweedfs) (version 1.44+). Inspired by:
- [tnextday/goseaweed](https://github.com/tnextday/goseaweed)
- [ginuerzh/weedo](https://github.com/ginuerzh/weedo)

## Installation
```
go get -u github.com/admpub/goseaweedfs
```

## Usage
Please refer to [Test Cases](https://github.com/admpub/goseaweedfs/blob/master/seaweed_test.go) for sample code.

## Supported

- [x] Grow
- [x] Status
- [x] Cluster Status
- [x] Filer
- [x] Upload
- [x] Submit
- [x] Delete
- [x] Replace
- [x] Upload large file with builtin manifest handler, auto file split and chunking
- [ ] Admin Operations (mount, unmount, delete volumn, etc)

## Contributing
Please issue me for things gone wrong or:

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D