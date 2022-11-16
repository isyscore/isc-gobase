package debug

import (
	"context"

	"github.com/isyscore/isc-gobase/extend/etcd"
	"github.com/isyscore/isc-gobase/logger"
	etcdClientV3 "go.etcd.io/etcd/client/v3"
)

var etcdClient *etcd.EtcdClientWrap
var keyListenerMap map[string][]KeyListener

type KeyListener func(key string, value string)

func Init() {
	if etcdClient != nil {
		return
	}
	etcdClientTem, err := etcd.NewEtcdClient()
	if err != nil {
		logger.Error("etcd初始化失败 %v", err.Error())
	}

	if etcdClientTem == nil {
		return
	}
	etcdClient = etcdClientTem
	go startWatch()
}

func AddWatcher(key string, keyListener KeyListener) {
	if etcdClient == nil {
		logger.Error("请先调用方法：debug.Init() 用于初始化调试模式")
		return
	}

	if keyListenerMap == nil {
		keyListenerMap = map[string][]KeyListener{}
	}
	if eventWatchers, exist := keyListenerMap[key]; exist {
		eventWatchers = append(eventWatchers, keyListener)
		keyListenerMap[key] = eventWatchers
	} else {
		eventWatchers = []KeyListener{}
		eventWatchers = append(eventWatchers, keyListener)
		keyListenerMap[key] = eventWatchers
	}
}

func Update(key, value string) {
	ctx := context.Background()
	_, err := etcdClient.Put(ctx, key, value)
	if err != nil {
		logger.Error("更新调试配置报错")
		return
	}
}

func startWatch() {
	for key, listeners := range keyListenerMap {
		watchKey := key
		keyListeners := listeners
		go func() {
			var currentVersion int64 = 0
			ctx := context.Background()

			// 首次启动获取最新的
			rsp, err := etcdClient.Get(ctx, watchKey)
			if err != nil {
				logger.Error("获取etcd的key异常, %v", err.Error())
				return
			} else {
				for _, kv := range rsp.Kvs {
					if watchKey == string(kv.Key) {
						currentVersion = kv.ModRevision
						go notifyWatcher(watchKey, string(kv.Value), keyListeners)
					}
				}
			}

			// 根据本地保存的最新版本进行watch
			watchRsp := etcdClient.Watch(ctx, watchKey, etcdClientV3.WithRev(currentVersion))
			for res := range watchRsp {
				for _, event := range res.Events {
					if watchKey != string(event.Kv.Key) {
						continue
					}
					latestModVersion := event.Kv.ModRevision

					if currentVersion == event.Kv.ModRevision {
						continue
					} else if currentVersion < latestModVersion {
						currentVersion = latestModVersion
						go notifyWatcher(watchKey, string(event.Kv.Value), keyListeners)
					}
				}
			}
		}()
	}
}

func notifyWatcher(key, value string, listeners []KeyListener) {
	for _, listener := range listeners {
		listener(key, value)
	}
}
