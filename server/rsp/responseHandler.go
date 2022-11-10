package rsp

import (
	"bytes"
	"encoding/json"
	"github.com/isyscore/isc-gobase/config"
	"io"
	"net/http"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqPrint := config.GetValueBoolDefault("base.server.request.print.enable", false)
		rspPrint := config.GetValueBoolDefault("base.server.response.print.enable", false)
		expPrint := config.GetValueBoolDefault("base.server.exception.print.enable", false)

		if !reqPrint && !rspPrint && !expPrint {
			return
		}

		// 开始时间
		startTime := time.Now()

		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("read request body failed,err = %s.", err)
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 状态码
		statusCode := c.Writer.Status()

		var body any
		bodyStr := string(data)
		if "" != bodyStr && unsafe.Sizeof(bodyStr) < 10240 {
			if strings.HasPrefix(bodyStr, "{") && strings.HasSuffix(bodyStr, "}") {
				bodys := map[string]any{}
				_ = isc.StrToObject(bodyStr, &bodys)
				body = bodys
			} else if strings.HasPrefix(bodyStr, "[") && strings.HasSuffix(bodyStr, "]") {
				var bodys []any
				_ = isc.StrToObject(bodyStr, &bodys)
				body = bodys
			}
		}

		request := Request{
			Method:     c.Request.Method,
			Uri:        c.Request.RequestURI,
			Ip:         c.ClientIP(),
			Parameters: c.Params,
			Headers:    c.Request.Header,
			Body:       body,
		}

		errMessage := ErrorMessage{
			Request:    request,
			StatusCode: statusCode,
			Cost:       time.Now().Sub(startTime).String(),
		}

		responseMessage := Response{
			Request:    request,
			StatusCode: statusCode,
			Cost:       time.Now().Sub(startTime).String(),
		}

		if reqPrint && !rspPrint && !expPrint {
			printReq(request.Uri, request)
		}

		// 1xx和2xx都是成功
		if (statusCode >= 300) && statusCode != 0 {
			datas := config.BaseCfg.Server.Exception.Print.Exclude
			for _, code := range datas {
				if code == statusCode {
					return
				}
			}
			if expPrint {
				logger.Error("返回异常, result：%v", isc.ObjectToJson(errMessage))
			}
		} else {
			var response DataResponse[any]
			if err := json.Unmarshal([]byte(blw.body.String()), &response); err != nil {
				return
			} else {
				c.Writer.Header().Add("isc-biz-code", isc.ToString(response.Code))
				c.Writer.Header().Add("isc-biz-message", response.Message)
				if response.Code != 0 && response.Code != 200 {
					errMessage.Response = response
					if expPrint {
						logger.Error("返回异常, result：%v", isc.ObjectToJson(errMessage))
					}
				} else {
					responseMessage.Response = response
					if rspPrint {
						printRsq(request.Uri, responseMessage)
					}
				}
			}
		}
	}
}

func printReq(requestUri string, requestData Request) {
	includeUri := config.GetValueArray("base.server.request.print.include-uri")
	printFlag := false
	if len(includeUri) != 0 {
		for _, uri := range includeUri {
			if strings.HasPrefix(requestUri, isc.ToString(uri)) {
				printFlag = true
				break
			}
		}
	}

	excludeUri := config.GetValueArray("base.server.request.print.exclude-uri")
	if len(excludeUri) != 0 {
		for _, uri := range excludeUri {
			if strings.HasPrefix(requestUri, isc.ToString(uri)) {
				printFlag = false
				break
			}
		}
	}

	reqLogLevel := config.GetValueString("base.server.request.print.level")
	if printFlag {
		logger.Record(reqLogLevel, "请求：%v", isc.ObjectToJson(requestData))
	}
	return
}

func printRsq(requestUri string, responseMessage Response) {
	includeUri := config.GetValueArray("base.server.response.print.include-uri")
	printFlag := false
	if len(includeUri) != 0 {
		for _, uri := range includeUri {
			if strings.HasPrefix(requestUri, isc.ToString(uri)) {
				printFlag = true
				break
			}
		}
	}

	excludeUri := config.GetValueArray("base.server.response.print.exclude-uri")
	if len(excludeUri) != 0 {
		for _, uri := range excludeUri {
			if strings.HasPrefix(requestUri, isc.ToString(uri)) {
				printFlag = false
				break
			}
		}
	}

	rspLogLevel := config.GetValueString("base.server.response.print.level")
	if printFlag {
		logger.Record(rspLogLevel, "响应：%v", isc.ObjectToJson(responseMessage))
	}
}

type Request struct {
	Method     string
	Uri        string
	Ip         string
	Headers    http.Header
	Parameters gin.Params
	Body       any
}

type ErrorMessage struct {
	Request    Request
	Response   DataResponse[any]
	Cost       string
	StatusCode int
}

type Response struct {
	Request    Request
	Response   DataResponse[any]
	Cost       string
	StatusCode int
}
