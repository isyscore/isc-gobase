package listener

var EventOfServerFinish = "event_of_server_finish"
var EventOfServerStop = "event_of_server_stop"
var EventOfConfigChange = "event_of_config_change"

// ServerFinishEvent 服务启动完成事件, 对应：event_of_server_finish
type ServerFinishEvent struct{}

// ServerStopEvent 服务关闭事件, 对应：event_of_server_stop
type ServerStopEvent struct{}

// ConfigChangeEvent 配置变更事件, 对应：event_of_config_change
type ConfigChangeEvent struct {
	Key   string
	Value string
}

func (e ServerFinishEvent) Name() string {
	return EventOfServerFinish
}

func (e ServerStopEvent) Name() string {
	return EventOfServerStop
}

func (e ConfigChangeEvent) Name() string {
	return EventOfConfigChange
}
