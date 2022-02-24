package coder

import (
	"crypto/dsa"
	"crypto/rand"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"os"
)

func DSAGenerateKeyPair(size dsa.ParameterSizes, privateKeyPath string, publicKeyPath string) error {
	var param dsa.Parameters
	err := dsa.GenerateParameters(&param, rand.Reader, size)
	if err != nil {
		return err
	}
	var privKey dsa.PrivateKey
	privKey.Parameters = param
	err = dsa.GenerateKey(&privKey, rand.Reader)
	if err != nil {
		return err
	}

	b, err := asn1.Marshal(privKey)
	if err != nil {
		return err
	}
	privBlock := pem.Block{
		Type:  "DSA Private Key",
		Bytes: b,
	}
	privFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	err = pem.Encode(privFile, &privBlock)
	if err != nil {
		return err
	}

	pubKey := privKey.PublicKey
	b, err = asn1.Marshal(pubKey)
	pubBlock := pem.Block{
		Type:  "DSA Public Key",
		Bytes: b,
	}
	pubFile, err := os.Create(publicKeyPath)
	if err != nil {
		return err
	}
	err = pem.Encode(pubFile, &pubBlock)
	if err != nil {
		return err
	}
	return nil
}

func DSASign(content string, privKeyPath string) (r, s *big.Int, err error) {
	file, err := os.Open(privKeyPath)
	if err != nil {
		return nil, nil, err
	}
	defer func(file *os.File) { _ = file.Close() }(file)
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return nil, nil, err
	}
	block, _ := pem.Decode(buf)
	var privKey dsa.PrivateKey
	_, err = asn1.Unmarshal(block.Bytes, &privKey)
	if err != nil {
		return nil, nil, err
	}
	rr, ss, err := dsa.Sign(rand.Reader, &privKey, []byte(content))
	if err != nil {
		return nil, nil, err
	}
	return rr, ss, nil
}

func DSAVerify(content string, pubKeyPath string, r, s *big.Int) (bool, error) {
	file, err := os.Open(pubKeyPath)
	if err != nil {
		return false, err
	}
	defer func(file *os.File) { _ = file.Close() }(file)
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return false, err
	}
	block, _ := pem.Decode(buf)
	var pubKey dsa.PublicKey
	_, err = asn1.Unmarshal(block.Bytes, &pubKey)
	if err != nil {
		return false, err
	}
	flag := dsa.Verify(&pubKey, []byte(content), r, s)
	return flag, nil
}
