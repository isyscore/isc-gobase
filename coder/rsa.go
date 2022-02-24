package coder

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	RSA_KEY_SIZE_256  = 256
	RSA_KEY_SIZE_512  = 512
	RSA_KEY_SIZE_1024 = 1024
	RSA_KEY_SIZE_2048 = 2048
)

func RSAGenerateKeyPair(size int, privateKeyPath string, publicKeyPath string) error {
	privKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return err
	}
	xpriv := x509.MarshalPKCS1PrivateKey(privKey)
	privFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	defer func(privFile *os.File) { _ = privFile.Close() }(privFile)
	privBlock := pem.Block{
		Type:  "RSA Private Key",
		Bytes: xpriv,
	}
	err = pem.Encode(privFile, &privBlock)
	if err != nil {
		return err
	}

	pubKey := privKey.PublicKey
	xpub, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		return err
	}
	pubFile, err := os.Create(publicKeyPath)
	if err != nil {
		return err
	}
	defer func(pubFile *os.File) { _ = pubFile.Close() }(pubFile)
	pubBlock := pem.Block{
		Type:  "RSA Public Key",
		Bytes: xpub,
	}
	err = pem.Encode(pubFile, &pubBlock)
	if err != nil {
		return err
	}
	return nil
}

func RSAEncrypt(content string, publicKeyPath string) (string, error) {
	file, err := os.Open(publicKeyPath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) { _ = file.Close() }(file)
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(buf)
	pubKeyIntf, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pubKey := pubKeyIntf.(*rsa.PublicKey)
	text, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(content))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", text), nil
}

func RSADecrypt(content string, privateKeyPath string) (string, error) {
	file, err := os.Open(privateKeyPath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) { _ = file.Close() }(file)
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(buf)
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	b, err := hex.DecodeString(content)
	if err != nil {
		return "", err
	}
	text, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, b)
	if err != nil {
		return "", err
	}
	return string(text), nil
}
