package codec

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5str(v string) string {
	m := md5.New()
	io.WriteString(m, v)
	return hex.EncodeToString(m.Sum(nil))
}

func Md5bytes(v []byte) []byte {
	m := md5.New()
	m.Write(v)
	value := m.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(value)))
	hex.Encode(dst, value)
	return dst
}
