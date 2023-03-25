// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x509

// RFC 1423 describes the encryption of PEM blocks. The algorithm used to
// generate a key from the password was derived by looking at the OpenSSL
// implementation.

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"hash"
	"io"
	"reflect"
)

type PEMCipher int

// Possible values for the EncryptPEMBlock encryption algorithm.
const (
	_ PEMCipher = iota
	PEMCipherDES
	PEMCipher3DES
	PEMCipherAES128
	PEMCipherAES192
	PEMCipherAES256
)

/*
 * reference to RFC5959 and RFC2898
 */

var (
	oidPBES1  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 3}  // pbeWithMD5AndDES-CBC(PBES1)
	oidPBES2  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 13} // id-PBES2(PBES2)
	oidPBKDF2 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 12} // id-PBKDF2

	oidKEYMD5    = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 5}
	oidKEYSHA1   = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 7}
	oidKEYSHA256 = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 9}
	oidKEYSHA512 = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 11}

	oidAES128CBC = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 2}
	oidAES256CBC = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 42}

	oidSM2 = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
)

// reference to https://www.rfc-editor.org/rfc/rfc5958.txt
type PrivateKeyInfo struct {
	Version             int // v1 or v2
	PrivateKeyAlgorithm []asn1.ObjectIdentifier
	PrivateKey          []byte
}

// reference to https://www.rfc-editor.org/rfc/rfc5958.txt
type EncryptedPrivateKeyInfo struct {
	EncryptionAlgorithm Pbes2Algorithms
	EncryptedData       []byte
}

// reference to https://www.ietf.org/rfc/rfc2898.txt
type Pbes2Algorithms struct {
	IdPBES2     asn1.ObjectIdentifier
	Pbes2Params Pbes2Params
}

// reference to https://www.ietf.org/rfc/rfc2898.txt
type Pbes2Params struct {
	KeyDerivationFunc Pbes2KDfs // PBES2-KDFs
	EncryptionScheme  Pbes2Encs // PBES2-Encs
}

// reference to https://www.ietf.org/rfc/rfc2898.txt
type Pbes2KDfs struct {
	IdPBKDF2    asn1.ObjectIdentifier
	Pkdf2Params Pkdf2Params
}

type Pbes2Encs struct {
	EncryAlgo asn1.ObjectIdentifier
	IV        []byte
}

// reference to https://www.ietf.org/rfc/rfc2898.txt
type Pkdf2Params struct {
	Salt           []byte
	IterationCount int
	Prf            pkix.AlgorithmIdentifier
}

// rfc1423Algo holds a method for enciphering a PEM block.
type rfc1423Algo struct {
	cipher     PEMCipher
	name       string
	cipherFunc func(key []byte) (cipher.Block, error)
	keySize    int
	blockSize  int
}

// rfc1423Algos holds a slice of the possible ways to encrypt a PEM
// block. The ivSize numbers were taken from the OpenSSL source.
var rfc1423Algos = []rfc1423Algo{{
	cipher:     PEMCipherDES,
	name:       "DES-CBC",
	cipherFunc: des.NewCipher,
	keySize:    8,
	blockSize:  des.BlockSize,
}, {
	cipher:     PEMCipher3DES,
	name:       "DES-EDE3-CBC",
	cipherFunc: des.NewTripleDESCipher,
	keySize:    24,
	blockSize:  des.BlockSize,
}, {
	cipher:     PEMCipherAES128,
	name:       "AES-128-CBC",
	cipherFunc: aes.NewCipher,
	keySize:    16,
	blockSize:  aes.BlockSize,
}, {
	cipher:     PEMCipherAES192,
	name:       "AES-192-CBC",
	cipherFunc: aes.NewCipher,
	keySize:    24,
	blockSize:  aes.BlockSize,
}, {
	cipher:     PEMCipherAES256,
	name:       "AES-256-CBC",
	cipherFunc: aes.NewCipher,
	keySize:    32,
	blockSize:  aes.BlockSize,
},
}

// deriveKey uses a key derivation function to stretch the password into a key
// with the number of bits our cipher requires. This algorithm was derived from
// the OpenSSL source.
func (c rfc1423Algo) deriveKey(password, salt []byte) []byte {
	hash := md5.New()
	out := make([]byte, c.keySize)
	var digest []byte

	for i := 0; i < len(out); i += len(digest) {
		hash.Reset()
		hash.Write(digest)
		hash.Write(password)
		hash.Write(salt)
		digest = hash.Sum(digest[:0])
		copy(out[i:], digest)
	}
	return out
}

// IsEncryptedPEMBlock returns if the PEM block is password encrypted.
func IsEncryptedPEMBlock(b *pem.Block) bool {
	_, ok := b.Headers["DEK-Info"]
	return ok
}

// IncorrectPasswordError is returned when an incorrect password is detected.
var IncorrectPasswordError = errors.New("x509: decryption password incorrect")

// DecryptPEMBlock takes a password encrypted PEM block and the password used to
// encrypt it and returns a slice of decrypted DER encoded bytes. It inspects
// the DEK-Info header to determine the algorithm used for decryption. If no
// DEK-Info header is present, an error is returned. If an incorrect password
// is detected an IncorrectPasswordError is returned. Because of deficiencies
// in the encrypted-PEM format, it's not always possible to detect an incorrect
// password. In these cases no error will be returned but the decrypted DER
// bytes will be random noise.
func DecryptPEMBlock(b *pem.Block, password []byte) ([]byte, error) {
	var keyInfo EncryptedPrivateKeyInfo

	_, err := asn1.Unmarshal(b.Bytes, &keyInfo)
	if err != nil {
		return nil, errors.New("x509: unknown format")
	}
	if !reflect.DeepEqual(keyInfo.EncryptionAlgorithm.IdPBES2, oidPBES2) {
		return nil, errors.New("x509: only support PBES2")
	}
	encryptionScheme := keyInfo.EncryptionAlgorithm.Pbes2Params.EncryptionScheme
	keyDerivationFunc := keyInfo.EncryptionAlgorithm.Pbes2Params.KeyDerivationFunc
	if !reflect.DeepEqual(keyDerivationFunc.IdPBKDF2, oidPBKDF2) {
		return nil, errors.New("x509: only support PBKDF2")
	}
	pkdf2Params := keyDerivationFunc.Pkdf2Params
	if !reflect.DeepEqual(encryptionScheme.EncryAlgo, oidAES128CBC) &&
		!reflect.DeepEqual(encryptionScheme.EncryAlgo, oidAES256CBC) {
		return nil, errors.New("x509: unknow encryption algorithm")
	}
	iv := encryptionScheme.IV
	salt := pkdf2Params.Salt
	iter := pkdf2Params.IterationCount
	encryptedKey := keyInfo.EncryptedData
	var key []byte
	switch {
	case pkdf2Params.Prf.Algorithm.Equal(oidKEYMD5):
		key = pbkdf(password, salt, iter, 32, md5.New)
		break
	case pkdf2Params.Prf.Algorithm.Equal(oidKEYSHA1):
		key = pbkdf(password, salt, iter, 32, sha1.New)
		break
	case pkdf2Params.Prf.Algorithm.Equal(oidKEYSHA256):
		key = pbkdf(password, salt, iter, 32, sha256.New)
		break
	case pkdf2Params.Prf.Algorithm.Equal(oidKEYSHA512):
		key = pbkdf(password, salt, iter, 32, sha512.New)
		break
	default:
		return nil, errors.New("x509: unknown hash algorithm")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptedKey, encryptedKey)

	//un-padding
	dataLen := len(encryptedKey)
	padLen := int(encryptedKey[dataLen-1])
	//check the padLen
	if dataLen <= padLen {
		return nil, errors.New("padding info incorrect")
	}
	for i := 0; i < padLen; i++ {
		if int(encryptedKey[dataLen-padLen+i]) != padLen {
			return nil, errors.New("padding info incorrect")
		}
	}
	return encryptedKey[:dataLen-padLen], nil
}

// EncryptPEMBlock returns a PEM block of the specified type holding the
// given DER-encoded data encrypted with the specified algorithm and
// password.
func EncryptPEMBlock(rand io.Reader, blockType string, data, password []byte, alg PEMCipher) (*pem.Block, error) {
	ciph := cipherByKey(alg)
	if ciph == nil {
		return nil, errors.New("x509: unknown encryption mode")
	}
	iter := 2048
	salt := make([]byte, 8)
	iv := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}
	key := pbkdf(password, salt, iter, 32, sha1.New) // SHA1
	padding := aes.BlockSize - len(data)%aes.BlockSize
	if padding > 0 {
		n := len(data)
		data = append(data, make([]byte, padding)...)
		for i := 0; i < padding; i++ {
			data[n+i] = byte(padding)
		}
	}
	encryptedKey := make([]byte, len(data))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encryptedKey, data)
	var algorithmIdentifier pkix.AlgorithmIdentifier
	algorithmIdentifier.Algorithm = oidKEYSHA1
	algorithmIdentifier.Parameters.Tag = 5
	algorithmIdentifier.Parameters.IsCompound = false
	algorithmIdentifier.Parameters.FullBytes = []byte{5, 0}
	keyDerivationFunc := Pbes2KDfs{
		oidPBKDF2,
		Pkdf2Params{
			salt,
			iter,
			algorithmIdentifier,
		},
	}
	encryptionScheme := Pbes2Encs{
		oidAES256CBC,
		iv,
	}
	pbes2Algorithms := Pbes2Algorithms{
		oidPBES2,
		Pbes2Params{
			keyDerivationFunc,
			encryptionScheme,
		},
	}
	encryptedPkey := EncryptedPrivateKeyInfo{
		pbes2Algorithms,
		encryptedKey,
	}

	encryptedBytes, err := asn1.Marshal(encryptedPkey)
	if err != nil {
		return nil, err
	}
	return &pem.Block{
		Type: blockType,
		Headers: map[string]string{
			"Proc-Type": "4,ENCRYPTED",
			"DEK-Info":  ciph.name + "," + hex.EncodeToString(iv),
		},
		Bytes: encryptedBytes,
	}, nil
}

func cipherByName(name string) *rfc1423Algo {
	for i := range rfc1423Algos {
		alg := &rfc1423Algos[i]
		if alg.name == name {
			return alg
		}
	}
	return nil
}

func cipherByKey(key PEMCipher) *rfc1423Algo {
	for i := range rfc1423Algos {
		alg := &rfc1423Algos[i]
		if alg.cipher == key {
			return alg
		}
	}
	return nil
}

// copy from crypto/pbkdf2.go
func pbkdf(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}
