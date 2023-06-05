package http

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "github.com/isyscore/isc-gobase/config"
    "github.com/isyscore/isc-gobase/logger"
    "mime/multipart"

    //"github.com/isyscore/isc-gobase/goid"
    "io"
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

var NetHttpHooks []GobaseHttpHook

func init() {
    NetHttpHooks = []GobaseHttpHook{}
}

type GobaseHttpHook interface {
    Before(ctx context.Context, req *http.Request) (context.Context, http.Header)
    After(ctx context.Context, rsp *http.Response, rspCode int, rspData any, err error)
}

func AddHook(httpHook GobaseHttpHook) {
    NetHttpHooks = append(NetHttpHooks, httpHook)
}

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
    config.LoadConfig()
    client := &http.Client{}

    // 从配置文件中载入配置
    loadClientFromConfig(client)

    return client
}

func loadClientFromConfig(client *http.Client) {
    if config.GetValueString("base.http.timeout") != "" {
        t, err := time.ParseDuration(config.GetValueString("base.http.timeout"))
        if err != nil {
            logger.Warn("读取配置【base.http.timeout】异常", err)
        } else {
            client.Timeout = t
        }
    }

    transport := &http.Transport{}
    if config.GetValueString("base.http.transport.tls-handshake-timeout") != "" {
        t, err := time.ParseDuration(config.GetValueString("base.http.transport.tls-handshake-timeout"))
        if err != nil {
            logger.Warn("读取配置【base.http.transport.tls-handshake-timeout】异常", err)
        } else {
            transport.TLSHandshakeTimeout = t
        }
    }

    if config.GetValueString("base.http.transport.disable-keep-alives") != "" {
        transport.DisableKeepAlives = config.GetValueBool("base.http.transport.disable-keep-alives")
    }

    if config.GetValueString("base.http.transport.disable-compression") != "" {
        transport.DisableCompression = config.GetValueBool("base.http.transport.disable-compression")
    }

    if config.GetValueString("base.http.transport.max-idle-conns") != "" {
        transport.MaxIdleConns = config.GetValueInt("base.http.transport.max-idle-conns")
    }

    if config.GetValueString("base.http.transport.max-idle-conns-per-host") != "" {
        transport.MaxIdleConnsPerHost = config.GetValueInt("base.http.transport.max-idle-conns-per-host")
    }

    if config.GetValueString("base.http.transport.max-conns-per-host") != "" {
        transport.MaxConnsPerHost = config.GetValueInt("base.http.transport.max-conns-per-host")
    }

    if config.GetValueString("base.http.transport.idle-conn-timeout") != "" {
        t, err := time.ParseDuration(config.GetValueString("base.http.transport.idle-conn-timeout"))
        if err != nil {
            logger.Warn("读取配置【base.http.transport.idle-conn-timeout】异常", err)
        } else {
            transport.IdleConnTimeout = t
        }
    }

    if config.GetValueString("base.http.transport.response-header-timeout") != "" {
        t, err := time.ParseDuration(config.GetValueString("base.http.transport.response-header-timeout"))
        if err != nil {
            logger.Warn("读取配置【base.http.transport.response-header-timeout】异常", err)
        } else {
            transport.ResponseHeaderTimeout = t
        }
    }

    if config.GetValueString("base.http.transport.expect-continue-timeout") != "" {
        t, err := time.ParseDuration(config.GetValueString("base.http.transport.expect-continue-timeout"))
        if err != nil {
            logger.Warn("读取配置【base.http.transport.expect-continue-timeout】异常", err)
        } else {
            transport.ExpectContinueTimeout = t
        }
    }

    if config.GetValueString("base.http.transport.max-response-header-bytes") != "" {
        transport.MaxResponseHeaderBytes = config.GetValueInt64("base.http.transport.max-response-header-bytes")
    }

    if config.GetValueString("base.http.transport.write-buffer-size") != "" {
        transport.WriteBufferSize = config.GetValueInt("base.http.transport.write-buffer-size")
    }

    if config.GetValueString("base.http.transport.read-buffer-size") != "" {
        transport.ReadBufferSize = config.GetValueInt("base.http.transport.read-buffer-size")
    }

    if config.GetValueString("base.http.transport.force-attempt-HTTP2") != "" {
        transport.ForceAttemptHTTP2 = config.GetValueBool("base.http.transport.force-attempt-HTTP2")
    }

    transport.DialContext = loadConfigOfDialContext()
    client.Transport = transport
}

func loadConfigOfDialContext() func(ctx context.Context, network, addr string) (net.Conn, error) {
    dialer := &net.Dialer{}
    if config.GetValueString("base.http.transport.dial-context.timeout") != "" {
        t, err := time.ParseDuration(config.GetValueString("base.http.transport.dial-context.timeout"))
        if err != nil {
            logger.Warn("读取配置【base.http.transport.dial-context.timeout】异常", err)
        } else {
            dialer.Timeout = t
        }
    }

    if config.GetValueString("base.http.transport.dial-context.keep-alive") != "" {
        t, err := time.ParseDuration(config.GetValueString("base.http.transport.dial-context.keep-alive"))
        if err != nil {
            logger.Warn("读取配置【base.http.transport.dial-context.keep-alive】异常", err)
        } else {
            dialer.KeepAlive = t
        }
    }
    return dialer.DialContext
}

func SetHttpClient(httpClientOuter *http.Client) {
    httpClient = httpClientOuter
}

func GetClient() *http.Client {
    return httpClient
}

func Do(httpRequest *http.Request) (int, http.Header, any, error) {
    ctx := context.Background()
    for _, hook := range NetHttpHooks {
        _ctx, httpHeader := hook.Before(ctx, httpRequest)
        httpRequest.Header = httpHeader
        ctx = _ctx
    }

    resp, err := httpClient.Do(httpRequest)
    rspCode, rspHead, rspData, err := doParseResponse(resp, err)
    for _, hook := range NetHttpHooks {
        hook.After(ctx, resp, rspCode, rspData, err)
    }
    return rspCode, rspHead, rspData, err
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

func PostForm(urlHost string, header http.Header, formMap map[string]any) (int, http.Header, any, error) {
    // 创建一个字节缓冲区，用于构建表单数据
    bodyBuffer := &bytes.Buffer{}
    writer := multipart.NewWriter(bodyBuffer)

    // 将表单数据写入multipart.Writer
    for k, v := range formMap {
        _ = writer.WriteField(k, fmt.Sprintf("%v", v))
    }

    // 关闭multipart.Writer，并获取最终的Content-Type
    err := writer.Close()
    if err != nil {
        logger.Error("关闭multipart.Writer时出错:", err)
        return 0, nil, nil, err
    }
    contentType := writer.FormDataContentType()

    // 创建一个请求对象
    httpReq, err := http.NewRequest("POST", urlHost, bodyBuffer)
    if err != nil {
        logger.Error("创建请求对象时出错:", err)
        return 0, nil, nil, err
    }

    // 设置请求头
    for k, values := range header {
        for _, value := range values {
            httpReq.Header.Add(k, value)
        }
    }
    httpReq.Header.Set("Content-Type", contentType)

    ctx := context.Background()
    for _, hook := range NetHttpHooks {
        _ctx, httpHeader := hook.Before(ctx, httpReq)
        httpReq.Header = httpHeader
        ctx = _ctx
    }

    resp, err := httpClient.Do(httpReq)
    rspCode, rspHead, rspData, err := doParseResponse(resp, err)
    for _, hook := range NetHttpHooks {
        hook.After(ctx, resp, rspCode, rspData, err)
    }

    return rspCode, rspHead, rspData, err
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
    ctx := context.Background()

    for _, hook := range NetHttpHooks {
        _ctx, httpHeader := hook.Before(ctx, httpRequest)
        httpRequest.Header = httpHeader
        ctx = _ctx
    }

    httpResponse, err := httpClient.Do(httpRequest)
    rspCode, rspHead, rspData, err := doParseResponse(httpResponse, err)

    for _, hook := range NetHttpHooks {
        hook.After(ctx, httpResponse, rspCode, rspData, err)
    }
    return rspCode, rspHead, rspData, err
}

func doParseResponse(httpResponse *http.Response, err error) (int, http.Header, any, error) {
    if err != nil && httpResponse == nil {
        log.Printf("Error sending request to API endpoint. %+v", err)
        return -1, nil, nil, &NetError{ErrMsg: "Error sending request, err" + err.Error()}
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
            body, _ := io.ReadAll(httpResponse.Body)
            return code, headers, nil, &NetError{ErrMsg: "remote error, url: code " + strconv.Itoa(code) + ", message: " + string(body)}
        }

        // We have seen inconsistencies even when we get 200 OK response
        body, err := io.ReadAll(httpResponse.Body)
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
    ctx := context.Background()

    for _, hook := range NetHttpHooks {
        _ctx, httpHeader := hook.Before(ctx, httpRequest)
        httpRequest.Header = httpHeader
        ctx = _ctx
    }

    httpResponse, err := httpClient.Do(httpRequest)
    rspCode, _, rspData, err := doParseResponse(httpResponse, err)

    for _, hook := range NetHttpHooks {
        hook.After(ctx, httpResponse, rspCode, rspData, err)
    }
    return err
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
