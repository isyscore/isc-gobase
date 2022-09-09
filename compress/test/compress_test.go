package test

import (
	"github.com/isyscore/isc-gobase/compress"
	"github.com/isyscore/isc-gobase/file"
	"github.com/isyscore/isc-gobase/isc"
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

func TestGZipFile(t *testing.T) {
	file.AppendFile("./zip/test3.txt", "cc")
	compress.GzipCompressFile("./zip/test3.txt", "./zip/test3.txt.gz")
	file.DeleteDirs("./zip/")
}

func TestUnGZipFile(t *testing.T) {
	file.AppendFile("./zip/test3.txt", "cc")
	compress.GzipCompressFile("./zip/test3.txt", "./zip/test3.txt.gz")
	compress.GzipDeCompressFile("./zip/test3.txt.gz", "./zip/test3_1.txt")
	file.DeleteDirs("./zip/")
}

func TestZip(t *testing.T) {
	f1 := "./zip/test1.txt"
	f2 := "./zip/test2.txt"
	f3 := "./zip/test3.txt"
	var files = []string{f1, f2, f3}
	dest := "./zip/test.zip"
	srcSize := isc.FormatSize(file.SizeList(files))
	if err := compress.Zip(dest, files); err == nil {
		dstSize := isc.FormatSize(file.Size("./zip/test.zip"))
		// Compress success, size: 714.55KB -> 5.88KB
		t.Logf("Compress success, zip size: %s -> %s", srcSize, dstSize)
	} else {
		t.Logf("Compress err: %v", err)
	}
	file.DeleteFile("./zip/test.zip")
}

func TestUnzip(t *testing.T) {
	// zip
	f1 := "./zip/test1.txt"
	f2 := "./zip/test2.txt"
	f3 := "./zip/test3.txt"
	var files = []string{f1, f2, f3}

	zipFile := "./zip/test.zip"
	compress.Zip(zipFile, files)

	// unzip
	dest := "./zip/uncomp"
	err := compress.Unzip(zipFile, dest)
	if err == nil {
		srcSize := isc.FormatSize(file.Size("./zip/test.zip"))
		dstSize := isc.FormatSize(file.SizeList(files))

		// Decompress success, unzip size: 5.88KB -> 714.55KB
		t.Logf("Decompress success, unzip size: %s -> %s", srcSize, dstSize)
	} else {
		t.Logf("Decompress err: %v", err)
	}
	file.DeleteFile("./zip/test.zip")
	file.DeleteDirs("./zip/uncomp")
}
