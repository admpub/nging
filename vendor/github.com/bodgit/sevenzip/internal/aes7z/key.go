package aes7z

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func calculateKey(password string, cycles int, salt []byte) []byte {
	b := bytes.NewBuffer(salt)

	// Convert password to UTF-16LE
	utf16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	t := transform.NewWriter(b, utf16le.NewEncoder())
	_, _ = t.Write([]byte(password))

	key := make([]byte, sha256.Size)
	if cycles == 0x3f {
		copy(key, b.Bytes())
	} else {
		h := sha256.New()
		for i := uint64(0); i < 1<<cycles; i++ {
			// These will never error
			_, _ = h.Write(b.Bytes())
			_ = binary.Write(h, binary.LittleEndian, i)
		}
		copy(key, h.Sum(nil))
	}

	return key
}
