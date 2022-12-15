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
	headKvMap := headKeyValueStorage.Get()

	if headKvMap == nil {
		headKvMap = cmap.New()
	}
	kvMap := headKvMap.(cmap.ConcurrentMap)
	kvMap.Set(constants.TRACE_HEAD_ID, httpHead.Get(constants.TRACE_HEAD_ID))
	kvMap.Set(constants.TRACE_HEAD_RPC_ID, httpHead.Get(constants.TRACE_HEAD_RPC_ID))
	kvMap.Set(constants.TRACE_HEAD_SAMPLED, httpHead.Get(constants.TRACE_HEAD_SAMPLED))
	kvMap.Set(constants.TRACE_HEAD_USER_ID, httpHead.Get(constants.TRACE_HEAD_USER_ID))
	kvMap.Set(constants.TRACE_HEAD_USER_NAME, httpHead.Get(constants.TRACE_HEAD_USER_NAME))
	kvMap.Set(constants.TRACE_HEAD_REMOTE_IP, httpHead.Get(constants.TRACE_HEAD_REMOTE_IP))
	kvMap.Set(constants.TRACE_HEAD_REMOTE_APPNAME, httpHead.Get(constants.TRACE_HEAD_REMOTE_APPNAME))
	kvMap.Set(constants.TRACE_HEAD_ORIGNAL_URL, httpHead.Get(constants.TRACE_HEAD_ORIGNAL_URL))
	headKeyValueStorage.Set(kvMap)
}

func Put(key string, value any) {
	headKvMap := headKeyValueStorage.Get()

	if headKvMap == nil {
		headKvMap = map[string]any{}
	}
	kvMap := headKvMap.(cmap.ConcurrentMap)
	kvMap.Set(key, value)
	headKeyValueStorage.Set(kvMap)
}

func Get(key string) any {
	headKvMap := headKeyValueStorage.Get()

	if headKvMap == nil {
		headKvMap = cmap.New()
		headKeyValueStorage.Set(headKvMap)
		return ""
	}
	kvMap := headKvMap.(cmap.ConcurrentMap)
	data, exist := kvMap.Get(key)
	if exist {
		return data
	}
	return ""
}

func Keys() []string {
	headKvMap := headKeyValueStorage.Get()

	if headKvMap == nil {
		return []string{}
	}

	kvMap := headKvMap.(cmap.ConcurrentMap)
	return kvMap.Keys()
}

func Clean() {
	headKeyValueStorage.Del()
}
