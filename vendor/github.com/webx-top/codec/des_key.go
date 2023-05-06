package codec

import (
	"bytes"
	"crypto/des"
)

func newDESKey() *desKey {
	return &desKey{key: NewSafeKeys(), keyFixer: FixedDESKey}
}

type desKey struct {
	key      *SafeKeys
	keyFixer KeyFixer
}

func (a *desKey) SetKeyFixer(fixer KeyFixer) {
	a.keyFixer = fixer
}

func (a *desKey) GetKey(key []byte) []byte {
	ckey := string(key)
	k, ok := a.key.Get(ckey)
	if !ok {
		k = a.GenKey(key)
		a.key.Set(ckey, k)
	}
	return k
}

func (a *desKey) GenKey(key []byte) []byte {
	if a.keyFixer == nil {
		return FixedDESKey(des.BlockSize, key)
	}
	return a.keyFixer(des.BlockSize, key)
}

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
