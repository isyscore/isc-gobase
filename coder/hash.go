package coder

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func MD5String(s string) string {
	b := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", b)
}

func MD5File(filePath string) (string, error) {
	if file, err := os.Open(filePath); err == nil {
		m := md5.New()
		_, _ = io.Copy(m, file)
		return fmt.Sprintf("%x", m.Sum(nil)), nil
	} else {
		return "", err
	}
}

func Sha1String(s string) string {
	b := sha1.Sum([]byte(s))
	return fmt.Sprintf("%x", b)
}

func Sha1File(filePath string) (string, error) {
	if file, err := os.Open(filePath); err == nil {
		s := sha1.New()
		_, _ = io.Copy(s, file)
		return fmt.Sprintf("%x", s.Sum(nil)), nil
	} else {
		return "", err
	}

}

func Sha256String(s string) string {
	b := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", b)
}

func Sha256File(filePath string) (string, error) {
	if file, err := os.Open(filePath); err == nil {
		s := sha256.New()
		_, _ = io.Copy(s, file)
		return fmt.Sprintf("%x", s.Sum(nil)), nil
	} else {
		return "", err
	}
}

func HMacMD5String(s string, key string) string {
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HMacMD5File(filePath string, key string) (string, error) {
	if file, err := os.Open(filePath); err == nil {
		h := hmac.New(md5.New, []byte(key))
		_, _ = io.Copy(h, file)
		return fmt.Sprintf("%x", h.Sum(nil)), nil
	} else {
		return "", err
	}
}

func HMacSha1String(s string, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HMacSha1File(filePath string, key string) (string, error) {
	if file, err := os.Open(filePath); err == nil {
		h := hmac.New(sha1.New, []byte(key))
		_, _ = io.Copy(h, file)
		return fmt.Sprintf("%x", h.Sum(nil)), nil
	} else {
		return "", err
	}
}

func HMacSha256String(s string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HMacSha256File(filePath string, key string) (string, error) {
	if file, err := os.Open(filePath); err == nil {
		h := hmac.New(sha256.New, []byte(key))
		_, _ = io.Copy(h, file)
		return fmt.Sprintf("%x", h.Sum(nil)), nil
	} else {
		return "", err
	}
}
