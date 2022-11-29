package store

import (
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/goid"
	"net/http"
)

var RequestStorage goid.LocalStorage
var MdcStorage goid.LocalStorage

func init() {
	RequestStorage = goid.NewLocalStorage()
	MdcStorage = goid.NewLocalStorage()
}

func RequestHeadAdd(key, value string) {
	req := RequestStorage.Get()
	if req == nil {
		return
	}
	_req := req.(*http.Request)
	_req.Header.Set(key, value)
}

func GetRequest() *http.Request {
	req := RequestStorage.Get()
	if req == nil {
		return nil
	}
	return req.(*http.Request)
}

func GetHeader() http.Header {
	req := RequestStorage.Get()
	if req == nil {
		return nil
	}
	reqS := req.(*http.Request)
	return reqS.Header
}

func GetRemoteAddr() string {
	req := RequestStorage.Get()
	if req == nil {
		return ""
	}
	reqS := req.(*http.Request)
	return reqS.RemoteAddr
}

func GetHeaderWithKey(headKey string) string {
	req := RequestStorage.Get()
	if req == nil {
		return ""
	}
	reqS := req.(*http.Request)
	return reqS.Header.Get(headKey)
}

func GetTraceId() string {
	req := RequestStorage.Get()
	if req == nil {
		return ""
	}
	reqS := req.(*http.Request)
	return reqS.Header.Get(constants.TRACE_HEAD_ID)
}

func GetUserId() string {
	req := RequestStorage.Get()
	if req == nil {
		return ""
	}
	reqS := req.(*http.Request)
	return reqS.Header.Get(constants.TRACE_HEAD_USER_ID)
}
