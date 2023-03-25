// Copyright 2020 cetc-30. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
// license that can be found in the LICENSE file.

// Package sm2 implements china crypto standards.
package sm2

import (
	"crypto"
	"crypto/elliptic"
	"io"
	"math/big"

	"github.com/admpub/ccs-gm/sm3"
)

type PublicKey struct {
	elliptic.Curve
	X, Y        *big.Int
	PreComputed *[37][64 * 8]uint64 //precomputation
}

type PrivateKey struct {
	PublicKey
	D    *big.Int
	DInv *big.Int //(1+d)^-1
}

var generateRandK = _generateRandK

//optMethod includes some optimized methods.
type optMethod interface {
	// CombinedMult implements fast multiplication S1*g + S2*p (g - generator, p - arbitrary point)
	CombinedMult(Precomputed *[37][64 * 8]uint64, baseScalar, scalar []byte) (x, y *big.Int)
	// InitPubKeyTable implements precomputed table of public key
	InitPubKeyTable(x, y *big.Int) (Precomputed *[37][64 * 8]uint64)
	// PreScalarMult implements fast multiplication of public key
	PreScalarMult(Precomputed *[37][64 * 8]uint64, scalar []byte) (x, y *big.Int)
}

// The SM2's private key contains the public key
func (priv *PrivateKey) Public() crypto.PublicKey {
	return &priv.PublicKey
}

var one = new(big.Int).SetInt64(1)

func randFieldElement(c elliptic.Curve, rand io.Reader) (k *big.Int, err error) {
	params := c.Params()
	b := make([]byte, params.BitSize/8+8)
	_, err = io.ReadFull(rand, b)
	if err != nil {
		return
	}
	k = new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)
	return
}

func GenerateKey(rand io.Reader) (*PrivateKey, error) {
	c := P256()

	k, err := randFieldElement(c, rand)
	if err != nil {
		return nil, err
	}
	priv := new(PrivateKey)
	priv.PublicKey.Curve = c
	priv.D = k
	//(1+d)^-1
	priv.DInv = new(big.Int).Add(k, one)
	priv.DInv.ModInverse(priv.DInv, c.Params().N)
	priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(k.Bytes())
	if opt, ok := c.(optMethod); ok {
		priv.PreComputed = opt.InitPubKeyTable(priv.PublicKey.X, priv.PublicKey.Y)
	}
	return priv, nil
}

func _generateRandK(rand io.Reader, c elliptic.Curve) (k *big.Int) {
	params := c.Params()
	b := make([]byte, params.BitSize/8+8)
	_, err := io.ReadFull(rand, b)
	if err != nil {
		return
	}
	k = new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)
	return
}

func getZById(pub *PublicKey, id []byte) []byte {
	c := P256()
	var lena = uint16(len(id) * 8) //bit len of IDA
	var ENTLa = []byte{byte(lena >> 8), byte(lena)}
	var z = make([]byte, 0, 1024)

	//判断公钥x,y坐标长度是否小于32字节，若小于则在前面补0
	xBuf := pub.X.Bytes()
	yBuf := pub.Y.Bytes()

	xPadding := make([]byte, 32)
	yPadding := make([]byte, 32)

	if n := len(xBuf); n < 32 {
		xBuf = append(xPadding[:32-n], xBuf...)
	}

	if n := len(yBuf); n < 32 {
		yBuf = append(yPadding[:32-n], yBuf...)
	}

	z = append(z, ENTLa...)
	z = append(z, id...)
	z = append(z, SM2PARAM_A...)
	z = append(z, c.Params().B.Bytes()...)
	z = append(z, c.Params().Gx.Bytes()...)
	z = append(z, c.Params().Gy.Bytes()...)
	z = append(z, xBuf...)
	z = append(z, yBuf...)

	//h := sm3.New()
	hash := sm3.SumSM3(z)
	return hash
}

//Za = sm3(ENTL||IDa||a||b||Gx||Gy||Xa||Xy)
func getZ(pub *PublicKey) []byte {
	return getZById(pub, []byte("1234567812345678"))
}

func Sign(rand io.Reader, priv *PrivateKey, msg []byte) (r, s *big.Int, err error) {
	//if len(hash) < 32 {
	//	err = errors.New("The length of hash has short than what SM2 need.")
	//	return
	//}

	var m = make([]byte, 32+len(msg))
	copy(m, getZ(&priv.PublicKey))
	copy(m[32:], msg)

	hash := sm3.SumSM3(m)
	e := new(big.Int).SetBytes(hash)
	k := generateRandK(rand, priv.PublicKey.Curve)

	x1, _ := priv.PublicKey.Curve.ScalarBaseMult(k.Bytes())

	n := priv.PublicKey.Curve.Params().N

	r = e.Add(e, x1)

	r.Mod(r, n)

	s1 := new(big.Int).Mul(r, priv.D)
	s1.Sub(k, s1)

	s2 := new(big.Int)
	if priv.DInv == nil {
		s2 = s2.Add(one, priv.D)
		s2.ModInverse(s2, n)
	} else {
		s2 = priv.DInv
	}

	s = s1.Mul(s1, s2)
	s.Mod(s, n)

	return
}

func SignWithDigest(rand io.Reader, priv *PrivateKey, digest []byte) (r, s *big.Int, err error) {
	//if len(hash) < 32 {
	//	err = errors.New("The length of hash has short than what SM2 need.")
	//	return
	//}
	e := new(big.Int).SetBytes(digest)
	k := generateRandK(rand, priv.PublicKey.Curve)

	x1, _ := priv.PublicKey.Curve.ScalarBaseMult(k.Bytes())

	n := priv.PublicKey.Curve.Params().N

	r = e.Add(e, x1)

	r.Mod(r, n)

	s1 := new(big.Int).Mul(r, priv.D)
	s1.Mod(s1, n)
	s1.Sub(k, s1)
	s1.Mod(s1, n)

	s2 := new(big.Int)
	if priv.DInv == nil {
		s2 = s2.Add(one, priv.D)
		s2.ModInverse(s2, n)
	} else {
		s2 = priv.DInv
	}

	s = s1.Mul(s1, s2)
	s.Mod(s, n)

	return
}

func Verify(pub *PublicKey, msg []byte, r, s *big.Int) bool {
	c := pub.Curve
	N := c.Params().N

	if r.Sign() <= 0 || s.Sign() <= 0 {
		return false
	}
	if r.Cmp(N) >= 0 || s.Cmp(N) >= 0 {
		return false
	}

	n := c.Params().N

	var m = make([]byte, 32+len(msg))
	copy(m, getZ(pub))
	copy(m[32:], msg)
	//h := sm3.New()
	//hash := h.Sum(m)
	hash := sm3.SumSM3(m)
	e := new(big.Int).SetBytes(hash[:])

	t := new(big.Int).Add(r, s)

	// Check if implements S1*g + S2*p
	//Using fast multiplication CombinedMult.
	var x1 *big.Int
	if opt, ok := c.(optMethod); ok && (pub.PreComputed != nil) {
		x1, _ = opt.CombinedMult(pub.PreComputed, s.Bytes(), t.Bytes())
	} else {
		x11, y11 := c.ScalarMult(pub.X, pub.Y, t.Bytes())
		x12, y12 := c.ScalarBaseMult(s.Bytes())
		x1, _ = c.Add(x11, y11, x12, y12)
	}

	e.Add(e, x1)
	e.Mod(e, n)

	return e.Cmp(r) == 0
}

func VerifyWithDigest(pub *PublicKey, digest []byte, r, s *big.Int) bool {
	c := pub.Curve
	N := c.Params().N

	if r.Sign() <= 0 || s.Sign() <= 0 {
		return false
	}
	if r.Cmp(N) >= 0 || s.Cmp(N) >= 0 {
		return false
	}

	n := pub.Curve.Params().N

	e := new(big.Int).SetBytes(digest)

	t := new(big.Int).Add(r, s)
	// Check if implements S1*g + S2*p
	//Using fast multiplication CombinedMult.
	var x1 *big.Int
	if opt, ok := c.(optMethod); ok && (pub.PreComputed != nil) {
		x1, _ = opt.CombinedMult(pub.PreComputed, s.Bytes(), t.Bytes())
	} else {
		x11, y11 := c.ScalarMult(pub.X, pub.Y, t.Bytes())
		x12, y12 := c.ScalarBaseMult(s.Bytes())
		x1, _ = c.Add(x11, y11, x12, y12)
	}
	e.Add(e, x1)
	e.Mod(e, n)

	return e.Cmp(r) == 0
}

type zr struct {
	io.Reader
}

func (z *zr) Read(dst []byte) (n int, err error) {
	for i := range dst {
		dst[i] = 0
	}
	return len(dst), nil
}

var zeroReader = &zr{}

//func OptSign(rand io.Reader, priv *PrivateKey, msg []byte) (r, s *big.Int, err error) {
//	//var one = new(big.Int).SetInt64(1)
//	//if len(hash) < 32 {
//	//	err = errors.New("The length of hash has short than what SM2 need.")
//	//	return
//	//}
//
//	var m = make([]byte, 32+len(msg))
//	copy(m, getZ(&priv.PublicKey))
//	copy(m[32:], msg)
//
//	//h := sm3.New()
//	//hash := h.Sum(m)
//	hash := sm3.SumSM3(m)
//	e := new(big.Int).SetBytes(hash[:])
//	k := generateRandK(rand, priv.PublicKey.Curve)
//
//	x1, _ := priv.PublicKey.Curve.ScalarBaseMult(k.Bytes())
//
//	n := priv.PublicKey.Curve.Params().N
//
//	r = new(big.Int).Add(e, x1)
//
//	r.Mod(r, n)
//
//	s1 := new(big.Int).Mul(r, priv.D)
//	//s1.Mod(s1, n)
//	s1.Sub(k, s1)
//	s1.Mod(s1, n)
//
//	//s2 := new(big.Int).Add(one, priv.D)
//	//s2.Mod(s2, n)
//	//s2.ModInverse(s2, n)
//	s = new(big.Int).Mul(s1, priv.DInv)
//	s.Mod(s, n)
//
//	return
//}
//
//func OptVerify(pub *PublicKey, msg []byte, r, s *big.Int) bool {
//	c := pub.Curve
//	N := c.Params().N
//
//	if r.Sign() <= 0 || s.Sign() <= 0 {
//		return false
//	}
//	if r.Cmp(N) >= 0 || s.Cmp(N) >= 0 {
//		return false
//	}
//
//	n := c.Params().N
//
//	var m = make([]byte, 32+len(msg))
//	copy(m, getZ(pub))
//	copy(m[32:], msg)
//	//h := sm3.New()
//	//hash := h.Sum(m)
//	hash := sm3.SumSM3(m)
//	e := new(big.Int).SetBytes(hash[:])
//
//	t := new(big.Int).Add(r, s)
//
//	// Check if implements S1*g + S2*p
//	//Using fast multiplication CombinedMult.
//	var x1 *big.Int
//	if opt, ok := c.(optMethod); ok {
//		//x11, y11 := opt.PreScalarMult(pub.PreComputed,t.Bytes())
//		//x12, y12 := c.ScalarBaseMult(s.Bytes())
//		//x1, _ = c.Add(x11, y11, x12, y12)
//		x1, _ = opt.CombinedMult(pub.PreComputed, s.Bytes(), t.Bytes())
//	} else {
//		x11, y11 := c.ScalarMult(pub.X, pub.Y, t.Bytes())
//		x12, y12 := c.ScalarBaseMult(s.Bytes())
//		x1, _ = c.Add(x11, y11, x12, y12)
//	}
//
//	x := new(big.Int).Add(e, x1)
//	x = x.Mod(x, n)
//
//	return x.Cmp(r) == 0
//}
