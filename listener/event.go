package listener

var EventOfConfigPerInit = "event_of_config_per_init"
var EventOfConfigPostInit = "event_of_config_post_init"
var EventOfServerPost = "event_of_server_post"
var EventOfServerStop = "event_of_server_stop"

// ConfigInitPreEvent config文件加载前事件
type ConfigInitPreEvent struct{}

// ConfigInitPostEvent config文件加载完成事件
type ConfigInitPostEvent struct{}

// ServerPostEvent 服务启动完成事件
type ServerPostEvent struct{}

// ServerStopEvent 服务关闭事件
type ServerStopEvent struct{}

func (e ConfigInitPreEvent) Name() string {
	return EventOfConfigPerInit
}

func (e ConfigInitPostEvent) Name() string {
	return EventOfConfigPostInit
}

func (e ServerPostEvent) Name() string {
	return EventOfServerPost
}

func (e ServerStopEvent) Name() string {
	return EventOfServerStop
}
