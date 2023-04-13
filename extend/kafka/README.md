# kafka
这个是基于 github/Shopify/sarama 这个客户端的封装

## 示例

### 配置
```yaml
base:
  kafka:
    addrs:
      - xx.xx.xx.xx:xxx
    # 是否激活，默认：false
    enable: true
    # 生产者配置
    producer:
      # 默认false
      return-success: true
```
### 代码

```go
// 生产者
package main

import "github.com/isyscore/isc-gobase/extend/kafka"

// 生产者
func TestProducerNew(t *testing.T) {
    // 获取同步生产者
    producer, err := kafka.NewSyncProducer()
    if err != nil {
        logger.Error("异常：%v", err.Error())
        return
    }
    
    // 发送消息
    msg := &sarama.ProducerMessage{
        Topic: "my_topic",
        Value: sarama.StringEncoder("Hello, world!"),
    }
    partition, offset, err := producer.SendMessage(msg)
    if err != nil {
        logger.Error("%v", err.Error())
        return
    }
    logger.Info("Message sent to partition %d at offset %d\n", partition, offset)
}

// 消费者
func TestConsumerNew(t *testing.T) {
    // 创建 Kafka 消费者
    consumer, err := kafka.NewConsumer()
    if err != nil {
        logger.Error("%v", err.Error())
        return
    }
    
    // 订阅主题
    topic := "my_topic"
    partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
    if err != nil {
        logger.Error("%v", err.Error())
        return
    }
    
    // 处理消息
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)
    for {
        select {
        case msg := <-partitionConsumer.Messages():
            logger.Info("Received message: %s\n", string(msg.Value))
        case err := <-partitionConsumer.Errors():
            logger.Info("Error: %s\n", err.Error())
        case <-signals:
            return
        }
    }
}
```

## 全部api
提供如下的一些api，简化生成
```go
func NewClient() (sarama.Client, error) {}

func NewAsyncProducer() (sarama.AsyncProducer, error) {}
func NewSyncProducer() (sarama.SyncProducer, error) {}

func NewConsumer() (sarama.Consumer, error) {}
func NewConsumerGroup(groupId string) (sarama.ConsumerGroup, error) {}

func NewClusterAdmin() (sarama.ClusterAdmin, error) {}
func GetKafkaConfig() *sarama.Config {}
```

## 全部配置 
```yaml
base:
  kafka:
    addrs:
      - {ip}:{port}
      - {ip}:{port}
      - {ip}:{port}
    # 是否激活，默认：false
    enable: false
    # 默认 sarama
    client-id: sarama
    # 默认 256
    channel-buffer-size: 256
    # 默认 true
    api-versions-request: true
    # 默认 V1_0_0_0，版本格式为V{x}_{x}_{x}_{x}；版本是否存在请见 Shopify/sarama 代码中的util.go包的版本
    version: V1_0_0_0
    admin:
      # 默认5
      retry-max: 5
      # 默认100ms
      retry-backoff: 100ms
      # 默认3s
      timeout: 3s
    net:
      # 默认5
      max-open-requests: 5
      # 默认3s
      dial-timeout: 3s
      # 默认3s
      read-timeout: 3s
      # 默认3s
      write-timeout: 3s
      # 默认true
      SASL-handshake: true
      # 默认0
      SASL-version: 0
    metadata:
      # 默认 3
      retry-max: 3
      # 默认250ms
      retry-backoff: 250ms
      # 默认10分钟，即10m
      refresh-frequency: 10m
      # 默认 true
      full: true
      # 默认 true
      allow-auto-topic-creation: true
    producer:
      # 默认1000000
      max-message-bytes: 1000000
      # 默认1，只可以为：-1, 0, 1
      required-acks: 1
      # 10s
      timeout: 10s
      # 默认3
      retry-max: 3
      # 默认100ms
      retry-backoff: 100ms
      # 默认true
      return-errors: true
      # 默认false
      return-success: false
      # 默认-1000
      compression-level: -1000
      # 默认1分钟
      transaction-timeout: 1m
      # 默认50
      transaction-retry-max: 50
      # 默认100毫秒
      transaction-retry-backoff: 100ms
    consumer:
      # 默认1
      fetch-min: 1
      # 默认1048576，即：1024*1024
      fetch-default: 1048576 
      # 默认2s
      retry-backoff: 2s
      # 默认500ms
      max-wait-time: 500ms
      # 默认100ms
      max-processing-time: 100ms
      # 默认false
      return-errors: false
      # 默认false
      offsets-auto-commit-enable: false
      # 默认1秒
      offsets-auto-commit-interval: 1s
      # 默认-1
      offsets-initial: -1
      # 默认3
      offsets-retry-max: 3
      # 消费组配置
      group:
        # 默认10s
        session-timeout: 10s
        # 默认3s
        heartbeat-interval: 3s
        # 默认60s
        rebalance-timeout: 60s
        # 默认4
        rebalance-retry-max: 4
        # 默认2秒
        rebalance-retry-backoff: 2s
        # 默认true
        reset-invalid-offsets: true
```
