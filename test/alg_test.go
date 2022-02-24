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

	file := "/Users/rarnu/Code/isyscore/opensource/isc-gobase/go.mod"
	s2, err := coder.MD5File(file)
	t.Logf("%s md5 is %s (%v)", file, s2, err)
}

func TestSha1(t *testing.T) {
	str := "abcdefg"
	s1 := coder.Sha1String(str)
	t.Logf("%s sha1 is %s", str, s1)

	file := "/Users/rarnu/Code/isyscore/opensource/isc-gobase/go.mod"
	s2, err := coder.Sha1File(file)
	t.Logf("%s sha1 is %s (%v)", file, s2, err)
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
