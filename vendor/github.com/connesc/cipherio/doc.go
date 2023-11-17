// Package cipherio allows to use block ciphers with io.Reader and io.Writer.
//
// Golang already provides io.Reader and io.Writer implementations for cipher.Stream, but not for
// cipher.BlockMode (such as AES-CBC). The purpose of this package is to fill the gap.
//
// Block ciphers require data size to be a multiple of the block size. The io.Reader and io.Writer
// implementations found here can either enforce this requirement or automatically apply a
// user-defined padding.
//
// This package has been written with performance in mind: buffering and copies are avoided as much
// as possible.
package cipherio
