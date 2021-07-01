package codec

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"io/ioutil"
)

// 设置私钥
func NewRSAPrivateKey(privateKey []byte) (r *RSAPrivateKey, err error) {
	r = &RSAPrivateKey{}
	if privateKey != nil {
		err = r.SetPrivateKeyBytes(privateKey)
	}
	return
}

type RSAPrivateKey struct {
	keyBytes   []byte          //私钥内容
	privateKey *rsa.PrivateKey //私钥
}

func (r *RSAPrivateKey) SetPrivateKeyBytes(privateKey []byte) error {
	r.keyBytes = privateKey
	_, err := r.GetPrivateKey()
	return err
}

func (r *RSAPrivateKey) SetPrivateKey(privateKey *rsa.PrivateKey) *RSAPrivateKey {
	r.privateKey = privateKey
	return r
}

// *rsa.PrivateKey
func (r *RSAPrivateKey) GetPrivateKey() (*rsa.PrivateKey, error) {
	var err error
	if r.privateKey == nil {
		r.privateKey, err = getPrivateKey(r.keyBytes)
	}
	return r.privateKey, err
}

// 私钥加密
func (rsas *RSAPrivateKey) Encrypt(input []byte) ([]byte, error) {
	if rsas.privateKey == nil {
		return nil, ErrPrivateKeyNotSet
	}
	output := bytes.NewBuffer(nil)
	err := privateKeyIO(rsas.privateKey, bytes.NewReader(input), output, true)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(output)
}

// 私钥解密
func (r *RSAPrivateKey) Decrypt(input []byte) ([]byte, error) {
	if r.privateKey == nil {
		return nil, ErrPrivateKeyNotSet
	}
	output := bytes.NewBuffer(nil)
	err := privateKeyIO(r.privateKey, bytes.NewReader(input), output, false)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(output)
}

/**
 * 使用RSAWithMD5算法签名
 */
func (r *RSAPrivateKey) SignMd5(data []byte) ([]byte, error) {
	md5Hash := md5.New()
	md5Hash.Write(data)
	hashed := md5Hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.MD5, hashed)
}

/**
 * 使用RSAWithSHA1算法签名
 */
func (r *RSAPrivateKey) SignSha1(data []byte) ([]byte, error) {
	sha1Hash := sha1.New()
	sha1Hash.Write(data)
	hashed := sha1Hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA1, hashed)
}

/**
 * 使用RSAWithSHA256算法签名
 */
func (r *RSAPrivateKey) SignSha256(data []byte) ([]byte, error) {
	sha256Hash := sha256.New()
	sha256Hash.Write(data)
	hashed := sha256Hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA256, hashed)
}
