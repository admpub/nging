package codec

import (
	"crypto/aes"
	"encoding/base64"
	"log"
)

func NewAesECBCrypto(keyTypes ...string) *AesECBCrypto {
	var keyType string
	if len(keyTypes) > 0 {
		keyType = keyTypes[0]
	}
	return &AesECBCrypto{key: make(map[string][]byte), keyType: keyType}
}

type AesECBCrypto struct {
	key     map[string][]byte
	keyType string
}

func (c *AesECBCrypto) aesKey(key []byte) []byte {
	if c.key == nil {
		c.key = make(map[string][]byte, 0)
	}
	ckey := string(key)
	k, ok := c.key[ckey]
	if !ok {
		k = AesGenKey(key, c.keyType)
		c.key[ckey] = k
	}
	return k
}

func (c *AesECBCrypto) Encode(rawData, authKey string) string {
	crypted := c.EncodeBytes([]byte(rawData), []byte(authKey))
	return base64.StdEncoding.EncodeToString(crypted)
}

func (c *AesECBCrypto) Decode(cryptedData, authKey string) string {
	crypted, err := base64.StdEncoding.DecodeString(cryptedData)
	if err != nil {
		log.Println(err)
		return ``
	}
	origData := c.DecodeBytes(crypted, []byte(authKey))
	return string(origData)
}

func (c *AesECBCrypto) EncodeBytes(rawData, authKey []byte) []byte {
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

	blockMode := NewECBEncrypter(block)
	crypted := make([]byte, len(in))
	blockMode.CryptBlocks(crypted, in)
	return crypted
}

func (c *AesECBCrypto) DecodeBytes(cryptedData, authKey []byte) []byte {
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
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(in))
	blockMode.CryptBlocks(origData, in)
	origData = PKCS5UnPadding(origData)
	return origData
}
