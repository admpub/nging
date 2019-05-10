/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package codec

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"log"
)

func NewDesCBCCrypto() *DesCBCCrypto {
	return &DesCBCCrypto{
		key: make(map[string][]byte),
	}
}

type DesCBCCrypto struct {
	key map[string][]byte
}

// Encode DES CBC加密
func (d *DesCBCCrypto) Encode(text, secret string) string {
	crypted := d.EncodeBytes([]byte(text), []byte(secret))
	return base64.StdEncoding.EncodeToString(crypted)
}

// Decode DES CBC解密
func (d *DesCBCCrypto) Decode(cryptedStr, secret string) string {
	crypted, err := base64.StdEncoding.DecodeString(cryptedStr)
	if err != nil {
		log.Println(err)
		return ``
	}
	origData := d.DecodeBytes(crypted, []byte(secret))
	return string(origData)
}

func (d *DesCBCCrypto) EncodeBytes(text, secret []byte) []byte {
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

func (d *DesCBCCrypto) DecodeBytes(crypted, secret []byte) []byte {
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

func (d *DesCBCCrypto) GenKey(key []byte) []byte {
	if d.key == nil {
		d.key = make(map[string][]byte, 0)
	}
	ckey := string(key)
	kkey, ok := d.key[ckey]
	if !ok {
		d.key[ckey] = DesGenKey(key)
	}
	return kkey
}
