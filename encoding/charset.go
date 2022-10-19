package encoding

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
	"io"
	"log"
)

const (
	UTF8        = "UTF-8"
	UTF16       = "UTF-16"
	UTF16LE     = "UTF-16LE"
	UTF16BE     = "UTF-16BE"
	GBK         = "gbk"
	GB2312      = "GB2312"
	BIG5        = "Big5"
	GB18030     = "gb18030"
	EUCJP       = "EUC-JP"
	ISO2022JP   = "ISO-2022-JP"
	SHIFTJIS    = "Shift_JIS"
	EUCKR       = "EUC-KR"
	ISO8859_2   = "ISO-8859-2"
	ISO8859_3   = "ISO-8859-3"
	ISO8859_4   = "ISO-8859-4"
	ISO8859_5   = "ISO-8859-5"
	ISO8859_7   = "ISO-8859-7"
	ISO8859_9   = "ISO-8859-9"
	ISO8859_10  = "ISO-8859-10"
	ISO8859_13  = "ISO-8859-13"
	ISO8859_14  = "ISO-8859-14"
	ISO8859_15  = "ISO-8859-15"
	ISO8859_16  = "ISO-8859-16"
	WINDOWS1250 = "windows-1250"
	WINDOWS1251 = "windows-1251"
	WINDOWS1252 = "windows-1252"
	WINDOWS1253 = "windows-1253"
	WINDOWS1254 = "windows-1254"
	WINDOWS1255 = "windows-1255"
	WINDOWS1256 = "windows-1256"
	WINDOWS1257 = "windows-1257"
	WINDOWS1258 = "windows-1258"
	WINDOWS874  = "windows-874"
	MACINTOSH   = "macintosh"
	KOI8R       = "KOI8-R"
	KOI8U       = "KOI8-U"
)

// 别名
var charsetAlias = map[string]string{"HZGB2312": "HZ-GB-2312", "hzgb2312": "HZ-GB-2312", "GB2312": "HZ-GB-2312", "gb2312": "HZ-GB-2312"}

// Supported 判断指定的编码是否被支持
func Supported(charset string) bool {
	return getEncoding(charset) != nil
}

func Convert(src string, srcCharset string, dstCharset string) (string, error) {
	if dstCharset == srcCharset {
		return src, nil
	}
	dst := src
	if srcCharset != "UTF-8" {
		if e := getEncoding(srcCharset); e != nil {
			tmp, err := io.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewDecoder()),
			)
			if err != nil {
				return "", fmt.Errorf(`convert string "%s" to utf8 failed`, srcCharset)
			}
			src = string(tmp)
		} else {
			return dst, fmt.Errorf(`unsupported charset "%s"`, srcCharset)
		}
	}
	if dstCharset != "UTF-8" {
		if e := getEncoding(dstCharset); e != nil {
			tmp, err := io.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewEncoder()),
			)
			if err != nil {
				return "", fmt.Errorf(`convert string from utf8 to "%s" failed`, dstCharset)
			}
			dst = string(tmp)
		} else {
			return dst, fmt.Errorf(`unsupported charset "%s"`, dstCharset)
		}
	} else {
		dst = src
	}
	return dst, nil
}

func StringToUTF8(src string, srcCharset string) (string, error) {
	return Convert(src, srcCharset, "UTF-8")
}

func UTF8ToString(src string, dstCharset string) (string, error) {
	return Convert(src, "UTF-8", dstCharset)
}

func getEncoding(charset string) encoding.Encoding {
	if c, ok := charsetAlias[charset]; ok {
		charset = c
	}
	enc, err := ianaindex.MIB.Encoding(charset)
	if err != nil {
		log.Printf("[WARN] charset %s not supported", charset)
	}
	return enc
}
