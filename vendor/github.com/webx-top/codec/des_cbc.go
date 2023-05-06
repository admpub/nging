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
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"log"
)

func NewDESCBC() *DESCBC {
	return &DESCBC{
		desKey: newDESKey(),
	}
}

type DESCBC struct {
	*desKey
}

// Encode DES CBC加密
func (d *DESCBC) Encode(text, secret string) string {
	crypted := d.EncodeBytes([]byte(text), []byte(secret))
	return base64.StdEncoding.EncodeToString(crypted)
}

// Decode DES CBC解密
func (d *DESCBC) Decode(cryptedStr, secret string) string {
	crypted, err := base64.StdEncoding.DecodeString(cryptedStr)
	if err != nil {
		log.Println(err)
		return ``
	}
	origData := d.DecodeBytes(crypted, []byte(secret))
	return string(origData)
}

func (d *DESCBC) EncodeBytes(text, secret []byte) []byte {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	key := secret
	key = d.GenKey(key)
	block, err := des.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}
	origData := text
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted
}

func (d *DESCBC) DecodeBytes(crypted, secret []byte) []byte {
	key := secret
	key = d.GenKey(key)
	block, err := des.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData
}

func (d *DESCBC) GenKey(key []byte) []byte {
	if d.desKey == nil {
		d.desKey = newDESKey()
	}
	return d.GetKey(key)
}
