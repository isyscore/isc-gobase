package test

import (
	"github.com/Shopify/sarama"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/extend/kafka"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"testing"
)

func TestCreateTopic(t *testing.T) {
	// 创建 Kafka 配置
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0

	// 创建 Kafka 管理者
	admin, err := sarama.NewClusterAdmin([]string{"10.30.30.78:29092"}, config)
	if err != nil {
		logger.Error("%v", err.Error())
	}
	defer admin.Close()

	// 创建 Topic
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	err = admin.CreateTopic("my_topic", topicDetail, true)
	if err != nil {
		logger.Error("%v", err.Error())
	}
}

func TestProducerOriginal(t *testing.T) {
	// 配置 Kafka 生产者
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	client, err := sarama.NewClient([]string{"10.30.30.78:29092"}, config)

	// 创建 Kafka 生产者
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		logger.Error("%v", err.Error())
	}
	defer producer.Close()

	// 发送消息
	msg := &sarama.ProducerMessage{
		Topic: "my_topic",
		Value: sarama.StringEncoder("Hello, world!"),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		logger.Error("%v", err.Error())
	}
	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}

func TestConsumerOriginal(t *testing.T) {
	// 配置 Kafka 消费者
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// 创建 Kafka 消费者
	consumer, err := sarama.NewConsumer([]string{"10.30.30.78:29092"}, config)
	if err != nil {
		logger.Error("%v", err.Error())
	}
	defer consumer.Close()

	// 订阅主题
	topic := "my_topic"
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		logger.Error("%v", err.Error())
	}
	defer partitionConsumer.Close()

	// 处理消息
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Received message: %s\n", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %s\n", err.Error())
		case <-signals:
			return
		}
	}
}

func TestConfig(t *testing.T) {
	config.LoadYamlFile("./application-parameter.yaml")
	if config.GetValueBoolDefault("base.kafka.enable", false) {
		err := config.GetValueObject("base.kafka", &config.KafkaCfg)
		if err != nil {
			return
		}
	}

	//fmt.Println(config.KafkaCfg)
}

func TestProducerNew(t *testing.T) {
	config.LoadYamlFile("./application-simple.yaml")
	if config.GetValueBoolDefault("base.kafka.enable", false) {
		err := config.GetValueObject("base.kafka", &config.KafkaCfg)
		if err != nil {
			return
		}
	}

	producer, err := kafka.NewSyncProducer()
	if err != nil {
		logger.Error("异常：%v", err.Error())
		return
	}

	// 发送消息
	msg := &sarama.ProducerMessage{
		Topic: "my_topic_demo",
		Value: sarama.StringEncoder("Hello, world!"),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		logger.Error("%v", err.Error())
		return
	}
	logger.Info("Message sent to partition %d at offset %d\n", partition, offset)
}

func TestProducerNew2(t *testing.T) {
	config.LoadYamlFile("./application-producer.yaml")
	if config.GetValueBoolDefault("base.kafka.enable", false) {
		err := config.GetValueObject("base.kafka", &config.KafkaCfg)
		if err != nil {
			return
		}
	}

	producer, err := kafka.NewAsyncProducer()
	if err != nil {
		logger.Error("异常：%v", err.Error())
		return
	}

	// 发送消息
	msg := &sarama.ProducerMessage{
		Topic: "my_topic",
		Value: sarama.StringEncoder("Hello, world!"),
	}
	producer.Input() <- msg

	// 处理 Kafka 异步回调
	for {
		select {
		case <-producer.Successes():
			logger.Info("Message sent successfully")
		case err := <-producer.Errors():
			logger.Info("Error producing message: %s", err.Error())
		}
	}
}

func TestConsumerNew(t *testing.T) {
	config.LoadYamlFile("./application-simple.yaml")
	if config.GetValueBoolDefault("base.kafka.enable", false) {
		err := config.GetValueObject("base.kafka", &config.KafkaCfg)
		if err != nil {
			return
		}
	}

	// 创建 Kafka 消费者
	consumer, err := kafka.NewConsumer()
	if err != nil {
		logger.Error("%v", err.Error())
		return
	}

	// 订阅主题
	topic := "my_topic_demo"
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

func TestCreateTopicNew(t *testing.T) {
	config.LoadYamlFile("./application-admin.yaml")
	if config.GetValueBoolDefault("base.kafka.enable", false) {
		err := config.GetValueObject("base.kafka", &config.KafkaCfg)
		if err != nil {
			return
		}
	}

	admin, err := kafka.NewClusterAdmin()
	if err != nil {
		logger.Error("%v", err.Error())
	}
	defer admin.Close()

	// 创建 Topic
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	err = admin.CreateTopic("my_topic2", topicDetail, true)
	if err != nil {
		logger.Error("%v", err.Error())
	}
}
