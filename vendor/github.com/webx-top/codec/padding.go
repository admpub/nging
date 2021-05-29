package codec

import "bytes"

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length > unpadding {
		return origData[:(length - unpadding)]
	}
	return []byte{}
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length > unpadding {
		return origData[:(length - unpadding)]
	}
	return []byte{}
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	return PKCS5Padding(ciphertext, blockSize)
}

func PKCS7UnPadding(origData []byte) []byte {
	return PKCS5UnPadding(origData)
}
