package compress

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io/ioutil"
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
	d, err := ioutil.ReadAll(gz)
	if strings.Contains(err.Error(), "unexpected EOF") {
		// 解gzip时，读到最后必定是unexpected EOF，这里做特殊处理
		err = nil
	}
	return d, err
}
