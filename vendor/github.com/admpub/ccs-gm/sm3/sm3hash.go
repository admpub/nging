/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	SPDX-License-Identifier: Apache-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package sm3

func leftRotate(x uint32, r uint32) uint32 { return (x<<(r%32) | x>>(32-r%32)) }

func ff0(X uint32, Y uint32, Z uint32) uint32 { return X ^ Y ^ Z }
func ff1(X uint32, Y uint32, Z uint32) uint32 { return (X & Y) | (X & Z) | (Y & Z) }

func gg0(X uint32, Y uint32, Z uint32) uint32 { return X ^ Y ^ Z }
func gg1(X uint32, Y uint32, Z uint32) uint32 { return (X & Y) | ((^X) & Z) }

func p0(X uint32) uint32 { return X ^ leftRotate(X, 9) ^ leftRotate(X, 17) }

func p1(X uint32) uint32 { return X ^ leftRotate(X, 15) ^ leftRotate(X, 23) }

func msgPadding(message []byte) []byte {
	// Pre-processing:
	chunk := message

	// Pre-processing: adding a single 1 bit
	chunk = append(chunk, byte(0x80))

	// Pre-processing: padding with zeros
	padding := 56 - len(chunk)%64
	for i := 0; i < padding; i++ {
		chunk = append(chunk, 0x00)
	}
	var l uint64
	l = uint64(len(message) * 8)

	//	l := byte((len(message) * 8))()
	chunk = append(chunk, byte((l>>56)&0xff))
	chunk = append(chunk, byte((l>>48)&0xff))
	chunk = append(chunk, byte((l>>40)&0xff))
	chunk = append(chunk, byte((l>>32)&0xff))
	chunk = append(chunk, byte((l>>24)&0xff))
	chunk = append(chunk, byte((l>>16)&0xff))
	chunk = append(chunk, byte((l>>8)&0xff))
	chunk = append(chunk, byte(l&0xff))

	//	hstr := biu.BytesToHexString(chunk)
	//	fmt.Println(len(hstr))
	//	fmt.Println("test" + hstr)

	//	return hstr
	return chunk
}

type W struct {
	W1 [68]uint32
	W2 [64]uint32
}

func msgExp(x [16]uint32) W {
	var i int
	var wtmp W
	for i = 0; i < 16; i++ {
		wtmp.W1[i] = x[i]
	}
	for i = 16; i < 68; i++ {
		wtmp.W1[i] = p1(wtmp.W1[i-16]^wtmp.W1[i-9]^leftRotate(wtmp.W1[i-3], 15)) ^ leftRotate(wtmp.W1[i-13], 7) ^ wtmp.W1[i-6]
	}
	for i = 0; i < 64; i++ {
		wtmp.W2[i] = wtmp.W1[i] ^ wtmp.W1[i+4]
	}
	return wtmp
}

func cF(V [8]uint32, Bmsg [16]uint32) [8]uint32 {
	var j int
	var A, B, C, D, E, F, G, H uint32
	A = V[0]
	B = V[1]
	C = V[2]
	D = V[3]
	E = V[4]
	F = V[5]
	G = V[6]
	H = V[7]
	wtmp := msgExp(Bmsg)
	for j = 0; j < 16; j++ {
		var jj int
		if j < 33 {
			jj = j
		} else {
			jj = j - 32
		}
		SS1 := leftRotate(leftRotate(A, 12)+E+leftRotate(0x79cc4519, uint32(jj)), 7)
		SS2 := SS1 ^ leftRotate(A, 12)
		TT1 := ff0(A, B, C) + D + SS2 + wtmp.W2[j]
		TT2 := gg0(E, F, G) + H + SS1 + wtmp.W1[j]
		D = C
		C = leftRotate(B, 9)
		B = A
		A = TT1
		H = G
		G = leftRotate(F, 19)
		F = E
		E = p0(TT2)
	}
	for j = 16; j < 64; j++ {
		var jj int
		if j < 33 {
			jj = j
		} else {
			jj = j - 32
		}
		SS1 := leftRotate(leftRotate(A, 12)+E+leftRotate(0x7a879d8a, uint32(jj)), 7)
		SS2 := SS1 ^ leftRotate(A, 12)
		TT1 := ff1(A, B, C) + D + SS2 + wtmp.W2[j]
		TT2 := gg1(E, F, G) + H + SS1 + wtmp.W1[j]
		D = C
		C = leftRotate(B, 9)
		B = A
		A = TT1
		H = G
		G = leftRotate(F, 19)
		F = E
		E = p0(TT2)
	}

	V[0] = A ^ V[0]
	V[1] = B ^ V[1]
	V[2] = C ^ V[2]
	V[3] = D ^ V[3]
	V[4] = E ^ V[4]
	V[5] = F ^ V[5]
	V[6] = G ^ V[6]
	V[7] = H ^ V[7]

	return V
}

func Block(dig *digest, p []byte) {
	var V [8]uint32
	for i := 0; i < 8; i++ {
		V[i] = dig.h[i]
	}
	for len(p) >= 64 {
		m := [16]uint32{}
		x := p[:64]
		xi := 0
		mi := 0
		for mi < 16 {
			m[mi] = (uint32(x[xi+3]) |
				(uint32(x[xi+2]) << 8) |
				(uint32(x[xi+1]) << 16) |
				(uint32(x[xi]) << 24))
			mi += 1
			xi += 4
		}
		V = cF(V, m)
		p = p[64:]
	}
	for i := 0; i < 8; i++ {
		dig.h[i] = V[i]
	}
}