package test

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/isyscore/isc-gobase/extend/emqx"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCut(t *testing.T) {
	fmt.Println(strings.Cut("oksfadf#sdf#fms", "#"))
	fmt.Println(strings.Cut("oksfadfdffms", "#"))

	urlFinal, _ := url.Parse("tcp://user:xxxsdf@localhost:8080")
	fmt.Println(urlFinal.User)
	fmt.Println(urlFinal.Path)
	fmt.Println(urlFinal.Host)
}

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

	time.Sleep(6 * time.Second)

	// 取消订阅
	if token := emqxClient.Unsubscribe("testtopic/#"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// 断开连接
	emqxClient.Disconnect(250)
	time.Sleep(1 * time.Second)
}
