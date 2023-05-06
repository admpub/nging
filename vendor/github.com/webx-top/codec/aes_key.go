package codec

import "bytes"

const (
	KeyAES128 = `AES-128`
	KeyAES192 = `AES-192`
	KeyAES256 = `AES-256`
)

const (
	aes128KeyLen = 128 / 8 // 16
	aes192KeyLen = 192 / 8 // 24
	aes256KeyLen = 256 / 8 // 32
)

// AESKeyTypes AES Key类型
var AESKeyTypes = map[string]int{
	KeyAES128: aes128KeyLen,
	KeyAES192: aes192KeyLen,
	KeyAES256: aes256KeyLen,
}

type KeyFixer func(keyLen int, key []byte) []byte

func newAESKey(keyTypes ...string) *aesKey {
	var keyType string
	if len(keyTypes) > 0 {
		keyType = keyTypes[0]
	}
	return &aesKey{key: NewSafeKeys(), keyType: keyType, keyFixer: FixedAESKey}
}

type aesKey struct {
	key      *SafeKeys
	keyType  string
	keyFixer KeyFixer
}

func (a *aesKey) SetKeyFixer(fixer KeyFixer) {
	a.keyFixer = fixer
}

func (a *aesKey) GetKey(key []byte) []byte {
	ckey := string(key)
	k, ok := a.key.Get(ckey)
	if !ok {
		k = a.GenKey(key, a.keyType)
		a.key.Set(ckey, k)
	}
	return k
}

func (a *aesKey) GenKey(key []byte, typ ...string) []byte {
	var keyType string
	if len(typ) > 0 {
		keyType = typ[0]
	}
	keyLen, ok := AESKeyTypes[keyType]
	if !ok {
		keyLen = aes128KeyLen
	}
	if a.keyFixer == nil {
		return FixedAESKey(keyLen, key)
	}
	return a.keyFixer(keyLen, key)
}

var FixedAESKey = FixedKeyDefault

func FixedKeyDefault(keyLen int, key []byte) []byte {
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

func FixedKeyByWhitespacePrefix(keyLen int, key []byte) []byte {
	if len(key) == keyLen {
		return key
	}
	k := make([]byte, keyLen)
	if len(key) < keyLen {
		remains := keyLen - len(key)
		copy(k, bytes.Repeat([]byte(` `), remains))
		copy(k[remains:], key)
	} else {
		copy(k, key)
	}
	return k
}
