package listener

var EventOfServerPost = "event_of_server_post"
var EventOfServerStop = "event_of_server_stop"

// ServerPostEvent 服务启动完成事件
type ServerPostEvent struct{}

// ServerStopEvent 服务关闭事件
type ServerStopEvent struct{}

func (e ServerPostEvent) Name() string {
	return EventOfServerPost
}

func (e ServerStopEvent) Name() string {
	return EventOfServerStop
}
