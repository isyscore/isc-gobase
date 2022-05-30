package listener

var EventOfServerFinish = "event_of_server_finish"
var EventOfServerStop = "event_of_server_stop"

// ServerFinishEvent 服务启动完成事件
type ServerFinishEvent struct{}

// ServerStopEvent 服务关闭事件
type ServerStopEvent struct{}

func (e ServerFinishEvent) Name() string {
	return EventOfServerFinish
}

func (e ServerStopEvent) Name() string {
	return EventOfServerStop
}
