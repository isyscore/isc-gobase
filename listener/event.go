package listener

var DefaultGroup = "DefaultGroup"
var EventOfServerRunStart = "event_of_server_run_start"
var EventOfServerRunFinish = "event_of_server_run_finish"
var EventOfServerStop = "event_of_server_stop"
var EventOfConfigChange = "event_of_config_change"

// ServerRunStartEvent 服务开始启动事件, 对应：event_of_server_run_start
type ServerRunStartEvent struct{}

// ServerRunFinishEvent 服务完成启动事件, 对应：event_of_server_run_finish
type ServerRunFinishEvent struct{}

// ServerStopEvent 服务关闭事件, 对应：event_of_server_stop
type ServerStopEvent struct{}

// ConfigChangeEvent 配置变更事件, 对应：event_of_config_change
type ConfigChangeEvent struct {
	Key   string
	Value string
}

func (e ServerRunStartEvent) Name() string {
	return EventOfServerRunStart
}

func (e ServerRunStartEvent) Group() string {
	return DefaultGroup
}

func (e ServerRunFinishEvent) Name() string {
	return EventOfServerRunFinish
}

func (e ServerRunFinishEvent) Group() string {
	return DefaultGroup
}

func (e ServerStopEvent) Name() string {
	return EventOfServerStop
}

func (e ServerStopEvent) Group() string {
	return DefaultGroup
}

func (e ConfigChangeEvent) Name() string {
	return EventOfConfigChange
}

func (e ConfigChangeEvent) Group() string {
	return DefaultGroup
}
