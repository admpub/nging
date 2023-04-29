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

import (
	"bytes"
	"crypto/des"
)

var (
	_ Codec = NewDESCBC()
	_ Codec = NewDESECB()
)

func GenDESKey(key []byte) []byte {
	return FixedDESKey(des.BlockSize, key)
}

var FixedDESKey = FixedKeyByRepeatContent

func FixedKeyByRepeatContent(keyLen int, key []byte) []byte {
	if len(key) == keyLen {
		return key
	}
	k := make([]byte, keyLen)
	length := len(key)
	if length == 0 {
		copy(k, bytes.Repeat([]byte(` `), keyLen))
	} else if length < keyLen {
		div := keyLen / length
		mod := keyLen % length
		for i := 0; i < div; i++ {
			copy(k[length*i:], key)
		}
		copy(k[length*div:], key[:mod])
	} else {
		copy(k, key)
	}
	return k
}
