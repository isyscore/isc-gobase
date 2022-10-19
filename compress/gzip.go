package compress

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

func GzipCompress(data []byte) ([]byte, error) {
	var res bytes.Buffer
	gz := gzip.NewWriter(&res)
	defer func(gz *gzip.Writer) { _ = gz.Close() }(gz)
	_, err := gz.Write(data)
	_ = gz.Flush()
	if err == nil {
		return res.Bytes(), nil
	} else {
		return nil, err
	}
}

func GzipDecompress(data []byte) ([]byte, error) {
	var res bytes.Buffer
	_ = binary.Write(&res, binary.LittleEndian, data)
	fmt.Printf("res: %v\n", res.Bytes())
	gz, _ := gzip.NewReader(&res)
	defer func(gz *gzip.Reader) { _ = gz.Close() }(gz)
	d, err := io.ReadAll(gz)
	if strings.Contains(err.Error(), "unexpected EOF") {
		// 解gzip时，读到最后必定是unexpected EOF，这里做特殊处理
		err = nil
	}
	return d, err
}

// GzipCompressFile 压缩文件Src到Dst
func GzipCompressFile(src string, dst string) error {
	newfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(newfile *os.File) {
		_ = newfile.Close()
	}(newfile)

	file, err := os.Open(src)
	if err != nil {
		return err
	}

	zw := gzip.NewWriter(newfile)

	filestat, err := file.Stat()
	if err != nil {
		return nil
	}

	zw.Name = filestat.Name()
	zw.ModTime = filestat.ModTime()
	_, err = io.Copy(zw, file)
	if err != nil {
		return nil
	}

	_ = zw.Flush()
	if err := zw.Close(); err != nil {
		return nil
	}
	return nil
}

// GzipDeCompressFile 解压文件Src到Dst
func GzipDeCompressFile(src string, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	newfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(newfile *os.File) {
		_ = newfile.Close()
	}(newfile)

	zr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	filestat, err := file.Stat()
	if err != nil {
		return err
	}

	zr.Name = filestat.Name()
	zr.ModTime = filestat.ModTime()
	_, err = io.Copy(newfile, zr)
	if err != nil {
		return err
	}

	if err := zr.Close(); err != nil {
		return err
	}
	return nil
}
