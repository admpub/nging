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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
)

type AesCrypto struct {
	key map[string][]byte
}

func NewAesCrypto() *AesCrypto {
	return &AesCrypto{key: make(map[string][]byte)}
}

const (
	aesKeyLen = 128
	keyLen    = aesKeyLen / 8
)

func (c *AesCrypto) aesKey(key []byte) []byte {
	if c.key == nil {
		c.key = make(map[string][]byte, 0)
	}
	ckey := string(key)
	k, ok := c.key[ckey]
	if !ok {
		if len(key) == keyLen {
			return key
		}

		k = make([]byte, keyLen)
		copy(k, key)
		for i := keyLen; i < len(key); {
			for j := 0; j < keyLen && i < len(key); j, i = j+1, i+1 {
				k[j] ^= key[i]
			}
		}
		c.key[ckey] = k
	}
	return k
}

func (c *AesCrypto) Encode(rawData, authKey string) string {
	crypted := c.EncodeBytes([]byte(rawData), []byte(authKey))
	return base64.StdEncoding.EncodeToString(crypted)
}

func (c *AesCrypto) Decode(cryptedData, authKey string) string {
	crypted, err := base64.StdEncoding.DecodeString(cryptedData)
	if err != nil {
		log.Println(err)
		return ``
	}
	origData := c.DecodeBytes(crypted, []byte(authKey))
	return string(origData)
}

func (c *AesCrypto) EncodeBytes(rawData, authKey []byte) []byte {
	in := rawData
	key := authKey
	key = c.aesKey(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}
	blockSize := block.BlockSize()
	in = PKCS5Padding(in, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(in))
	blockMode.CryptBlocks(crypted, in)
	return crypted
}

func (c *AesCrypto) DecodeBytes(cryptedData, authKey []byte) []byte {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	in := cryptedData
	key := authKey
	key = c.aesKey(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(in))
	blockMode.CryptBlocks(origData, in)
	origData = PKCS5UnPadding(origData)
	return origData
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length > unpadding {
		return origData[:(length - unpadding)]
	}
	return []byte{}
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length > unpadding {
		return origData[:(length - unpadding)]
	}
	return []byte{}
}
