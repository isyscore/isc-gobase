package store

import (
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/goid"
	cmap "github.com/orcaman/concurrent-map"
	"net/http"
)

var headKeyValueStorage goid.LocalStorage

func init() {
	headKeyValueStorage = goid.NewLocalStorage()
}

func PutFromHead(httpHead http.Header) {
	mdcMapTem := headKeyValueStorage.Get()

	if mdcMapTem == nil {
		mdcMapTem = cmap.New()
	}
	mdcMap := mdcMapTem.(cmap.ConcurrentMap)
	mdcMap.Set(constants.TRACE_HEAD_ID, httpHead.Get(constants.TRACE_HEAD_ID))
	mdcMap.Set(constants.TRACE_HEAD_RPC_ID, httpHead.Get(constants.TRACE_HEAD_RPC_ID))
	mdcMap.Set(constants.TRACE_HEAD_SAMPLED, httpHead.Get(constants.TRACE_HEAD_SAMPLED))
	mdcMap.Set(constants.TRACE_HEAD_USER_ID, httpHead.Get(constants.TRACE_HEAD_USER_ID))
	mdcMap.Set(constants.TRACE_HEAD_USER_NAME, httpHead.Get(constants.TRACE_HEAD_USER_NAME))
	mdcMap.Set(constants.TRACE_HEAD_REMOTE_IP, httpHead.Get(constants.TRACE_HEAD_REMOTE_IP))
	mdcMap.Set(constants.TRACE_HEAD_REMOTE_APPNAME, httpHead.Get(constants.TRACE_HEAD_REMOTE_APPNAME))
	mdcMap.Set(constants.TRACE_HEAD_ORIGNAL_URL, httpHead.Get(constants.TRACE_HEAD_ORIGNAL_URL))
	headKeyValueStorage.Set(mdcMap)
}

func Put(key string, value any) {
	mdcMapTem := headKeyValueStorage.Get()

	if mdcMapTem == nil {
		mdcMapTem = map[string]any{}
	}
	mdcMap := mdcMapTem.(cmap.ConcurrentMap)
	mdcMap.Set(key, value)
	headKeyValueStorage.Set(mdcMap)
}

func Get(key string) any {
	mdcMapTem := headKeyValueStorage.Get()

	if mdcMapTem == nil {
		mdcMapTem = cmap.New()
		headKeyValueStorage.Set(mdcMapTem)
		return ""
	}
	mdcMap := mdcMapTem.(cmap.ConcurrentMap)
	data, exist := mdcMap.Get(key)
	if exist {
		return data
	}
	return ""
}

func Clean() {
	headKeyValueStorage.Del()
}
