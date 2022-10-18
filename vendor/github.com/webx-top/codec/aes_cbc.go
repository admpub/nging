package codec

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
)

func NewAESCBC(keyTypes ...string) *AESCBC {
	var keyType string
	if len(keyTypes) > 0 {
		keyType = keyTypes[0]
	}
	return &AESCBC{key: NewSafeKeys(), keyType: keyType}
}

type AESCBC struct {
	key     *SafeKeys
	keyType string
}

func (c *AESCBC) genKey(key []byte) []byte {
	if c.key == nil {
		c.key = NewSafeKeys()
	}
	ckey := string(key)
	k, ok := c.key.Get(ckey)
	if !ok {
		k = GenAESKey(key, c.keyType)
		c.key.Set(ckey, k)
	}
	return k
}

func (c *AESCBC) Encode(rawData, authKey string) string {
	crypted := c.EncodeBytes([]byte(rawData), []byte(authKey))
	return base64.StdEncoding.EncodeToString(crypted)
}

func (c *AESCBC) Decode(cryptedData, authKey string) string {
	crypted, err := base64.StdEncoding.DecodeString(cryptedData)
	if err != nil {
		log.Println(err)
		return ``
	}
	origData := c.DecodeBytes(crypted, []byte(authKey))
	return string(origData)
}

func (c *AESCBC) EncodeBytes(rawData, authKey []byte) []byte {
	in := rawData
	key := authKey
	key = c.genKey(key)
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

func (c *AESCBC) DecodeBytes(cryptedData, authKey []byte) []byte {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	in := cryptedData
	key := authKey
	key = c.genKey(key)
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
