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

package codec

import "strings"

func NewAES(keyTypes ...string) *AES {
	c := &AES{}
	var keyType string
	if len(keyTypes) > 0 {
		keyType = keyTypes[0]
	}
	if len(keyType) > 0 { // AES-128-CBC
		args := strings.SplitN(keyType, `-`, 3)
		if len(args) == 3 {
			switch args[2] {
			case `CBC`:
				keyType := strings.Join(args[0:2], `-`)
				c.Codec = NewAESCBC(keyType)
			case `ECB`:
				keyType := strings.Join(args[0:2], `-`)
				c.Codec = NewAESECB(keyType)
			default:
				panic("Unsuppored: " + keyType)
			}
		}
	}
	if c.Codec == nil {
		c.Codec = NewAESCBC(keyType)
	}
	return c
}

type AES struct {
	Codec
}

const (
	aes128KeyLen = 128 / 8 // 16
	aes192KeyLen = 192 / 8 // 24
	aes256KeyLen = 256 / 8 // 32
)

// AESKeyTypes AES Key类型
var AESKeyTypes = map[string]int{
	`AES-128`: aes128KeyLen,
	`AES-192`: aes192KeyLen,
	`AES-256`: aes256KeyLen,
}

func GenAESKey(key []byte, typ ...string) []byte {
	var keyType string
	if len(typ) > 0 {
		keyType = typ[0]
	}
	keyLen, ok := AESKeyTypes[keyType]
	if !ok {
		keyLen = aes128KeyLen
	}
	if len(key) == keyLen {
		return key
	}

	k := make([]byte, keyLen)
	copy(k, key)
	for i := keyLen; i < len(key); {
		for j := 0; j < keyLen && i < len(key); j, i = j+1, i+1 {
			k[j] ^= key[i]
		}
	}
	return k
}
