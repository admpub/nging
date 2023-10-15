package securecookie

import (
	"crypto/cipher"
	"hash"
)

func SetMaxLength(codecs []Codec, l int) {
	for _, c := range codecs {
		if codec, ok := c.(*SecureCookie); ok {
			codec.MaxLength(l)
		}
	}
}

func SetMaxAge(codecs []Codec, l int) {
	for _, c := range codecs {
		if codec, ok := c.(*SecureCookie); ok {
			codec.MaxAge(l)
		}
	}
}

func SetMinAge(codecs []Codec, l int) {
	for _, c := range codecs {
		if codec, ok := c.(*SecureCookie); ok {
			codec.MinAge(l)
		}
	}
}

func SetHashFunc(codecs []Codec, f func() hash.Hash) {
	for _, c := range codecs {
		if codec, ok := c.(*SecureCookie); ok {
			codec.HashFunc(f)
		}
	}
}

func SetBlockFunc(codecs []Codec, f func([]byte) (cipher.Block, error)) {
	for _, c := range codecs {
		if codec, ok := c.(*SecureCookie); ok {
			codec.BlockFunc(f)
		}
	}
}

func SetSerializer(codecs []Codec, sz Serializer) {
	for _, c := range codecs {
		if codec, ok := c.(*SecureCookie); ok {
			codec.SetSerializer(sz)
		}
	}
}
