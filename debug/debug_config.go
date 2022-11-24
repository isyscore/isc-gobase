package debug

import (
	"context"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/logger"
	etcdClientV3 "go.etcd.io/etcd/client/v3"
	"net"
	"os"
	"strings"
	"time"
)

const (
	DEBUG_ETCD_ENDPOINTS  = "debug.etcd.endpoints"
	DEBUG_ETCD_USER       = "debug.etcd.user"
	DEBUG_ETCD_PASSWORD   = "debug.etcd.password"
	DEFAULT_ETCD_ENDPOINT = "etcd-service:22379"
)

var etcdClient *etcdClientV3.Client
var keyListenerMap map[string][]KeyListener

type KeyListener func(key string, value string)

func Init() {
	if !config.GetValueBoolDefault("base.debug.enable", true) {
		return
	}
	InitWithParameter(GetEtcdConfig())
}

func InitWithParameter(endpoints []string, user, password string) {
	if etcdClient != nil {
		return
	}

	for _, endpoint := range endpoints {
		if !ipPortAvailable(endpoint) {
			// 如果想使用调试模式，请在环境变量里面或者配置文件里面配置如下
			// debug.etcd.endpoints：多个{ip}:{port}格式，中间以逗号（英文逗号）分隔
			// debug.etcd.user
			// debug.etcd.password
			logger.Warn("调试模式【%v】不可用，调试功能暂时不支持", endpoint)
			return
		}
	}

	_etcdClient := getEtcdClient(endpoints, user, password)
	if _etcdClient == nil {
		return
	}

	etcdClient = _etcdClient
}

func ipPortAvailable(ipAndPort string) bool {
	conn, err := net.DialTimeout("tcp", ipAndPort, 2 * time.Second)
	if err != nil {
		return false
	} else {
		if conn != nil {
			conn.Close()
			return true
		} else {
			return false
		}
	}
}

func getEtcdClient(etcdPoints []string, user, password string) *etcdClientV3.Client {
	// 客户端配置
	etcdCfg := etcdClientV3.Config{
		Endpoints: etcdPoints,
		Username:  user,
		Password:  password,
	}

	etcdClient, err := etcdClientV3.New(etcdCfg)
	if err != nil {
		logger.Error("生成etcd-client失败：%v", err.Error())
		return nil
	}

	return etcdClient
}

// 优先级：vm配置 > 环境变量
func GetEtcdConfig() ([]string, string, string) {
	etcdEndpointStr := os.Getenv(DEBUG_ETCD_ENDPOINTS)
	etcdEndpointStrOfConfig := config.GetValueString(DEBUG_ETCD_ENDPOINTS)
	if etcdEndpointStrOfConfig != "" {
		etcdEndpointStr = etcdEndpointStrOfConfig
	}

	if etcdEndpointStr == "" {
		etcdEndpointStr = DEFAULT_ETCD_ENDPOINT
	}

	etcdEndpointsOriginal := strings.Split(etcdEndpointStr, ",")
	var etcdEndpoints []string
	for _, etcdEndpoint := range etcdEndpointsOriginal {
		etcdEndpoint = strings.TrimSpace(etcdEndpoint)
		if strings.HasPrefix(etcdEndpoint, "http://") {
			etcdEndpoint = etcdEndpoint[len("http://"):]
		}

		etcdEndpoints = append(etcdEndpoints, etcdEndpoint)
	}

	etcdUser := os.Getenv(DEBUG_ETCD_USER)
	etcdUserOfConfig := config.GetValueString(DEBUG_ETCD_USER)
	if etcdUserOfConfig != "" {
		etcdUser = etcdUserOfConfig
	}

	etcdPassword := os.Getenv(DEBUG_ETCD_PASSWORD)
	etcdPasswordOfConfig := config.GetValueString(DEBUG_ETCD_PASSWORD)
	if etcdPasswordOfConfig != "" {
		etcdPassword = etcdPasswordOfConfig
	}

	return etcdEndpoints, etcdUser, etcdPassword
}

func AddWatcher(key string, keyListener KeyListener) {
	if etcdClient == nil {
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

func StartWatch() {
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
