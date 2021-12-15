// Copyright 2020 cetc-30. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
// license that can be found in the LICENSE file.

package sm2

import (
	"bytes"
	"crypto"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"math/big"

	"github.com/admpub/ccs-gm/sm3"
)

var EncryptionErr = errors.New("sm2: encryption error")
var DecryptionErr = errors.New("sm2: decryption error")

func (key *PrivateKey) Decrypt(rand io.Reader, msg []byte, opts crypto.DecrypterOpts) (plaintext []byte, err error) {
	return Decrypt(msg, key)
}

func keyDerivation(Z []byte, klen int) []byte {
	var ct = 1
	if klen%8 != 0 {
		return nil
	}

	K := make([]byte, int(math.Ceil(float64(klen)/(sm3.Size*8))*sm3.Size))
	v := sm3.Size * 8

	l := int(math.Ceil(float64(klen) / float64(v)))

	var m = make([]byte, len(Z)+4)
	var vBytes = make([]byte, 4)
	copy(m, Z)

	for ; ct <= l; ct++ {
		binary.BigEndian.PutUint32(vBytes, uint32(ct))
		copy(m[len(Z):], vBytes)

		hash := sm3.SumSM3(m)
		copy(K[(ct-1)*sm3.Size:], hash[:])
	}
	return K[:klen/8]
}

func Encrypt(rand io.Reader, key *PublicKey, msg []byte) (cipher []byte, err error) {
	x, y, c2, c3, err := doEncrypt(rand, key, msg)
	if err != nil {
		return nil, err
	}

	c1 := pointToBytes(x, y)

	//c = c1||c2||c3,len(c1)=65,len(c3)=32
	cipher = append(c1, c2...)
	cipher = append(cipher, c3...)

	return
}

func doEncrypt(rand io.Reader, key *PublicKey, msg []byte) (x, y *big.Int, c2, c3 []byte, err error) {
	k := generateRandK(rand, key.Curve)

regen:
	x1, y1 := key.Curve.ScalarBaseMult(k.Bytes())

	var x2, y2 *big.Int
	if opt, ok := key.Curve.(optMethod); ok && (key.PreComputed != nil) {
		x2, y2 = opt.PreScalarMult(key.PreComputed, k.Bytes())
	} else {
		x2, y2 = key.Curve.ScalarMult(key.X, key.Y, k.Bytes())
	}

	xBuf := x2.Bytes()
	yBuf := y2.Bytes()

	xPadding := make([]byte, 32)
	yPadding := make([]byte, 32)
	if n := len(xBuf); n < 32 {
		xBuf = append(xPadding[:32-n], xBuf...)
	}

	if n := len(yBuf); n < 32 {
		yBuf = append(yPadding[:32-n], yBuf...)
	}

	//z=x2||y2
	Z := make([]byte, 64)
	copy(Z, xBuf)
	copy(Z[32:], yBuf)

	t := keyDerivation(Z, len(msg)*8)
	if t == nil {
		return nil, nil, nil, nil, EncryptionErr
	}
	for i, v := range t {
		if v != 0 {
			break
		}
		if i == len(t)-1 {
			goto regen
		}
	}

	//M^t
	for i, v := range t {
		t[i] = v ^ msg[i]
	}

	m3 := make([]byte, 64+len(msg))
	copy(m3, xBuf)
	copy(m3[32:], msg)
	copy(m3[32+len(msg):], yBuf)
	h := sm3.SumSM3(m3)
	c3 = h[:]

	return x1, y1, t, c3, nil
}

func Decrypt(c []byte, key *PrivateKey) ([]byte, error) {
	x1, y1 := pointFromBytes(c[:65])

	//dB*C1
	x2, y2 := key.Curve.ScalarMult(x1, y1, key.D.Bytes())

	xBuf := x2.Bytes()
	yBuf := y2.Bytes()

	xPadding := make([]byte, 32)
	yPadding := make([]byte, 32)
	if n := len(xBuf); n < 32 {
		xBuf = append(xPadding[:32-n], xBuf...)
	}

	if n := len(yBuf); n < 32 {
		yBuf = append(yPadding[:32-n], yBuf...)
	}

	//z=x2||y2
	Z := make([]byte, 64)
	copy(Z, xBuf)
	copy(Z[32:], yBuf)

	t := keyDerivation(Z, (len(c)-97)*8)
	if t == nil {
		return nil, DecryptionErr
	}
	for i, v := range t {
		if v != 0 {
			break
		}
		if i == len(t)-1 {
			return nil, DecryptionErr
		}
	}

	// m` = c2 ^ t
	c2 := c[65:(len(c) - 32)]
	for i, v := range t {
		t[i] = v ^ c2[i]
	}

	//validate
	_u := make([]byte, 64+len(t))
	copy(_u, xBuf)
	copy(_u[32:], t)
	copy(_u[32+len(t):], yBuf)
	u := sm3.SumSM3(_u)
	if !bytes.Equal(u[:], c[65+len(c2):]) {
		return nil, DecryptionErr
	}

	return t, nil
}

//uncompressed form, s=04||x||y
func pointToBytes(x, y *big.Int) []byte {
	buf := []byte{}

	xBuf := x.Bytes()
	yBuf := y.Bytes()

	xPadding := make([]byte, 32)
	yPadding := make([]byte, 32)
	if n := len(xBuf); n < 32 {
		xBuf = append(xPadding[:32-n], xBuf...)
	}

	if n := len(yBuf); n < 32 {
		yBuf = append(yPadding[:32-n], yBuf...)
	}

	//s = 04||x||y
	buf = append(buf, 0x4)
	buf = append(buf, xBuf...)
	buf = append(buf, yBuf...)

	return buf
}

func pointFromBytes(buf []byte) (x, y *big.Int) {
	if len(buf) != 65 || buf[0] != 0x4 {
		return nil, nil
	}

	x = new(big.Int).SetBytes(buf[1:33])
	y = new(big.Int).SetBytes(buf[33:])

	return
}
