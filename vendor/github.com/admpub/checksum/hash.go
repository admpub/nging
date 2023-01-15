package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"hash/crc32"
	"io"

	"golang.org/x/crypto/blake2s"
)

const bufferSize = 65536

// MD5sumReader returns MD5 checksum of content in reader
func MD5sumReader(reader io.Reader) (string, error) {
	return sumReader(md5.New(), reader)
}

// SHA256sumReader returns SHA256 checksum of content in reader
func SHA256sumReader(reader io.Reader) (string, error) {
	return sumReader(sha256.New(), reader)
}

// SHA1sumReader returns SHA1 checksum of content in reader
func SHA1sumReader(reader io.Reader) (string, error) {
	return sumReader(sha1.New(), reader)
}

// Blake2s256Reader returns SHA1 checksum of content in reader
func Blake2s256Reader(reader io.Reader) (string, error) {
	hash, _ := blake2s.New256([]byte{})
	return sumReader(hash, reader)
}

// CRCReader returns CRC-32-IEEE checksum of content in reader
func CRCReader(reader io.Reader) (string, error) {
	table := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum([]byte(""), table)
	buf := make([]byte, bufferSize)
	for {
		switch n, err := reader.Read(buf); err {
		case nil:
			checksum = crc32.Update(checksum, table, buf[:n])
		case io.EOF:
			return fmt.Sprintf("%x", checksum), nil
		default:
			return "", err
		}
	}
}

// sumReader calculates the hash based on a provided hash provider
func sumReader(hashAlgorithm hash.Hash, reader io.Reader) (string, error) {
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
