// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sm3 implements china crypto standards.
package sm3

import (
	"hash"
)

var hashFunc func() hash.Hash

func init() {
	//crypto.RegisterHash(crypto.tjSM3, New)
	hashFunc = New
}

// The size of a SM2 checksum in bytes.
const Size = 32

// The blocksize of SHA256 and SHA224 in bytes.
const BlockSize = 64

const (
	chunk = 64
	init0 = 0x7380166f
	init1 = 0x4914b2b9
	init2 = 0x172442d7
	init3 = 0xda8a0600
	init4 = 0xa96f30bc
	init5 = 0x163138aa
	init6 = 0xe38dee4d
	init7 = 0xb0fb0e4e
)

// digest represents the partial evaluation of a checksum.
type digest struct {
	h   [8]uint32
	x   [chunk]byte
	nx  int
	len uint64
}

func (d *digest) Reset() {
	d.h[0] = init0
	d.h[1] = init1
	d.h[2] = init2
	d.h[3] = init3
	d.h[4] = init4
	d.h[5] = init5
	d.h[6] = init6
	d.h[7] = init7
	d.nx = 0
	d.len = 0
}

func GetFunc() (func() hash.Hash){
	return hashFunc
}

func New() hash.Hash {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Size() int {
	return Size
}

func (d *digest) BlockSize() int { return BlockSize }

func (d *digest) Write(p []byte) (nn int, err error) {
	nn = len(p)
	d.len += uint64(nn)
	//var n int
	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		d.nx += n
		if d.nx == chunk {
			Block(d, d.x[:])
			d.nx = 0
		}
		p = p[n:]
	}

	if len(p) >= chunk {
		n := len(p) &^ (chunk - 1)
		Block(d, p)
		p = p[n:]
	}
	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}
	return
}

func (d0 *digest) Sum(in []byte) []byte {
	// Make a copy of d0 so that caller can keep writing and summing.
	d := *d0
	hash := d.checkSum()
	return append(in, hash[:]...)
}

func (d0 *digest)ConstantTimeSum(b []byte) []byte {
	return d0.Sum(b)
}

func (d *digest) checkSum() []byte {
	len := d.len
	// Padding. Add a 1 bit and 0 bits until 56 bytes mod 64.
	var tmp [64]byte
	tmp[0] = 0x80
	if len%64 < 56 {
		d.Write(tmp[0 : 56-len%64])
	} else {
		d.Write(tmp[0 : 64+56-len%64])
	}

	// Length in bits.
	len <<= 3
	for i := uint(0); i < 8; i++ {
		tmp[i] = byte(len >> (56 - 8*i))
	}
	d.Write(tmp[0:8])

	if d.nx != 0 {
		panic("d.nx != 0")
	}

	h := d.h[:]

	var digest []byte = make([]byte, Size)
	for i, s := range h {
		digest[i*4] = byte(s >> 24)
		digest[i*4+1] = byte(s >> 16)
		digest[i*4+2] = byte(s >> 8)
		digest[i*4+3] = byte(s)
	}
	return digest
}

func SumSM3(data []byte) []byte {
	var d digest
	d.Reset()
	d.Write(data)
	return d.checkSum()
}