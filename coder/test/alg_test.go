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
