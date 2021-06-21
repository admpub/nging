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
func NewRSAPrivateKey(priStr string) (r *RSAPrivateKey, err error) {
	r = &RSAPrivateKey{}
	r.priStr = priStr
	r.prikey, err = r.GetPrivatekey()
	return
}

type RSAPrivateKey struct {
	priStr string          //私钥字符串
	prikey *rsa.PrivateKey //私钥
}

// *rsa.PublicKey
func (r *RSAPrivateKey) GetPrivatekey() (*rsa.PrivateKey, error) {
	return getPriKey([]byte(r.priStr))
}

// 私钥加密
func (rsas *RSAPrivateKey) Encrypt(input []byte) ([]byte, error) {
	if rsas.prikey == nil {
		return nil, ErrPrivateKeyNotSet
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(rsas.prikey, bytes.NewReader(input), output, true)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(output)
}

// 私钥解密
func (r *RSAPrivateKey) Decrypt(input []byte) ([]byte, error) {
	if r.prikey == nil {
		return nil, ErrPrivateKeyNotSet
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(r.prikey, bytes.NewReader(input), output, false)
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
	return rsa.SignPKCS1v15(rand.Reader, r.prikey, crypto.MD5, hashed)
}

/**
 * 使用RSAWithSHA1算法签名
 */
func (r *RSAPrivateKey) SignSha1(data []byte) ([]byte, error) {
	sha1Hash := sha1.New()
	sha1Hash.Write(data)
	hashed := sha1Hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.prikey, crypto.SHA1, hashed)
}

/**
 * 使用RSAWithSHA256算法签名
 */
func (r *RSAPrivateKey) SignSha256(data []byte) ([]byte, error) {
	sha256Hash := sha256.New()
	s_data := []byte(data)
	sha256Hash.Write(s_data)
	hashed := sha256Hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.prikey, crypto.SHA256, hashed)
}
