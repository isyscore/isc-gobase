package encoding

import "net/url"

func UrlEncoding(str string, charset string) (string, error) {
	if s0, err := Convert(str, "UTF-8", charset); err == nil {
		return url.QueryEscape(s0), nil
	} else {
		return "", err
	}
}

func UrlDecoding(str string, charset string) (string, error) {
	if s0, err := url.QueryUnescape(str); err != nil {
		return "", err
	} else {
		return Convert(s0, charset, "UTF-8")
	}
}
