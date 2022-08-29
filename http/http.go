package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var httpClient = createHTTPClient()

const (
	MaxIdleConns          int    = 100
	MaxIdleConnsPerHost   int    = 100
	IdleConnTimeout       int    = 90
	ContentTypeJson       string = "application/json; charset=utf-8"
	ContentTypeHtml       string = "text/html; charset=utf-8"
	ContentTypeText       string = "text/plain; charset=utf-8"
	ContentTypeCss        string = "text/css; charset=utf-8"
	ContentTypeJavaScript string = "application/x-javascript; charset=utf-8"
	ContentTypeJpeg       string = "image/jpeg"
	ContentTypePng        string = "image/png"
	ContentTypeGif        string = "image/gif"
	ContentTypeAll        string = "*/*"
	ContentPostForm       string = "application/x-www-form-urlencoded"
)

type NetError struct {
	ErrMsg string
}

func (error *NetError) Error() string {
	return error.ErrMsg
}

type DataResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
		},

		Timeout: 20 * time.Second,
	}
	return client
}

func SetHttpClient(httpClientOuter *http.Client) {
	httpClient = httpClientOuter
}

func GetClient() *http.Client {
	return httpClient
}

// ------------------ get ------------------

func GetSimple(url string) (int, http.Header, any, error) {
	return Get(url, nil, nil)
}

func GetSimpleOfStandard(url string) (int, http.Header, any, error) {
	return GetOfStandard(url, nil, nil)
}

func Get(url string, header http.Header, parameterMap map[string]string) (int, http.Header, any, error) {
	httpRequest, err := http.NewRequest("GET", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return call(httpRequest, url)
}

func GetOfStandard(url string, header http.Header, parameterMap map[string]string) (int, http.Header, any, error) {
	httpRequest, err := http.NewRequest("GET", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return callToStandard(httpRequest, url)
}

// ------------------ head ------------------

func HeadSimple(url string) error {
	return Head(url, nil, nil)
}

func Head(url string, header http.Header, parameterMap map[string]string) error {
	httpRequest, err := http.NewRequest("GET", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return callIgnoreReturn(httpRequest, url)
}

// ------------------ post ------------------

func PostSimple(url string, body any) (int, http.Header, any, error) {
	return Post(url, nil, nil, body)
}

func PostSimpleOfStandard(url string, body any) (int, http.Header, any, error) {
	return PostOfStandard(url, nil, nil, body)
}

func Post(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, any, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("POST", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return call(httpRequest, url)
}

func PostOfStandard(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, any, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("POST", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return callToStandard(httpRequest, url)
}

func PostForm(url string, header http.Header, parameterMap map[string]string) (int, http.Header, any, error) {
	var httpRequest http.Request
	_ = httpRequest.ParseForm()
	if parameterMap != nil {
		_ = httpRequest.ParseForm()
		for k, v := range parameterMap {
			httpRequest.Form.Add(k, v)
		}
	}
	if header != nil {
		httpRequest.Header = header
	}
	body := strings.NewReader(httpRequest.Form.Encode())
	resp, err := httpClient.Post(url, ContentPostForm, body)
	if err != nil {
		return -1, nil, nil, err
	}

	code := resp.StatusCode
	head := resp.Header

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return code, head, nil, err
	}

	return code, head, b, nil
}

// ------------------ put ------------------

func PutSimple(url string, body any) (int, http.Header, any, error) {
	return Put(url, nil, nil, body)
}

func PutSimpleOfStandard(url string, body any) (int, http.Header, any, error) {
	return PutOfStandard(url, nil, nil, body)
}

func Put(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, any, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PUT", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return call(httpRequest, url)
}

func PutOfStandard(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, any, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PUT", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return callToStandard(httpRequest, url)
}

// ------------------ delete ------------------

func DeleteSimple(url string) (int, http.Header, any, error) {
	return Get(url, nil, nil)
}

func DeleteSimpleOfStandard(url string) (int, http.Header, any, error) {
	return GetOfStandard(url, nil, nil)
}

func Delete(url string, header http.Header, parameterMap map[string]string) (int, http.Header, any, error) {
	httpRequest, err := http.NewRequest("DELETE", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return call(httpRequest, url)
}

func DeleteOfStandard(url string, header http.Header, parameterMap map[string]string) (int, http.Header, any, error) {
	httpRequest, err := http.NewRequest("DELETE", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return callToStandard(httpRequest, url)
}

// ------------------ patch ------------------

func PatchSimple(url string, body any) (int, http.Header, any, error) {
	return Post(url, nil, nil, body)
}

func PatchSimpleOfStandard(url string, body any) (int, http.Header, any, error) {
	return PostOfStandard(url, nil, nil, body)
}

func Patch(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, any, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PATCH", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return call(httpRequest, url)
}

func PatchOfStandard(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, any, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PATCH", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return callToStandard(httpRequest, url)
}

func call(httpRequest *http.Request, url string) (int, http.Header, any, error) {
	if httpResponse, err := httpClient.Do(httpRequest); err != nil && httpResponse == nil {
		log.Printf("Error sending request to API endpoint. %+v", err)
		return -1, nil, nil, &NetError{ErrMsg: "Error sending request, url: " + url + ", err" + err.Error()}
	} else {
		if httpResponse == nil {
			log.Printf("httpResponse is nil\n")
			return -1, nil, nil, nil
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Body close error(%v)", err)
			}
		}(httpResponse.Body)

		code := httpResponse.StatusCode
		headers := httpResponse.Header
		if code != http.StatusOK {
			body, _ := ioutil.ReadAll(httpResponse.Body)
			return code, headers, nil, &NetError{ErrMsg: "remote error, url: " + url + ", code " + strconv.Itoa(code) + ", message: " + string(body)}
		}

		// We have seen inconsistencies even when we get 200 OK response
		body, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			log.Printf("Couldn't parse response body(%v)", err)
			return code, headers, nil, &NetError{ErrMsg: "Couldn't parse response body, err: " + err.Error()}
		}

		return code, headers, body, nil
	}
}

// ------------------ trace ------------------
// ------------------ options ------------------
// 暂时先不处理

func callIgnoreReturn(httpRequest *http.Request, url string) error {
	if httpResponse, err := httpClient.Do(httpRequest); err != nil && httpResponse == nil {
		log.Printf("Error sending request to API endpoint. %v", err)
		return &NetError{ErrMsg: "Error sending request, url: " + url + ", err" + err.Error()}
	} else {
		if httpResponse == nil {
			log.Printf("httpResponse is nil\n")
			return nil
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Body close error(%v)", err)
			}
		}(httpResponse.Body)

		code := httpResponse.StatusCode
		if code != http.StatusOK {
			body, _ := ioutil.ReadAll(httpResponse.Body)
			return &NetError{ErrMsg: "remote error, url: " + url + ", code " + strconv.Itoa(code) + ", message: " + string(body)}
		}
		return nil
	}
}

func callToStandard(httpRequest *http.Request, url string) (int, http.Header, any, error) {
	return parseStandard(call(httpRequest, url))
}

func parseStandard(statusCode int, headers http.Header, responseResult any, errs error) (int, http.Header, any, error) {
	if errs != nil {
		return statusCode, headers, nil, errs
	}
	var standRsp DataResponse[any]
	err := json.Unmarshal(responseResult.([]byte), &standRsp)
	if err != nil {
		return statusCode, headers, nil, err
	}

	// 判断业务的失败信息
	if standRsp.Code != 0 && standRsp.Code != 200 {
		return statusCode, headers, nil, &NetError{ErrMsg: fmt.Sprintf("remote err, bizCode=%d, message=%s", standRsp.Code, standRsp.Message)}
	}

	return statusCode, headers, standRsp.Data, nil
}

func urlWithParameter(url string, parameterMap map[string]string) string {
	if parameterMap == nil || len(parameterMap) == 0 {
		return url
	}

	url += "?"

	var parameters []string
	for key, value := range parameterMap {
		parameters = append(parameters, key+"="+value)
	}

	return url + strings.Join(parameters, "&")
}
