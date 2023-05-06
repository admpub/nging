package codec

import (
	"crypto/aes"
	"encoding/base64"
	"log"
)

func NewAESECB(keyTypes ...string) *AESECB {
	return &AESECB{aesKey: newAESKey(keyTypes...)}
}

type AESECB struct {
	*aesKey
}

func (c *AESECB) genKey(key []byte) []byte {
	if c.aesKey == nil {
		c.aesKey = newAESKey()
	}
	return c.GetKey(key)
}

func (c *AESECB) Encode(rawData, authKey string) string {
	crypted := c.EncodeBytes([]byte(rawData), []byte(authKey))
	return base64.StdEncoding.EncodeToString(crypted)
}

func (c *AESECB) Decode(cryptedData, authKey string) string {
	crypted, err := base64.StdEncoding.DecodeString(cryptedData)
	if err != nil {
		log.Println(err)
		return ``
	}
	origData := c.DecodeBytes(crypted, []byte(authKey))
	return string(origData)
}

func (c *AESECB) EncodeBytes(rawData, authKey []byte) []byte {
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

	blockMode := NewECBEncrypter(block)
	crypted := make([]byte, len(in))
	blockMode.CryptBlocks(crypted, in)
	return crypted
}

func (c *AESECB) DecodeBytes(cryptedData, authKey []byte) []byte {
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
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(in))
	blockMode.CryptBlocks(origData, in)
	origData = PKCS5UnPadding(origData)
	return origData
}
