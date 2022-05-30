## listener
listener包：是isc-gobase的事件监听模块

### 用法

提供api
```go
// 添加事件监听器
listener.AddListener(eventName string, eventListener EventListener)

// 发布事件
listener.PublishEvent(event listener.BaseEvent)
```

### 快速使用
#### 1. 定义事件
事件要实现接口 `listener.BaseEvent` 的方法

```go
type Event1 struct {
    Company string
}

func (e1 Event1) Name() string {
    return "event1"
}
```

```go
import (
    "fmt"
    "github.com/isyscore/isc-gobase/listener"
    "testing"
)

func TestPublish(t *testing.T) {
    listener.AddListener("event1", Event1Lister1)
    listener.AddListener("event1", Event1Lister2)
    listener.AddListener("event1", Event1Lister3)

    listener.PublishEvent(Event1{Company: "公司"})
}

// 事件监听器1
func Event1Lister1(event listener.BaseEvent) {
    ev := event.(Event1)
    fmt.Println("Event1Lister1: " + ev.Company)
}

// 事件监听器2
func Event1Lister2(event listener.BaseEvent) {
    ev := event.(Event1)
    fmt.Println("Event1Lister2: " + ev.Company)
}

// 事件监听器3
func Event1Lister3(event listener.BaseEvent) {
    ev := event.(Event1)
    fmt.Println("Event1Lister3: " + ev.Company)
}
```

### 内置监听器
isc-gobase内置了几类事件
- ServerPostEvent: 服务启动完成事件
- ServerFinishEvent: 服务关闭事件

常用示例：
```go
// 添加服务器启动完成事件监听
listener.AddListener(listener.EventOfServerFinish, func(event listener.BaseEvent) {
    logger.Info("应用启动完成")
})
```
