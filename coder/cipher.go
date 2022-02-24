package coder

import (
	"bytes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

type Crypto interface {
	Encrypt(plainText []byte) (string, error)
	Decrypt(cipherText string) (string, error)
}

type Cipher struct {
	GroupMode  int
	FillMode   FillMode
	DecodeType int
	Key        []byte
	Iv         []byte
	Output     CipherText
}

func (c *Cipher) Encrypt(block cipher.Block, plainData []byte) (err error) {
	c.Output = make([]byte, len(plainData))
	if c.GroupMode == CBCMode {
		cipher.NewCBCEncrypter(block, c.Iv).CryptBlocks(c.Output, plainData)
		return
	}
	if c.GroupMode == ECBMode {
		c.NewECBEncrypter(block, plainData)
		return
	}
	return
}

func (c *Cipher) Decrypt(block cipher.Block, cipherData []byte) (err error) {
	c.Output = make([]byte, len(cipherData))
	if c.GroupMode == CBCMode {
		cipher.NewCBCDecrypter(block, c.Iv).CryptBlocks(c.Output, cipherData)
		return
	}
	if c.GroupMode == ECBMode {
		c.NewECBDecrypter(block, cipherData)
		return
	}
	return
}

// Encode default print format is base64
func (c *Cipher) Encode() string {
	if c.DecodeType == PrintHex {
		return c.Output.hexEncode()
	} else {
		return c.Output.base64Encode()
	}
}

func (c *Cipher) Decode(cipherText string) ([]byte, error) {
	if c.DecodeType == PrintBase64 {
		return base64Decode(cipherText)
	} else if c.DecodeType == PrintHex {
		return hexDecode(cipherText)
	} else {
		return nil, errors.New("unsupported print type")
	}
}

func (c *Cipher) Fill(plainText []byte, blockSize int) []byte {
	if c.FillMode == PkcsZero {
		return c.FillMode.zeroPadding(plainText, blockSize)
	} else {
		return c.FillMode.pkcs7Padding(plainText, blockSize)
	}
}

func (c *Cipher) UnFill(plainText []byte) (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	if c.FillMode == Pkcs7 {
		return c.FillMode.pkcsUnPadding(plainText), nil
	} else if c.FillMode == PkcsZero {
		return c.FillMode.unZeroPadding(plainText), nil
	} else {
		return nil, errors.New("unsupported fill mode")
	}
}

func (c *Cipher) NewECBEncrypter(block cipher.Block, plainData []byte) {
	tempText := c.Output
	for len(plainData) > 0 {
		block.Encrypt(tempText, plainData[:block.BlockSize()])
		plainData = plainData[block.BlockSize():]
		tempText = tempText[block.BlockSize():]
	}
}

func (c *Cipher) NewECBDecrypter(block cipher.Block, cipherData []byte) {
	tempText := c.Output
	for len(cipherData) > 0 {
		block.Decrypt(tempText, cipherData[:block.BlockSize()])
		cipherData = cipherData[block.BlockSize():]
		tempText = tempText[block.BlockSize():]
	}
}

const (
	CBCMode = iota
	CFBMode
	CTRMode
	ECBMode
	OFBMode
)

type FillMode int

const (
	PkcsZero FillMode = iota
	Pkcs7
)

func (fm FillMode) pkcs7Padding(plainText []byte, blockSize int) []byte {
	paddingSize := blockSize - len(plainText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(plainText, paddingText...)
}

func (fm FillMode) pkcsUnPadding(plainText []byte) []byte {
	length := len(plainText)
	number := int(plainText[length-1])
	return plainText[:length-number]
}

func (fm FillMode) zeroPadding(plainText []byte, blockSize int) []byte {
	if plainText[len(plainText)-1] == 0 {
		return nil
	}
	paddingSize := blockSize - len(plainText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(0)}, paddingSize)
	return append(plainText, paddingText...)
}

func (fm FillMode) unZeroPadding(plainText []byte) []byte {
	length := len(plainText)
	count := 1
	for i := length - 1; i > 0; i-- {
		if plainText[i] == 0 && plainText[i-1] == plainText[i] {
			count++
		}
	}
	return plainText[:length-count]
}

type CipherText []byte

const (
	PrintHex = iota
	PrintBase64
)

func (ct CipherText) hexEncode() string {
	return hex.EncodeToString(ct)
}

func (ct CipherText) base64Encode() string {
	return base64.StdEncoding.EncodeToString(ct)
}

func hexDecode(cipherText string) ([]byte, error) {
	return hex.DecodeString(cipherText)
}

func base64Decode(cipherText string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(cipherText)
}
