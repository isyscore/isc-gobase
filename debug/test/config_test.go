package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/debug"
	"github.com/isyscore/isc-gobase/time"
	"testing"
	t0 "time"
)

func TestWatcher(t *testing.T) {
	config.LoadYamlFile("./application-test1.yaml")
	if config.GetValueBoolDefault("base.etcd.enable", false) {
		err := config.GetValueObject("base.etcd", &config.EtcdCfg)
		if err != nil {
			return
		}
	}

	debug.Init()
	debug.AddWatcher("test", func(key string, value string) {
		fmt.Println("有变化 key=", key, ", value=", value)
	})

	t0.Sleep(1000000000000)
}

func TestPush(t *testing.T) {
	config.LoadYamlFile("./application-test1.yaml")
	if config.GetValueBoolDefault("base.etcd.enable", false) {
		err := config.GetValueObject("base.etcd", &config.EtcdCfg)
		if err != nil {
			return
		}
	}

	debug.Init()
	var tim = time.TimeToStringYmdHmsS(time.Now())
	fmt.Println(tim)
	debug.Update("test", tim)
}
