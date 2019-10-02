package checksum

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
)

const bufferSize = 65536

// MD5sumReader returns MD5 checksum of filename
func MD5sumReader(reader io.Reader) (string, error) {
	return SumReader(md5.New(), reader)
}

// SHA256sumReader returns SHA256 checksum of filename
func SHA256sumReader(reader io.Reader) (string, error) {
	return SumReader(sha256.New(), reader)
}

// SumReader calculates the hash based on a provided hash provider
func SumReader(hashAlgorithm hash.Hash, reader io.Reader) (string, error) {
	buf := make([]byte, bufferSize)
	for {
		switch n, err := reader.Read(buf); err {
		case nil:
			hashAlgorithm.Write(buf[:n])
		case io.EOF:
			return fmt.Sprintf("%x", hashAlgorithm.Sum(nil)), nil
		default:
			return "", err
		}
	}
}
