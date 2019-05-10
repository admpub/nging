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
	"crypto/des"
	"encoding/base64"
	"log"
)

func NewDesECBCrypto() *DesECBCrypto {
	return &DesECBCrypto{
		key: make(map[string][]byte),
	}
}

type DesECBCrypto struct {
	key map[string][]byte
}

func (d *DesECBCrypto) Encode(text, secret string) string {
	out := d.EncodeBytes([]byte(text), []byte(secret))
	return base64.StdEncoding.EncodeToString(out)
}

func (d *DesECBCrypto) Decode(crypted, secret string) string {
	data, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		log.Println(err)
		return ``
	}
	out := d.DecodeBytes(data, []byte(secret))
	return string(out)
}

func (d *DesECBCrypto) EncodeBytes(text, secret []byte) []byte {
	data := text
	key := secret
	key = d.GenKey(key)
	block, err := des.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}
	bs := block.BlockSize()
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		log.Println("Need a multiple of the blocksize")
		return nil
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out
}

func (d *DesECBCrypto) DecodeBytes(crypted, secret []byte) []byte {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	data := crypted
	key := secret
	key = d.GenKey(key)
	block, err := des.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		log.Println("crypto/cipher: input not full blocks")
		return nil
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = PKCS5UnPadding(out)
	return out
}

func (d *DesECBCrypto) GenKey(key []byte) []byte {
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
