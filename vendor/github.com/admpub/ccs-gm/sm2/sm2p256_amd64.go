// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the Go wrapper for the constant-time, 64-bit assembly
// implementation of sm2 curve. The optimizations performed here are described in
// detail in:
// S.Gueron and V.Krasnov, "Fast prime field elliptic-curve cryptography with
//                          256-bit primes"
// http://link.springer.com/article/10.1007%2Fs13389-014-0090-x
// https://eprint.iacr.org/2013/816.pdf

// +build amd64

package sm2

import (
	"crypto/elliptic"
	"fmt"
	"math/big"
	"sync"
)

type (
	p256Curve struct {
		*elliptic.CurveParams
	}

	p256Point struct {
		xyz [12]uint64
	}
)

var (
	p256            p256Curve
	p256Precomputed *[37][64 * 8]uint64
	precomputeOnce  sync.Once
)

func initP256() {
	// See FIPS 186-3, section D.2.3
	p256.CurveParams = &elliptic.CurveParams{Name: "SM2-P-256"}
	p256.P, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	p256.N, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	p256.B, _ = new(big.Int).SetString("28E9FA9E9D9F5E344D5A9E4BCF6509A7F39789F515AB8F92DDBCBD414D940E93", 16)
	p256.Gx, _ = new(big.Int).SetString("32C4AE2C1F1981195F9904466A39C9948FE30BBFF2660BE1715A4589334C74C7", 16)
	p256.Gy, _ = new(big.Int).SetString("BC3736A2F4F6779C59BDCEE36B692153D0A9877CC62A474002DF32E52139F0A0", 16)
	p256.BitSize = 256
}

func (curve p256Curve) Params() *elliptic.CurveParams {
	return curve.CurveParams
}

// Functions implemented in sm2p256_amd64.s
// Montgomery multiplication modulo P256
func sm2p256Mul(res, in1, in2 []uint64)
func p256TestMul(res, in1, in2 []uint64)

// Montgomery square modulo P256
func sm2p256Sqr(res, in []uint64)

// Montgomery multiplication by 1
func sm2p256FromMont(res, in []uint64)

// iff cond == 1  val <- -val
func sm2p256NegCond(val []uint64, cond int)

// if cond == 0 res <- b; else res <- a
func sm2p256MovCond(res, a, b []uint64, cond int)

// Endianness swap
func sm2p256BigToLittle(res []uint64, in []byte)
func sm2p256LittleToBig(res []byte, in []uint64)

// Constant time table access
func sm2p256Select(point, table []uint64, idx int)
func sm2p256SelectBase(point, table []uint64, idx int)

// Montgomery multiplication modulo Ord(G)
func sm2p256OrdMul(res, in1, in2 []uint64)

// Montgomery square modulo Ord(G), repeated n times
func sm2p256OrdSqr(res, in []uint64, n int)

// Point add with in2 being affine point
// If sign == 1 -> in2 = -in2
// If sel == 0 -> res = in1
// if zero == 0 -> res = in2
func sm2p256PointAddAffineAsm(res, in1, in2 []uint64, sign, sel, zero int)

// Point add
func sm2p256PointAddAsm(res, in1, in2 []uint64) int

// Point double
func sm2p256PointDoubleAsm(res, in []uint64)

//Test Internal Func
func sm2p256TestSubInternal(res, in1, in2 []uint64)
func sm2p256TestMulInternal(res, in1, in2 []uint64)
func sm2p256TestMulBy2Inline(res, in1 []uint64)
func sm2p256TestSqrInternal(res, in1 []uint64)
func sm2p256TestAddInline(res, in1, in2 []uint64)

func (curve p256Curve) Inverse(k *big.Int) *big.Int {
	if k.Sign() < 0 {
		// This should never happen.
		k = new(big.Int).Neg(k)
	}

	if k.Cmp(p256.N) >= 0 {
		// This should never happen.
		k = new(big.Int).Mod(k, p256.N)
	}

	// table will store precomputed powers of x. The four words at index
	// 4×i store x^(i+1).
	var table [4 * 15]uint64

	x := make([]uint64, 4)
	fromBig(x[:], k)
	// This code operates in the Montgomery domain where R = 2^256 mod n
	// and n is the order of the scalar field. (See initP256 for the
	// value.) Elements in the Montgomery domain take the form a×R and
	// multiplication of x and y in the calculates (x × y × R^-1) mod n. RR
	// is R×R mod n thus the Montgomery multiplication x and RR gives x×R,
	// i.e. converts x into the Montgomery domain.
	//	RR := []uint64{0x83244c95be79eea2, 0x4699799c49bd6fa6, 0x2845b2392b6bec59, 0x66e12d94f3d95620}
	RR := []uint64{0x901192AF7C114F20, 0x3464504ADE6FA2FA, 0x620FC84C3AFFE0D4, 0x1EB5E412A22B3D3B}

	sm2p256OrdMul(table[:4], x, RR)

	// Prepare the table, no need in constant time access, because the
	// power is not a secret. (Entry 0 is never used.)
	for i := 2; i < 16; i += 2 {
		sm2p256OrdSqr(table[4*(i-1):], table[4*((i/2)-1):], 1)
		sm2p256OrdMul(table[4*i:], table[4*(i-1):], table[:4])
	}

	x[0] = table[4*14+0] // f
	x[1] = table[4*14+1]
	x[2] = table[4*14+2]
	x[3] = table[4*14+3]

	sm2p256OrdSqr(x, x, 4)
	sm2p256OrdMul(x, x, table[4*14:4*14+4]) // ff
	t := make([]uint64, 4, 4)
	t[0] = x[0]
	t[1] = x[1]
	t[2] = x[2]
	t[3] = x[3]

	sm2p256OrdSqr(x, x, 8)
	sm2p256OrdMul(x, x, t) // ffff
	t[0] = x[0]
	t[1] = x[1]
	t[2] = x[2]
	t[3] = x[3]

	sm2p256OrdSqr(x, x, 16)
	sm2p256OrdMul(x, x, t) // ffffffff
	t[0] = x[0]
	t[1] = x[1]
	t[2] = x[2]
	t[3] = x[3]

	sm2p256OrdSqr(x, x, 64) // ffffffff0000000000000000
	sm2p256OrdMul(x, x, t)  // ffffffff00000000ffffffff
	sm2p256OrdSqr(x, x, 32) // ffffffff00000000ffffffff00000000
	sm2p256OrdMul(x, x, t)  // ffffffff00000000ffffffffffffffff

	// Remaining 32 windows
	expLo := [32]byte{0xb, 0xc, 0xe, 0x6, 0xf, 0xa, 0xa, 0xd, 0xa, 0x7, 0x1, 0x7, 0x9, 0xe, 0x8, 0x4, 0xf, 0x3, 0xb, 0x9, 0xc, 0xa, 0xc, 0x2, 0xf, 0xc, 0x6, 0x3, 0x2, 0x5, 0x4, 0xf}
	for i := 0; i < 32; i++ {
		sm2p256OrdSqr(x, x, 4)
		sm2p256OrdMul(x, x, table[4*(expLo[i]-1):])
	}

	// Multiplying by one in the Montgomery domain converts a Montgomery
	// value out of the domain.
	one := []uint64{1, 0, 0, 0}
	sm2p256OrdMul(x, x, one)

	xOut := make([]byte, 32)
	sm2p256LittleToBig(xOut, x)
	return new(big.Int).SetBytes(xOut)
}

// fromBig converts a *big.Int into a format used by this code.
func fromBig(out []uint64, big *big.Int) {
	for i := range out {
		out[i] = 0
	}

	for i, v := range big.Bits() {
		out[i] = uint64(v)
	}
}

// p256GetScalar endian-swaps the big-endian scalar value from in and writes it
// to out. If the scalar is equal or greater than the order of the group, it's
// reduced modulo that order.
func p256GetScalar(out []uint64, in []byte) {
	n := new(big.Int).SetBytes(in)

	if n.Cmp(p256.N) >= 0 {
		n.Mod(n, p256.N)
	}
	fromBig(out, n)
}

// sm2p256Mul operates in a Montgomery domain with R = 2^256 mod p, where p is the
// underlying field of the curve. (See initP256 for the value.) Thus rr here is
// R×R mod p. See comment in Inverse about how this is used.
var rr = []uint64{0x0000000200000003, 0x00000002FFFFFFFF, 0x0000000100000001, 0x0000000400000002}

func maybeReduceModP(in *big.Int) *big.Int {
	if in.Cmp(p256.P) < 0 {
		return in
	}
	return new(big.Int).Mod(in, p256.P)
}

//CombinedMult implements fast multiplication baseScalar*G+scalar*P.
//CombinedMult returns baseScalar*G+scalar*P, where G is the base point of the group
//and P is the point (base or non-base point) of the group,
//baseScalar and scalar are integers in big-endian form.
//func (curve p256Curve) CombinedMult(bigX, bigY *big.Int, baseScalar, scalar []byte) (x, y *big.Int) {
//	scalarReversed := make([]uint64, 4)
//	var r1, r2 p256Point
//	p256GetScalar(scalarReversed, baseScalar)
//	r1IsInfinity := scalarIsZero(scalarReversed)
//	r1.p256BaseMult(scalarReversed)
//
//	p256GetScalar(scalarReversed, scalar)
//	r2IsInfinity := scalarIsZero(scalarReversed)
//	fromBig(r2.xyz[0:4], maybeReduceModP(bigX))
//	fromBig(r2.xyz[4:8], maybeReduceModP(bigY))
//	sm2p256Mul(r2.xyz[0:4], r2.xyz[0:4], rr[:])
//	sm2p256Mul(r2.xyz[4:8], r2.xyz[4:8], rr[:])
//
//	// This sets r2's Z value to 1, in the Montgomery domain.
//	//	r2.xyz[8] = 0x0000000000000001
//	//	r2.xyz[9] = 0xffffffff00000000
//	//	r2.xyz[10] = 0xffffffffffffffff
//	//	r2.xyz[11] = 0x00000000fffffffe
//	r2.xyz[8] = 0x0000000000000001
//	r2.xyz[9] = 0x00000000FFFFFFFF
//	r2.xyz[10] = 0x0000000000000000
//	r2.xyz[11] = 0x0000000100000000
//
//	//r2.p256ScalarMult(scalarReversed)
//	//sm2p256PointAddAsm(r1.xyz[:], r1.xyz[:], r2.xyz[:])
//
//	r2.p256ScalarMult(scalarReversed)
//
//	var sum, double p256Point
//	pointsEqual := sm2p256PointAddAsm(sum.xyz[:], r1.xyz[:], r2.xyz[:])
//	sm2p256PointDoubleAsm(double.xyz[:], r1.xyz[:])
//	sum.CopyConditional(&double, pointsEqual)
//	sum.CopyConditional(&r1, r2IsInfinity)
//	sum.CopyConditional(&r2, r1IsInfinity)
//	return sum.p256PointToAffine()
//}
func (curve p256Curve) CombinedMult(Precomputed *[37][64*8]uint64, baseScalar, scalar []byte) (x, y *big.Int) {
	scalarReversed := make([]uint64, 4)
	var r1 p256Point
	r2 := new(p256Point)
	p256GetScalar(scalarReversed, baseScalar)
	r1IsInfinity := scalarIsZero(scalarReversed)
	r1.p256BaseMult(scalarReversed)

	p256GetScalar(scalarReversed, scalar)
	r2IsInfinity := scalarIsZero(scalarReversed)
	//fromBig(r2.xyz[0:4], maybeReduceModP(bigX))
	//fromBig(r2.xyz[4:8], maybeReduceModP(bigY))
	//sm2p256Mul(r2.xyz[0:4], r2.xyz[0:4], rr[:])
	//sm2p256Mul(r2.xyz[4:8], r2.xyz[4:8], rr[:])
	//
	//// This sets r2's Z value to 1, in the Montgomery domain.
	////	r2.xyz[8] = 0x0000000000000001
	////	r2.xyz[9] = 0xffffffff00000000
	////	r2.xyz[10] = 0xffffffffffffffff
	////	r2.xyz[11] = 0x00000000fffffffe
	//r2.xyz[8] = 0x0000000000000001
	//r2.xyz[9] = 0x00000000FFFFFFFF
	//r2.xyz[10] = 0x0000000000000000
	//r2.xyz[11] = 0x0000000100000000
	//
	////r2.p256ScalarMult(scalarReversed)
	////sm2p256PointAddAsm(r1.xyz[:], r1.xyz[:], r2.xyz[:])

	//r2.p256ScalarMult(scalarReversed)
	r2.p256PreMult(Precomputed,scalarReversed)

	var sum, double p256Point
	pointsEqual := sm2p256PointAddAsm(sum.xyz[:], r1.xyz[:], r2.xyz[:])
	sm2p256PointDoubleAsm(double.xyz[:], r1.xyz[:])
	sum.CopyConditional(&double, pointsEqual)
	sum.CopyConditional(&r1, r2IsInfinity)
	sum.CopyConditional(r2, r1IsInfinity)
	return sum.p256PointToAffine()
}

func (curve p256Curve) ScalarBaseMult(scalar []byte) (x, y *big.Int) {
	scalarReversed := make([]uint64, 4)
	p256GetScalar(scalarReversed, scalar)

	var r p256Point
	r.p256BaseMult(scalarReversed)
	return r.p256PointToAffine()
}

func (curve p256Curve) ScalarMult(bigX, bigY *big.Int, scalar []byte) (x, y *big.Int) {
	scalarReversed := make([]uint64, 4)
	p256GetScalar(scalarReversed, scalar)

	var r p256Point
	fromBig(r.xyz[0:4], maybeReduceModP(bigX))
	fromBig(r.xyz[4:8], maybeReduceModP(bigY))
	sm2p256Mul(r.xyz[0:4], r.xyz[0:4], rr[:])
	sm2p256Mul(r.xyz[4:8], r.xyz[4:8], rr[:])
	// This sets r2's Z value to 1, in the Montgomery domain.
	//	r.xyz[8] = 0x0000000000000001
	//	r.xyz[9] = 0xffffffff00000000
	//	r.xyz[10] = 0xffffffffffffffff
	//	r.xyz[11] = 0x00000000fffffffe
	r.xyz[8] = 0x0000000000000001
	r.xyz[9] = 0x00000000FFFFFFFF
	r.xyz[10] = 0x0000000000000000
	r.xyz[11] = 0x0000000100000000

	r.p256ScalarMult(scalarReversed)
	return r.p256PointToAffine()
}

// uint64IsZero returns 1 if x is zero and zero otherwise.
func uint64IsZero(x uint64) int {
	x = ^x
	x &= x >> 32
	x &= x >> 16
	x &= x >> 8
	x &= x >> 4
	x &= x >> 2
	x &= x >> 1
	return int(x & 1)
}

// scalarIsZero returns 1 if scalar represents the zero value, and zero
// otherwise.
func scalarIsZero(scalar []uint64) int {
	return uint64IsZero(scalar[0] | scalar[1] | scalar[2] | scalar[3])
}

func (p *p256Point) p256PointToAffine() (x, y *big.Int) {
	zInv := make([]uint64, 4)
	zInvSq := make([]uint64, 4)
	p256Inverse(zInv, p.xyz[8:12])
	sm2p256Sqr(zInvSq, zInv)
	sm2p256Mul(zInv, zInv, zInvSq)

	sm2p256Mul(zInvSq, p.xyz[0:4], zInvSq)
	sm2p256Mul(zInv, p.xyz[4:8], zInv)

	sm2p256FromMont(zInvSq, zInvSq)
	sm2p256FromMont(zInv, zInv)

	xOut := make([]byte, 32)
	yOut := make([]byte, 32)
	sm2p256LittleToBig(xOut, zInvSq)
	sm2p256LittleToBig(yOut, zInv)

	return new(big.Int).SetBytes(xOut), new(big.Int).SetBytes(yOut)
}

// p256Inverse sets out to in^-1 mod p.
/*func p256Inverse(out, in []uint64) {
	var stack [6 * 4]uint64
	p2 := stack[4*0 : 4*0+4]
	p4 := stack[4*1 : 4*1+4]
	p8 := stack[4*2 : 4*2+4]
	p16 := stack[4*3 : 4*3+4]
	p32 := stack[4*4 : 4*4+4]

	sm2p256Sqr(out, in)//2^1
	sm2p256Mul(p2, out, in) // 2^2-2^0

	sm2p256Sqr(out, p2)
	sm2p256Sqr(out, out)
	sm2p256Mul(p4, out, p2) // f*p 2^4-2^0

	sm2p256Sqr(out, p4)
	sm2p256Sqr(out, out)
	sm2p256Sqr(out, out)
	sm2p256Sqr(out, out)
	sm2p256Mul(p8, out, p4) // ff*p 2^8-2^0

	sm2p256Sqr(out, p8)

	for i := 0; i < 7; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(p16, out, p8) // ffff*p 2^16-2^0

	sm2p256Sqr(out, p16)
	for i := 0; i < 15; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(p32, out, p16) // ffffffff*p 2^32-2^0

	sm2p256Sqr(out, p32)//ffffffffffffffff*p 2^64-2^0

	for i := 0; i < 31; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, in)

	for i := 0; i < 32*4; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p32)

	for i := 0; i < 32; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p32)

	for i := 0; i < 16; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p16)

	for i := 0; i < 8; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p8)

	sm2p256Sqr(out, out)
	sm2p256Sqr(out, out)
	sm2p256Sqr(out, out)
	sm2p256Sqr(out, out)
	sm2p256Mul(out, out, p4)

	sm2p256Sqr(out, out)
	sm2p256Sqr(out, out)
	sm2p256Mul(out, out, p2)

	sm2p256Sqr(out, out)
	sm2p256Sqr(out, out)
	sm2p256Mul(out, out, in)
}*/

// CopyConditional copies overwrites p with src if v == 1, and leaves p
// unchanged if v == 0.
func (p *p256Point) CopyConditional(src *p256Point, v int) {
	pMask := uint64(v) - 1
	srcMask := ^pMask

	for i, n := range p.xyz {
		p.xyz[i] = (n & pMask) | (src.xyz[i] & srcMask)
	}
}

func p256Inverse(out, in []uint64) {

	var stack [10 * 4]uint64
	p2 := stack[4*0 : 4*0+4]
	p4 := stack[4*1 : 4*1+4]
	p8 := stack[4*2 : 4*2+4]
	p16 := stack[4*3 : 4*3+4]
	p32 := stack[4*4 : 4*4+4]

	p3 := stack[4*5 : 4*5+4]
	p7 := stack[4*6 : 4*6+4]
	p15 := stack[4*7 : 4*7+4]
	p31 := stack[4*8 : 4*8+4]

	sm2p256Sqr(out, in) //2^1

	sm2p256Mul(p2, out, in) // 2^2-2^0
	sm2p256Sqr(out, p2)
	sm2p256Mul(p3, out, in)
	sm2p256Sqr(out, out)
	sm2p256Mul(p4, out, p2) // f*p 2^4-2^0

	sm2p256Sqr(out, p4)
	sm2p256Sqr(out, out)
	sm2p256Sqr(out, out)
	sm2p256Mul(p7, out, p3)
	sm2p256Sqr(out, out)
	sm2p256Mul(p8, out, p4) // ff*p 2^8-2^0

	sm2p256Sqr(out, p8)

	for i := 0; i < 6; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(p15, out, p7)
	sm2p256Sqr(out, out)
	sm2p256Mul(p16, out, p8) // ffff*p 2^16-2^0

	sm2p256Sqr(out, p16)
	for i := 0; i < 14; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(p31, out, p15)
	sm2p256Sqr(out, out)
	sm2p256Mul(p32, out, p16) // ffffffff*p 2^32-2^0

	//(2^31-1)*2^33+2^32-1
	sm2p256Sqr(out, p31)
	for i := 0; i < 32; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p32)

	//x*2^32+p32
	for i := 0; i < 32; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p32)
	//x*2^32+p32
	for i := 0; i < 32; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p32)
	//x*2^32+p32
	for i := 0; i < 32; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p32)
	//x*2^32
	for i := 0; i < 32; i++ {
		sm2p256Sqr(out, out)
	}

	//x*2^32+p32
	for i := 0; i < 32; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p32)

	//x*2^16+p16
	for i := 0; i < 16; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p16)

	//x*2^8+p8
	for i := 0; i < 8; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p8)

	//x*2^4+p4
	for i := 0; i < 4; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p4)

	//x*2^2+p2
	for i := 0; i < 2; i++ {
		sm2p256Sqr(out, out)
	}
	sm2p256Mul(out, out, p2)

	sm2p256Sqr(out, out)
	sm2p256Sqr(out, out)
	sm2p256Mul(out, out, in)
}

func (p *p256Point) p256StorePoint(r *[16 * 4 * 3]uint64, index int) {
	copy(r[index*12:], p.xyz[:])
}

func boothW5(in uint) (int, int) {
	var s uint = ^((in >> 5) - 1)
	var d uint = (1 << 6) - in - 1
	d = (d & s) | (in & (^s))
	d = (d >> 1) + (d & 1)
	return int(d), int(s & 1)
}

func boothW7(in uint) (int, int) {
	var s uint = ^((in >> 7) - 1)
	var d uint = (1 << 8) - in - 1
	d = (d & s) | (in & (^s))
	d = (d >> 1) + (d & 1)
	return int(d), int(s & 1)
}

func initTable() {
	p256Precomputed = new([37][64 * 8]uint64)

	/*	basePoint := []uint64{
		0x79e730d418a9143c, 0x75ba95fc5fedb601, 0x79fb732b77622510, 0x18905f76a53755c6,
		0xddf25357ce95560a, 0x8b4ab8e4ba19e45c, 0xd2e88688dd21f325, 0x8571ff1825885d85,
		0x0000000000000001, 0xffffffff00000000, 0xffffffffffffffff, 0x00000000fffffffe,
	}*/
	basePoint := []uint64{
		0x61328990F418029E, 0x3E7981EDDCA6C050, 0xD6A1ED99AC24C3C3, 0x91167A5EE1C13B05,
		0xC1354E593C2D0DDD, 0xC1F5E5788D3295FA, 0x8D4CFB066E2A48F8, 0x63CD65D481D735BD,
		0x0000000000000001, 0x00000000FFFFFFFF, 0x0000000000000000, 0x0000000100000000,
	}
	t1 := make([]uint64, 12)
	t2 := make([]uint64, 12)
	copy(t2, basePoint)

	zInv := make([]uint64, 4)
	zInvSq := make([]uint64, 4)
	for j := 0; j < 64; j++ {
		copy(t1, t2)
		for i := 0; i < 37; i++ {
			// The window size is 7 so we need to double 7 times.
			if i != 0 {
				for k := 0; k < 7; k++ {
					sm2p256PointDoubleAsm(t1, t1)
				}
			}
			// Convert the point to affine form. (Its values are
			// still in Montgomery form however.)
			p256Inverse(zInv, t1[8:12])
			sm2p256Sqr(zInvSq, zInv)
			sm2p256Mul(zInv, zInv, zInvSq)

			sm2p256Mul(t1[:4], t1[:4], zInvSq)
			sm2p256Mul(t1[4:8], t1[4:8], zInv)

			copy(t1[8:12], basePoint[8:12])
			// Update the table entry
			copy(p256Precomputed[i][j*8:], t1[:8])
		}
		if j == 0 {
			sm2p256PointDoubleAsm(t2, basePoint)
		} else {
			sm2p256PointAddAsm(t2, t2, basePoint)
		}
	}
}

func (p *p256Point) p256BaseMult(scalar []uint64) {
	precomputeOnce.Do(initTable)

	wvalue := (scalar[0] << 1) & 0xff
	sel, sign := boothW7(uint(wvalue))
	sm2p256SelectBase(p.xyz[0:8], p256Precomputed[0][0:], sel)
	sm2p256NegCond(p.xyz[4:8], sign)

	// (This is one, in the Montgomery domain.)
	//p.xyz[8] = 0x0000000000000001
	//p.xyz[9] = 0xffffffff00000000
	//p.xyz[10] = 0xffffffffffffffff
	//p.xyz[11] = 0x00000000fffffffe
	p.xyz[8] = 0x0000000000000001
	p.xyz[9] = 0x00000000FFFFFFFF
	p.xyz[10] = 0x0000000000000000
	p.xyz[11] = 0x0000000100000000
	var t0 p256Point
	// (This is one, in the Montgomery domain.)
	//t0.xyz[8] = 0x0000000000000001
	//t0.xyz[9] = 0xffffffff00000000
	//t0.xyz[10] = 0xffffffffffffffff
	//t0.xyz[11] = 0x00000000fffffffe
	t0.xyz[8] = 0x0000000000000001
	t0.xyz[9] = 0x00000000FFFFFFFF
	t0.xyz[10] = 0x0000000000000000
	t0.xyz[11] = 0x0000000100000000
	index := uint(6)
	zero := sel

	for i := 1; i < 37; i++ {
		if index < 192 {
			wvalue = ((scalar[index/64] >> (index % 64)) + (scalar[index/64+1] << (64 - (index % 64)))) & 0xff
		} else {
			wvalue = (scalar[index/64] >> (index % 64)) & 0xff
		}
		index += 7
		sel, sign = boothW7(uint(wvalue))
		sm2p256SelectBase(t0.xyz[0:8], p256Precomputed[i][0:], sel)
		sm2p256PointAddAffineAsm(p.xyz[0:12], p.xyz[0:12], t0.xyz[0:8], sign, sel, zero)
		zero |= sel
	}
}

func (p *p256Point) p256ScalarMult(scalar []uint64) {
	// precomp is a table of precomputed points that stores powers of p
	// from p^1 to p^16.
	var precomp [16 * 4 * 3]uint64
	var t0, t1, t2, t3 p256Point

	// Prepare the table
	p.p256StorePoint(&precomp, 0) // 1

	sm2p256PointDoubleAsm(t0.xyz[:], p.xyz[:])
	sm2p256PointDoubleAsm(t1.xyz[:], t0.xyz[:])
	sm2p256PointDoubleAsm(t2.xyz[:], t1.xyz[:])
	sm2p256PointDoubleAsm(t3.xyz[:], t2.xyz[:])
	t0.p256StorePoint(&precomp, 1)  // 2
	t1.p256StorePoint(&precomp, 3)  // 4
	t2.p256StorePoint(&precomp, 7)  // 8
	t3.p256StorePoint(&precomp, 15) // 16

	sm2p256PointAddAsm(t0.xyz[:], t0.xyz[:], p.xyz[:])
	sm2p256PointAddAsm(t1.xyz[:], t1.xyz[:], p.xyz[:])
	sm2p256PointAddAsm(t2.xyz[:], t2.xyz[:], p.xyz[:])
	t0.p256StorePoint(&precomp, 2) // 3
	t1.p256StorePoint(&precomp, 4) // 5
	t2.p256StorePoint(&precomp, 8) // 9

	sm2p256PointDoubleAsm(t0.xyz[:], t0.xyz[:])
	sm2p256PointDoubleAsm(t1.xyz[:], t1.xyz[:])
	t0.p256StorePoint(&precomp, 5) // 6
	t1.p256StorePoint(&precomp, 9) // 10

	sm2p256PointAddAsm(t2.xyz[:], t0.xyz[:], p.xyz[:])
	sm2p256PointAddAsm(t1.xyz[:], t1.xyz[:], p.xyz[:])
	t2.p256StorePoint(&precomp, 6)  // 7
	t1.p256StorePoint(&precomp, 10) // 11

	sm2p256PointDoubleAsm(t0.xyz[:], t0.xyz[:])
	sm2p256PointDoubleAsm(t2.xyz[:], t2.xyz[:])
	t0.p256StorePoint(&precomp, 11) // 12
	t2.p256StorePoint(&precomp, 13) // 14

	sm2p256PointAddAsm(t0.xyz[:], t0.xyz[:], p.xyz[:])
	sm2p256PointAddAsm(t2.xyz[:], t2.xyz[:], p.xyz[:])
	t0.p256StorePoint(&precomp, 12) // 13
	t2.p256StorePoint(&precomp, 14) // 15

	// Start scanning the window from top bit
	index := uint(254)
	var sel, sign int

	wvalue := (scalar[index/64] >> (index % 64)) & 0x3f
	sel, _ = boothW5(uint(wvalue))

	sm2p256Select(p.xyz[0:12], precomp[0:], sel)
	zero := sel

	for index > 4 {
		index -= 5
		sm2p256PointDoubleAsm(p.xyz[:], p.xyz[:])
		sm2p256PointDoubleAsm(p.xyz[:], p.xyz[:])
		sm2p256PointDoubleAsm(p.xyz[:], p.xyz[:])
		sm2p256PointDoubleAsm(p.xyz[:], p.xyz[:])
		sm2p256PointDoubleAsm(p.xyz[:], p.xyz[:])

		if index < 192 {
			wvalue = ((scalar[index/64] >> (index % 64)) + (scalar[index/64+1] << (64 - (index % 64)))) & 0x3f
		} else {
			wvalue = (scalar[index/64] >> (index % 64)) & 0x3f
		}

		sel, sign = boothW5(uint(wvalue))

		sm2p256Select(t0.xyz[0:], precomp[0:], sel)
		sm2p256NegCond(t0.xyz[4:8], sign)
		sm2p256PointAddAsm(t1.xyz[:], p.xyz[:], t0.xyz[:])
		sm2p256MovCond(t1.xyz[0:12], t1.xyz[0:12], p.xyz[0:12], sel)
		sm2p256MovCond(p.xyz[0:12], t1.xyz[0:12], t0.xyz[0:12], zero)
		zero |= sel
	}

	sm2p256PointDoubleAsm(p.xyz[:], p.xyz[:])
	sm2p256PointDoubleAsm(p.xyz[:], p.xyz[:])
	sm2p256PointDoubleAsm(p.xyz[:], p.xyz[:])
	sm2p256PointDoubleAsm(p.xyz[:], p.xyz[:])
	sm2p256PointDoubleAsm(p.xyz[:], p.xyz[:])

	wvalue = (scalar[0] << 1) & 0x3f
	sel, sign = boothW5(uint(wvalue))

	sm2p256Select(t0.xyz[0:], precomp[0:], sel)
	sm2p256NegCond(t0.xyz[4:8], sign)
	sm2p256PointAddAsm(t1.xyz[:], p.xyz[:], t0.xyz[:])
	sm2p256MovCond(t1.xyz[0:12], t1.xyz[0:12], p.xyz[0:12], sel)
	sm2p256MovCond(p.xyz[0:12], t1.xyz[0:12], t0.xyz[0:12], zero)
}

func Hexprint(in []byte) {
	for i := 0; i < len(in); i++ {
		fmt.Printf("%02x", in[i])
	}
	fmt.Println()
}

func AffineToP256Point(x, y *big.Int) (out p256Point) {
	z, _ := new(big.Int).SetString("0100000000000000000000000000000000FFFFFFFF0000000000000001", 16)
	p, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	tmpx, _ := new(big.Int).SetString("0", 16)
	tmpy, _ := new(big.Int).SetString("0", 16)
	tmpx.Mul(x, z)
	tmpy.Mul(y, z)
	tmpx.Mod(tmpx, p)
	tmpy.Mod(tmpy, p)
	fromBig(out.xyz[0:4], tmpx)
	fromBig(out.xyz[4:8], tmpy)
	fromBig(out.xyz[8:12], z)
	return out
}

func Uint64ToAffine(in []uint64) (x, y *big.Int) {
	var r p256Point
	for i := 0; i < 12; i++ {
		r.xyz[i] = in[i]
	}
	tmpx, tmpy := r.p256PointToAffine()
	return tmpx, tmpy
}

//precompute public key table
func (curve p256Curve) InitPubKeyTable(x,y *big.Int) (Precomputed *[37][64*8]uint64) {
	Precomputed = new([37][64 * 8]uint64)

	var r p256Point
	fromBig(r.xyz[0:4], maybeReduceModP(x))
	fromBig(r.xyz[4:8], maybeReduceModP(y))
	sm2p256Mul(r.xyz[0:4], r.xyz[0:4], rr[:])
	sm2p256Mul(r.xyz[4:8], r.xyz[4:8], rr[:])
	r.xyz[8] = 0x0000000000000001
	r.xyz[9] = 0x00000000FFFFFFFF
	r.xyz[10] = 0x0000000000000000
	r.xyz[11] = 0x0000000100000000
	basePoint := []uint64{
		r.xyz[0], r.xyz[1],r.xyz[2],r.xyz[3],
		r.xyz[4],r.xyz[5],r.xyz[6],r.xyz[7],
		r.xyz[8],r.xyz[9],r.xyz[10],r.xyz[11],
	}
	t1 := make([]uint64, 12)
	t2 := make([]uint64, 12)
	copy(t2, basePoint)

	zInv := make([]uint64, 4)
	zInvSq := make([]uint64, 4)
	for j := 0; j < 64; j++ {
		copy(t1, t2)
		for i := 0; i < 37; i++ {
			// The window size is 7 so we need to double 7 times.
			if i != 0 {
				for k := 0; k < 7; k++ {
					sm2p256PointDoubleAsm(t1, t1)
				}
			}
			// Convert the point to affine form. (Its values are
			// still in Montgomery form however.)
			p256Inverse(zInv, t1[8:12])
			sm2p256Sqr(zInvSq, zInv)
			sm2p256Mul(zInv, zInv, zInvSq)

			sm2p256Mul(t1[:4], t1[:4], zInvSq)
			sm2p256Mul(t1[4:8], t1[4:8], zInv)

			copy(t1[8:12], basePoint[8:12])
			// Update the table entry
			copy(Precomputed[i][j*8:], t1[:8])
		}
		if j == 0 {
			sm2p256PointDoubleAsm(t2, basePoint)
		} else {
			sm2p256PointAddAsm(t2, t2, basePoint)
		}
	}
	return
}

//fast sm2p256Mult with public key table
func (p *p256Point) p256PreMult(Precomputed *[37][64*8]uint64, scalar []uint64) {
	wvalue := (scalar[0] << 1) & 0xff
	sel, sign := boothW7(uint(wvalue))
	sm2p256SelectBase(p.xyz[0:8], Precomputed[0][0:], sel)
	sm2p256NegCond(p.xyz[4:8], sign)

	// (This is one, in the Montgomery domain.)
	//p.xyz[8] = 0x0000000000000001
	//p.xyz[9] = 0xffffffff00000000
	//p.xyz[10] = 0xffffffffffffffff
	//p.xyz[11] = 0x00000000fffffffe
	p.xyz[8] = 0x0000000000000001
	p.xyz[9] = 0x00000000FFFFFFFF
	p.xyz[10] = 0x0000000000000000
	p.xyz[11] = 0x0000000100000000
	var t0 p256Point
	// (This is one, in the Montgomery domain.)
	//t0.xyz[8] = 0x0000000000000001
	//t0.xyz[9] = 0xffffffff00000000
	//t0.xyz[10] = 0xffffffffffffffff
	//t0.xyz[11] = 0x00000000fffffffe
	t0.xyz[8] = 0x0000000000000001
	t0.xyz[9] = 0x00000000FFFFFFFF
	t0.xyz[10] = 0x0000000000000000
	t0.xyz[11] = 0x0000000100000000
	index := uint(6)
	zero := sel

	for i := 1; i < 37; i++ {
		if index < 192 {
			wvalue = ((scalar[index/64] >> (index % 64)) + (scalar[index/64+1] << (64 - (index % 64)))) & 0xff
		} else {
			wvalue = (scalar[index/64] >> (index % 64)) & 0xff
		}
		index += 7
		sel, sign = boothW7(uint(wvalue))
		sm2p256SelectBase(t0.xyz[0:8], Precomputed[i][0:], sel)
		sm2p256PointAddAffineAsm(p.xyz[0:12], p.xyz[0:12], t0.xyz[0:8], sign, sel, zero)
		zero |= sel
	}
}

//fast scalarmult with public key table
func (curve p256Curve) PreScalarMult(Precomputed *[37][64*8]uint64, scalar []byte) (x,y *big.Int) {
	scalarReversed := make([]uint64, 4)
	p256GetScalar(scalarReversed, scalar)

	r := new(p256Point)
	r.p256PreMult(Precomputed,scalarReversed)
	x,y = r.p256PointToAffine()
	return
}

//func TestP256_Point() {
//	fmt.Println("========================Test p256point==========================")
//	x, _ := new(big.Int).SetString("32C4AE2C1F1981195F9904466A39C9948FE30BBFF2660BE1715A4589334C74C7", 16)
//	y, _ := new(big.Int).SetString("BC3736A2F4F6779C59BDCEE36B692153D0A9877CC62A474002DF32E52139F0A0", 16)
//	r := AffineToP256Point(x, y)
//	x1, y1 := Uint64ToAffine(r.xyz[:])
//	fmt.Println("p256PointToAffine test:")
//	Hexprint(x1.Bytes())
//	Hexprint(y1.Bytes())
//
//	//r1:=AffineToP256Point(x,y)
//	fmt.Println("p256PointAdd test:")
//	res := make([]uint64, 12)
//	sm2p256PointAddAsm(res, r.xyz[0:12], r.xyz[0:12])
//	x2, y2 := Uint64ToAffine(res)
//	Hexprint(x2.Bytes())
//	Hexprint(y2.Bytes())
//
//	fmt.Println("p256PointDouble test:")
//
//	sm2p256PointDoubleAsm(r.xyz[:], r.xyz[:])
//	sm2p256PointDoubleAsm(r.xyz[:], r.xyz[:])
//	x3, y3 := Uint64ToAffine(r.xyz[:])
//	Hexprint(x3.Bytes())
//	Hexprint(y3.Bytes())
//}
//
//func Test_p256InternalFunc() {
//	fmt.Println("========================Test p256InternalFunction==========================")
//	x, _ := new(big.Int).SetString("6B17D1F2E12C4247F8BCE6E563A440F277037D812DEB33A0F4A13945D898C296", 16)
//	y, _ := new(big.Int).SetString("4FE342E2FE1A7F9B8EE7EB4A7C0F9E162BCE33576B315ECECBB6406837BF51F5", 16)
//	x1 := make([]uint64, 4)
//	y1 := make([]uint64, 4)
//	res := make([]uint64, 4)
//	fromBig(x1, x)
//	fromBig(y1, y)
//	//	fmt.Println(x1)
//	//	fmt.Println(y1)
//	sm2p256TestSubInternal(res, y1, x1)
//	res1 := make([]byte, 32)
//	sm2p256LittleToBig(res1, res)
//	fmt.Println("y-x mod p =")
//	Hexprint(res1)
//	fmt.Println("correct result =")
//	fmt.Println("E4CB70EF1CEE3D53962B0465186B5D23B4CAB5D53D462B2ED71507225F268F5E")
//
//	sm2p256TestMulInternal(res, x1, y1)
//	sm2p256LittleToBig(res1, res)
//	fmt.Println("x*y*2^-256 mod p =")
//	Hexprint(res1)
//	fmt.Println("correct result =")
//	fmt.Println("06A2D54FEDCC004C7A71BD0F44D288D4704F4BB8A8B155015971251DB3308ED8")
//
//	sm2p256TestMulBy2Inline(res, x1)
//	sm2p256LittleToBig(res1, res)
//	fmt.Println("x*2 mod p =")
//	Hexprint(res1)
//	fmt.Println("correct result =")
//	fmt.Println("D62FA3E5C258848FF179CDCAC74881E4EE06FB025BD66741E942728BB131852C")
//
//	sm2p256TestSqrInternal(res, x1)
//	sm2p256LittleToBig(res1, res)
//	fmt.Println("x^2*2^-256 mod p =")
//	Hexprint(res1)
//	fmt.Println("correct result =")
//	fmt.Println("8D7FA6B26A0495C732358637E3E851B017E632BBCE5D441F0BBB2B27B1B9C7C3")
//
//	sm2p256TestAddInline(res, x1, y1)
//	sm2p256LittleToBig(res1, res)
//	fmt.Println("x+y mod p =")
//	Hexprint(res1)
//	fmt.Println("correct result =")
//	fmt.Println("BAFB14D5DF46C1E387A4D22FDFB3DF08A2D1B0D8991C926FC05779AE1058148B")
//
//}
//
//func Test_p256Func() {
//	fmt.Println("========================Test p256Function==========================")
//	x, _ := new(big.Int).SetString("6B17D1F2E12C4247F8BCE6E563A440F277037D812DEB33A0F4A13945D898C296", 16)
//	y, _ := new(big.Int).SetString("4FE342E2FE1A7F9B8EE7EB4A7C0F9E162BCE33576B315ECECBB6406837BF51F5", 16)
//	x1 := make([]uint64, 4)
//	y1 := make([]uint64, 4)
//
//	res := make([]uint64, 4)
//	fromBig(x1, x)
//	fromBig(y1, y)
//	sm2p256Mul(res, x1, y1)
//	res1 := make([]byte, 32)
//	sm2p256LittleToBig(res1, res)
//	fmt.Println("x * y * 2^-256 mod p =")
//	Hexprint(res1)
//	fmt.Println("corect result:")
//	fmt.Println("06A2D54FEDCC004C7A71BD0F44D288D4704F4BB8A8B155015971251DB3308ED8")
//
//	sm2p256NegCond(res, 1)
//	sm2p256LittleToBig(res1, res)
//	fmt.Println("-x mod p =")
//	Hexprint(res1)
//	fmt.Println("corect result:")
//	fmt.Println("F95D2AAF1233FFB3858E42F0BB2D772B8FB0B446574EAAFFA68EDAE24CCF7127")
//
//	sm2p256Sqr(res, x1)
//	sm2p256LittleToBig(res1, res)
//	fmt.Println("x^2 * 2^-256 mod p =")
//	Hexprint(res1)
//	fmt.Println("corect result:")
//	fmt.Println("8D7FA6B26A0495C732358637E3E851B017E632BBCE5D441F0BBB2B27B1B9C7C3")
//
//	sm2p256FromMont(res, x1)
//	sm2p256LittleToBig(res1, res)
//	fmt.Println("x * 2^-256 mod p =")
//	Hexprint(res1)
//	fmt.Println("corect result:")
//	fmt.Println("A23CF658216A55A1C19D98E4F04219041D84F64CD13B0BADCCB2ED6D2259F06B")
//
//	sm2p256OrdMul(res, x1, y1)
//	sm2p256LittleToBig(res1, res)
//	fmt.Println("x * y * 2^-256 mod n =")
//	Hexprint(res1)
//	fmt.Println("corect result:")
//	fmt.Println("FEF2C8BAB69D869936CD281381435F925909E136C08350759D54DE6202B29D48")
//
//	sm2p256OrdSqr(res, x1, 1)
//	sm2p256LittleToBig(res1, res)
//	fmt.Println("x^2 * 2^-256 mod n =")
//	Hexprint(res1)
//	fmt.Println("corect result:")
//	fmt.Println("D5FD94D54F3640D652C6A56164B2DD5577D6A5EAA980FD231EFEB7AB23FC0E60")
//}
//
//func Test_amd64() {
//	fmt.Println("========================Test Amd64_Function==========================")
//	x, _ := new(big.Int).SetString("6B17D1F2E12C4247F8BCE6E563A440F277037D812DEB33A0F4A13945D898C296", 16)
//	x1 := make([]uint64, 4)
//	x_1 := make([]uint64, 4)
//	fromBig(x1, x)
//	p256Inverse(x_1, x1)
//	res1 := make([]byte, 32)
//	sm2p256LittleToBig(res1, x_1)
//	Hexprint(res1)
//}