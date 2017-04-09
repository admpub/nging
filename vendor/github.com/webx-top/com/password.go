/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package com

import (
	"crypto/hmac"
	"crypto/sha1"
	"hash"
)

// Hash 生成哈希值
func Hash(str string) string {
	return Sha256(str)
}

// Salt 盐值加密(生成64个字符)
func Salt() string {
	return Hash(RandStr(64))
}

// MakePassword 创建密码(生成64个字符)
// 可以指定positions用来在hash处理后的密码的不同位置插入salt片段(数量取决于positions的数量)，然后再次hash
func MakePassword(password string, salt string, positions ...uint) string {
	length := len(positions)
	if length < 1 {
		return Hash(salt + password)
	}
	saltLength := len(salt)
	if saltLength < length {
		return Hash(salt + password)
	}
	saltChars := saltLength / length
	hashedPassword := Hash(password)
	maxIndex := len(hashedPassword) - 1
	saltMaxIndex := saltLength - 1
	var result string
	var lastPos int
	for k, pos := range positions {
		end := int(pos)
		start := lastPos
		if start > end {
			start, end = end, start
		}
		if start > maxIndex {
			continue
		}
		lastPos = end
		saltStart := k * saltChars
		saltEnd := saltStart + saltChars
		if end > maxIndex {
			result += hashedPassword[start:] + salt[saltStart:saltEnd]
			continue
		}
		result += hashedPassword[start:end] + salt[saltStart:saltEnd]
		if k == length-1 {
			if end <= maxIndex {
				result += hashedPassword[end:]
			}
			if saltEnd <= saltMaxIndex {
				result += salt[saltEnd:]
			}
		}
	}
	return Hash(result)
}

// CheckPassword 检查密码(密码原文，数据库中保存的哈希过后的密码，数据库中保存的盐值)
func CheckPassword(rawPassword string, hashedPassword string, salt string, positions ...uint) bool {
	return MakePassword(rawPassword, salt, positions...) == hashedPassword
}

// PBKDF2Key derives a key from the password, salt and iteration count, returning a
// []byte of length keylen that can be used as cryptographic key. The key is
// derived based on the method described as PBKDF2 with the HMAC variant using
// the supplied hash function.
//
// For example, to use a HMAC-SHA-1 based PBKDF2 key derivation function, you
// can get a derived key for e.g. AES-256 (which needs a 32-byte key) by
// doing:
//
// 	dk := pbkdf2.Key([]byte("some password"), salt, 4096, 32, sha1.New)
//
// Remember to get a good random salt. At least 8 bytes is recommended by the
// RFC.
//
// Using a higher iteration count will increase the cost of an exhaustive
// search but will also make derivation proportionally slower.
func PBKDF2Key(password, salt []byte, iter, keyLen int, hFunc ...func() hash.Hash) string {
	var h func() hash.Hash
	if len(hFunc) > 0 {
		h = hFunc[0]
	}
	if h == nil {
		h = sha1.New
	}
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return Base64Encode(string(dk[:keyLen]))
}
