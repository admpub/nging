package cipherio

import "fmt"

// Padding defines how to fill an incomplete block.
type Padding interface {
	Fill(dst []byte)
}

// PaddingFunc allows to implement the Padding interface with a padding function.
type PaddingFunc func(dst []byte)

// Fill an incomplete block.
func (p PaddingFunc) Fill(dst []byte) {
	p(dst)
}

// ZeroPadding fills an incomplete block with zeroes.
var ZeroPadding = PaddingFunc(zeroPadding)

// BitPadding fills an incomplete block with 0x80 followed by zeroes.
//
// This is defined by ISO/IEC 9797-1 as Padding Method 2 and is also known as ISO padding.
var BitPadding = PaddingFunc(bitPadding)

// PKCS7Padding fills an incomplete block by repeating the total number of padding bytes.
//
// PKCS#7 is described by RFC 5652.
//
// WARNING: this padding method MUST NOT be used with a block size larger than 256 bytes.
var PKCS7Padding = PaddingFunc(pkcs7Padding)

func fill(dst []byte, val byte) {
	for i := range dst {
		dst[i] = val
	}
}

func zeroPadding(dst []byte) {
	fill(dst, 0)
}

func bitPadding(dst []byte) {
	dst[0] = 0x80
	fill(dst[1:], 0)
}

func pkcs7Padding(dst []byte) {
	n := len(dst)
	if n > 255 {
		panic(fmt.Errorf("cipherio: PKCS#7 padding cannot fill more than 255 bytes: %d > 255", n))
	}
	fill(dst, byte(n))
}
