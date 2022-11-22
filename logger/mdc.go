package logger

import (
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/goid"
	"net/http"
)

var MdcStorage goid.LocalStorage

func init() {
	MdcStorage = goid.NewLocalStorage()
}

func PutHead(httpHead http.Header) {
	mdcMapTem := MdcStorage.Get()

	if mdcMapTem == nil {
		mdcMapTem = map[string]any{}
	}
	mdcMap := mdcMapTem.(map[string]any)
	mdcMap[constants.TRACE_HEAD_ID] = httpHead.Get(constants.TRACE_HEAD_ID)
	mdcMap[constants.TRACE_HEAD_RPC_ID] = httpHead.Get(constants.TRACE_HEAD_RPC_ID)
	mdcMap[constants.TRACE_HEAD_SAMPLED] = httpHead.Get(constants.TRACE_HEAD_SAMPLED)
	mdcMap[constants.TRACE_HEAD_USER_ID] = httpHead.Get(constants.TRACE_HEAD_USER_ID)
	mdcMap[constants.TRACE_HEAD_USER_NAME] = httpHead.Get(constants.TRACE_HEAD_USER_NAME)
	mdcMap[constants.TRACE_HEAD_REMOTE_IP] = httpHead.Get(constants.TRACE_HEAD_REMOTE_IP)
	mdcMap[constants.TRACE_HEAD_REMOTE_APPNAME] = httpHead.Get(constants.TRACE_HEAD_REMOTE_APPNAME)
	mdcMap[constants.TRACE_HEAD_ORIGNAL_URL] = httpHead.Get(constants.TRACE_HEAD_ORIGNAL_URL)
	MdcStorage.Set(mdcMap)
}

func PutMdc(key string, value any) {
	mdcMapTem := MdcStorage.Get()

	if mdcMapTem == nil {
		mdcMapTem = map[string]any{}
	}
	mdcMap := mdcMapTem.(map[string]any)
	mdcMap[key] = value
	MdcStorage.Set(mdcMap)
}

func GetMdc(key string) any {
	mdcMapTem := MdcStorage.Get()

	if mdcMapTem == nil {
		mdcMapTem = map[string]any{}
		MdcStorage.Set(mdcMapTem)
		return ""
	}
	mdcMap := mdcMapTem.(map[string]any)
	data, exist := mdcMap[key]
	if exist {
		return data
	}
	return ""
}
