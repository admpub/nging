# cipherio

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c)](https://pkg.go.dev/github.com/connesc/cipherio)
[![Go Report Card](https://goreportcard.com/badge/github.com/connesc/cipherio)](https://goreportcard.com/report/github.com/connesc/cipherio)
[![GitHub tag](https://img.shields.io/github/v/tag/connesc/cipherio?sort=semver)](https://github.com/connesc/cipherio/tags)
[![License](https://img.shields.io/github/license/connesc/cipherio)](LICENSE)

This Golang package allows to use block ciphers with `io.Reader` and `io.Writer`.

Golang already provides [`io.Reader`](https://golang.org/pkg/io/#Reader) and [`io.Writer`](https://golang.org/pkg/io/#Writer) implementations for [`cipher.Stream`](https://golang.org/pkg/crypto/cipher/#Stream), but not for [`cipher.BlockMode`](https://golang.org/pkg/crypto/cipher/#BlockMode) (such as AES-CBC). The purpose of this package is to fill the gap.

Block ciphers require data size to be a multiple of the block size. The `io.Reader` and `io.Writer` implementations found here can either enforce this requirement or automatically apply a user-defined padding.

This package has been written with performance in mind: buffering and copies are avoided as much as possible.
