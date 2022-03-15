package test

import (
	"github.com/isyscore/isc-gobase/compress"
	"testing"
)

func TestGZip(t *testing.T) {
	str := []byte("Hello World!")
	if b, err := compress.GzipCompress(str); err == nil {
		t.Logf("%v", b)
	} else {
		t.Logf("GzipCompress err: %v", err)
	}
}

func TestUnGzip(t *testing.T) {
	b := []byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xf2, 0x48, 0xcd, 0xc9, 0xc9, 0x57,
		0x08, 0xcf, 0x2f, 0xca, 0x49, 0x51, 0x04, 0x00,
		0x00, 0x00, 0xff, 0xff}
	if s, err := compress.GzipDecompress(b); err == nil {
		t.Logf("%s", string(s))
	} else {
		t.Logf("GzipDecompress err: %v", err)
	}
}

func TestZip(t *testing.T) {
	f1 := "./zip/test1.txt"
	f2 := "./zip/test2.txt"
	f3 := "./zip/test3.txt"
	var files = []string{f1, f2, f3}
	dest := "./zip/test.zip"
	if err := compress.Compress(dest, "./zip/", files); err == nil {
		t.Logf("Compress success")
	} else {
		t.Logf("Compress err: %v", err)
	}
}

func TestUnzip(t *testing.T) {
	zf := "./zip/test.zip"
	dest := "./zip/uncomp"
	err := compress.Decompress(zf, dest)
	if err == nil {
		t.Logf("Decompress success")
	} else {
		t.Logf("Decompress err: %v", err)
	}
}
