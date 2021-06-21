package codec

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
)

// 设置公钥
func NewRSAPublicKey(pubStr string) (r *RSAPublicKey, err error) {
	r = &RSAPublicKey{}
	r.pubStr = pubStr
	r.pubkey, err = r.GetPublickey()
	return
}

type RSAPublicKey struct {
	pubStr string         //公钥字符串
	pubkey *rsa.PublicKey //公钥
}

// *rsa.PrivateKey
func (r *RSAPublicKey) GetPublickey() (*rsa.PublicKey, error) {
	return getPubKey([]byte(r.pubStr))
}

// 公钥加密
func (r *RSAPublicKey) Encrypt(input []byte) ([]byte, error) {
	if r.pubkey == nil {
		return nil, ErrPublicKeyNotSet
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(r.pubkey, bytes.NewReader(input), output, true)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(output)
}

// 公钥解密
func (r *RSAPublicKey) Decrypt(input []byte) ([]byte, error) {
	if r.pubkey == nil {
		return nil, ErrPublicKeyNotSet
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(r.pubkey, bytes.NewReader(input), output, false)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(output)
}

/**
 * 使用RSAWithMD5验证签名
 */
func (r *RSAPublicKey) VerifySignMd5(data string, signData string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	hash := md5.New()
	hash.Write([]byte(data))
	return rsa.VerifyPKCS1v15(r.pubkey, crypto.MD5, hash.Sum(nil), sign)
}

/**
 * 使用RSAWithSHA1验证签名
 */
func (r *RSAPublicKey) VerifySignSha1(data []byte, sign []byte) error {
	hash := sha1.New()
	hash.Write(data)
	return rsa.VerifyPKCS1v15(r.pubkey, crypto.SHA1, hash.Sum(nil), sign)
}

/**
 * 使用RSAWithSHA256验证签名
 */
func (r *RSAPublicKey) VerifySignSha256(data []byte, sign []byte) error {
	hash := sha256.New()
	hash.Write(data)

	return rsa.VerifyPKCS1v15(r.pubkey, crypto.SHA256, hash.Sum(nil), sign)
}
