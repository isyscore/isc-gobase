# emqx

对业内的emqx客户端进行配置化封装，用于简化获取
### 全部配置
```yaml
base:
  emqx:
    # 是否开启emqx，默认关闭
    enable: true
    servers:
      # 域名格式1
      - "tcp://{user}:{password}@{host}:{port}"
      # 域名格式2
      - "tcp://{host}:{port}"
    client-id: "xxxx"
    username: "xxxx"
    password: "xxxx"
    # 是否清理session，默认为true
    clean-session: true
    order: true
    will-enabled: true
    will-topic: "xxx-topic"
    will-qos: 0
    will-retained: false
    protocol-version: 0
    keep-alive: 30
    ping-timeout: "10s"
    connect-timeout: "30s"
    max-reconnect-interval: "10m"
    auto-reconnect: true
    connect-retry-interval: "30s"
    connect-retry: false
    write-timeout: 0
    resume-subs: false
    max-resume-pub-in-flight: 0
    auto-ack-disabled: false
```
提供封装的 `emqx客户端api`
```go
func NewEmqxClient() (mqtt.Client, error) {}
```
#### 示例：
```yaml

```yaml
base:
  emqx:
    enable: true
    servers:
      - "tcp://{host}:{port}"
    client-id: "xxxx"
    username: "xxxx"
    password: "xxxx"
```
```go
import (
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/isyscore/isc-gobase/extend/emqx"
 )

// 消息回调函数
var msgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("TOPIC: %s\n", msg.Topic())
    fmt.Printf("MSG: %s\n", msg.Payload())
}

func TestConnect(t *testing.T) {
    // 获取emqx的客户端
    emqxClient, _ := emqx.NewEmqxClient()
    
    // 订阅主题
    if token := emqxClient.Subscribe("testtopic/#", 0, msgHandler); token.Wait() && token.Error() != nil {
        fmt.Println(token.Error())
        os.Exit(1)
    }
    
    // 发布消息
    token := emqxClient.Publish("testtopic/1", 0, false, "Hello World")
    token.Wait()
    
    time.Sleep(1 * time.Second)
    
    // 取消订阅
    if token := emqxClient.Unsubscribe("testtopic/#"); token.Wait() && token.Error() != nil {
        fmt.Println(token.Error())
        os.Exit(1)
    }
    
    // 断开连接
    emqxClient.Disconnect(250)
    time.Sleep(1 * time.Second)
}
```
