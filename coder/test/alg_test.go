package test

import (
	"crypto/dsa"
	"github.com/isyscore/isc-gobase/coder"
	"github.com/isyscore/isc-gobase/file"
	"testing"
)

func TestMd5(t *testing.T) {
	str := "abcdefg"
	s1 := coder.MD5String(str)
	t.Logf("%s md5 is %s", str, s1)

	f := "/Users/rarnu/Code/isyscore/opensource/isc-gobase/go.mod"
	s2, err := coder.MD5File(f)
	t.Logf("%s md5 is %s (%v)", f, s2, err)
}

func TestSha1(t *testing.T) {
	str := "abcdefg"
	s1 := coder.Sha1String(str)
	t.Logf("%s sha1 is %s", str, s1)

	f := "/Users/rarnu/Code/isyscore/opensource/isc-gobase/go.mod"
	s2, err := coder.Sha1File(f)
	t.Logf("%s sha1 is %s (%v)", f, s2, err)
}

func TestHMac(t *testing.T) {
	key := "isyscore"
	str := "abcdefg"
	s1 := coder.HMacMD5String(str, key)
	t.Logf("%s hmac is %s", str, s1)
}

func TestRC4(t *testing.T) {
	key := "isyscore"
	str := "abcdefg"
	s1 := coder.RC4Encrypt(str, key)
	t.Logf("%s rc4 is %s", str, s1)

	s2 := coder.RC4Decrypt(s1, key)
	t.Logf("%s rc4 is %s", s1, s2)
}

func TestDES(t *testing.T) {
	// CBC
	key := "isyscore"
	iv := "12345678"
	str := "abcdefg"
	s1 := coder.DESEncryptCBC(str, key, iv)
	t.Logf("%s des is %s", str, s1)
	s2 := coder.DESDecryptCBC(s1, key, iv)
	t.Logf("%s des is %s", s1, s2)

	// ECB
	ss1 := coder.DESEncryptECB(str, key)
	t.Logf("%s des is %s", str, ss1)
	ss2 := coder.DESDecryptECB(ss1, key)
	t.Logf("%s des is %s", ss1, ss2)
}

func TestRSA(t *testing.T) {
	privKeyPath := "/Users/rarnu/Code/isyscore/opensource/isc-gobase/test/rsa/private.pem"
	pubKeyPath := "/Users/rarnu/Code/isyscore/opensource/isc-gobase/test/rsa/public.pem"
	if !file.FileExists(privKeyPath) {
		err := coder.RSAGenerateKeyPair(coder.RSA_KEY_SIZE_1024, privKeyPath, pubKeyPath)
		if err != nil {
			t.Logf("generate key pair error: %v", err)
			return
		}
	}

	str := "abcdefg"
	s1, err := coder.RSAEncrypt(str, pubKeyPath)
	if err != nil {
		t.Logf("encrypt error: %v", err)
		return
	}
	t.Logf("%s rsa is %s", str, s1)

	s2, err := coder.RSADecrypt(s1, privKeyPath)
	t.Logf("%s rsa is %s (%v)", s1, s2, err)
}

func TestDSA(t *testing.T) {
	privKeyPath := "/Users/rarnu/Code/isyscore/opensource/isc-gobase/test/dsa/private.pem"
	pubKeyPath := "/Users/rarnu/Code/isyscore/opensource/isc-gobase/test/dsa/public.pem"
	if !file.FileExists(privKeyPath) {
		err := coder.DSAGenerateKeyPair(dsa.L1024N160, privKeyPath, pubKeyPath)
		if err != nil {
			t.Logf("generate key pair error: %v", err)
			return
		}
	}

	str := "abcdefg"
	r, s, err := coder.DSASign(str, privKeyPath)
	if err != nil {
		t.Logf("sign error: %v", err)
		return
	}
	verify, err := coder.DSAVerify(str, pubKeyPath, r, s)
	if err != nil {
		t.Logf("verify error: %v", err)
		return
	}
	t.Logf("%s dsa is %v", str, verify)
}

func TestAes(t *testing.T) {
	// encrypt
	content := "abcdefg"
	// ISC-PKI-BRIDGE~1

	// key的长度必须是 16，24，32
	key := "isyscore12345678"
	iv := "0102030405060708"
	s1 := coder.AesEncrypt(content, key, iv)
	t.Logf("%s aes is %s", content, s1)

	// decrypt
	s2 := coder.AesDecrypt(s1, key, iv)
	t.Logf("%s aes is %s", s1, s2)
}

func TestAesJava(t *testing.T) {
	content := "abcdefg"
	key := "isyscore12345678"
	s1 := coder.AesEncryptECB(content, key)
	t.Logf("%s aes is %s", content, s1)
	s2 := coder.AesDecryptECB(s1, key)
	t.Logf("%s aes is %s", s1, s2)
}

func TestAesCBC(t *testing.T) {
	content := "abcdefg"
	key := "isyscore12345678"
	s1 := coder.AesEncryptCBC(content, key)
	t.Logf("%s aes is %s", content, s1)
	s2 := coder.AesDecryptCBC(s1, key)
	t.Logf("%s aes is %s", s1, s2)
}

func Test2333(t *testing.T) {
	content := "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"
	key := "ISC-PKI-BRIDGE~1"
	s1 := coder.AesDecryptECB(content, key)
	t.Logf(s1)
}
