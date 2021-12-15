package dgoogauth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"
)

var (
	Size           = "400x400" // dimensions for the QrCode()
	Issuer         = "admpub"
	QrApiUrl       = "https://chart.googleapis.com/chart?chs=%s&cht=qr&choe=UTF-8&chl=%s"
	Length   uint8 = 10 //length of the generated tokens
	LifeTime int64 = 30 //seconds
)

type KeyData struct {
	Original string `json:"original"`
	Encoded  string `json:"encoded"`
}

func (d *KeyData) OTP(account string) string {
	return OTPData(account, d.Encoded)
}

func (d *KeyData) Size() string {
	return Size
}

func OTPData(account, base32encode string) string {
	otpStr := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", Issuer, account, base32encode, Issuer)
	return otpStr
}

func QrCode(account, base32encode, qrApiURL string) string {
	if len(qrApiURL) == 0 {
		qrApiURL = QrApiUrl
	}
	return fmt.Sprintf(qrApiURL, Size, url.QueryEscape(OTPData(account, base32encode)))
}

//GenKeyData 生成KeyData
func GenKeyData() *KeyData {
	ts := uint64(time.Now().Unix() / LifeTime)
	text := counterToBytes(ts)
	secret := randomString(100)
	hash := hmacSHA1([]byte(secret), text)
	binary := truncate(hash)
	otp := int64(binary) % int64(math.Pow10(int(Length)))
	key := fmt.Sprintf(fmt.Sprintf("%%0%dd", Length), otp)

	encode := base32.StdEncoding.EncodeToString([]byte(key))
	encode = strings.TrimRight(encode, "=")
	return &KeyData{key, encode}
}

//GenQrCode 生成QrCode
func GenQrCode(account string, qrApiURL string) (*KeyData, string) {
	keyData := GenKeyData()
	return keyData, keyData.OTP(account)
}

func VerifyFrom(keyData *KeyData, password string) (bool, error) {
	return Verify(keyData.Encoded, password)
}

func NewConfig(secret string) *OTPConfig {
	return &OTPConfig{WindowSize: 2, Secret: secret}
}

//Verify 验证
func Verify(secret string, password string) (bool, error) {
	conf := NewConfig(secret)
	return conf.Authenticate(password)
}

func randomString(size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func counterToBytes(counter uint64) (text []byte) {
	text = make([]byte, 8)
	for i := 7; i >= 0; i-- {
		text[i] = byte(counter & 0xff)
		counter = counter >> 8
	}
	return
}

func hmacSHA1(key, text []byte) []byte {
	H := hmac.New(sha1.New, key)
	H.Write([]byte(text))
	return H.Sum(nil)
}

func truncate(hash []byte) int {
	offset := int(hash[len(hash)-1] & 0xf)
	return ((int(hash[offset]) & 0x7f) << 24) |
		((int(hash[offset+1] & 0xff)) << 16) |
		((int(hash[offset+2] & 0xff)) << 8) |
		(int(hash[offset+3]) & 0xff)
}
