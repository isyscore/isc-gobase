package coder

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rc4"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

func AesDecrypt(content string, key string, iv string) string {
	b, _ := base64.StdEncoding.DecodeString(content)
	block, _ := aes.NewCipher([]byte(key))
	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	originData := make([]byte, len(b))
	mode.CryptBlocks(originData, b)
	origData := pkcs5UnPadding(originData)
	return origData
}

func AesEncrypt(content string, key string, iv string) string {
	origData := pkcs5Padding([]byte(content), aes.BlockSize)
	block, _ := aes.NewCipher([]byte(key))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(origData))
	mode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted)
}

// AesDecryptECB 兼容java的AES解密方式
func AesDecryptECB(content string, key string) string {
	b, _ := base64.StdEncoding.DecodeString(content)
	cp, _ := aes.NewCipher([]byte(key))
	d := make([]byte, len(b))
	size := 16
	for bs, be := 0, size; bs < len(b); bs, be = bs+size, be+size {
		cp.Decrypt(d[bs:be], b[bs:be])
	}
	return strings.TrimSpace(string(d))
}

// AesEncryptECB 兼容java的AES加密方式
func AesEncryptECB(content string, key string) string {
	b := padding([]byte(content))
	cp, _ := aes.NewCipher([]byte(key))
	d := make([]byte, len(b))
	size := 16
	for bs, be := 0, size; bs < len(b); bs, be = bs+size, be+size {
		cp.Encrypt(d[bs:be], b[bs:be])
	}
	return base64.StdEncoding.EncodeToString(d)
}

func DESEncryptCBC(content string, key string, iv string) string {
	block, _ := des.NewCipher([]byte(key))
	data := pkcs5Padding([]byte(content), block.BlockSize())
	dest := make([]byte, len(data))
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	blockMode.CryptBlocks(dest, data)
	return fmt.Sprintf("%x", dest)
}

func DESDecryptCBC(content string, key string, iv string) string {
	b, _ := hex.DecodeString(content)
	block, _ := des.NewCipher([]byte(key))
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	originData := make([]byte, len(b))
	blockMode.CryptBlocks(originData, b)
	origData := pkcs5UnPadding(originData)
	return origData
}

func DESEncryptECB(content string, key string) string {
	block, _ := des.NewCipher([]byte(key))
	size := block.BlockSize()
	data := pkcs5Padding([]byte(content), size)
	if len(data)%size != 0 {
		return ""
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:size])
		data = data[size:]
		dst = dst[size:]
	}
	return fmt.Sprintf("%x", out)
}

func DESDecryptECB(content string, key string) string {
	b, _ := hex.DecodeString(content)
	block, _ := des.NewCipher([]byte(key))
	size := block.BlockSize()
	out := make([]byte, len(b))
	dst := out
	for len(b) > 0 {
		block.Decrypt(dst, b[:size])
		b = b[size:]
		dst = dst[size:]
	}
	return pkcs5UnPadding(out)
}

func RC4Encrypt(content string, key string) string {
	dest := make([]byte, len(content))
	cp, _ := rc4.NewCipher([]byte(key))
	cp.XORKeyStream(dest, []byte(content))
	return fmt.Sprintf("%x", dest)
}

func RC4Decrypt(content string, key string) string {
	b, _ := hex.DecodeString(content)
	dest := make([]byte, len(b))
	cp, _ := rc4.NewCipher([]byte(key))
	cp.XORKeyStream(dest, b)
	return string(dest)
}

func Base64Encrypt(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}

func Base64Decrypt(content string) string {
	b, _ := base64.StdEncoding.DecodeString(content)
	return string(b)
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) string {
	length := len(origData)
	unpadding := int(origData[length-1])
	return string(origData[:(length - unpadding)])
}

func padding(src []byte) []byte {
	paddingCount := aes.BlockSize - len(src)%aes.BlockSize
	if paddingCount == 0 {
		return src
	} else {
		return append(src, bytes.Repeat([]byte{byte(0)}, paddingCount)...)
	}
}
