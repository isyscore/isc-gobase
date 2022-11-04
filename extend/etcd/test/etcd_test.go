package test

import (
	"context"
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/extend/etcd"
	"testing"
)

func Test1(t *testing.T) {
	config.LoadYamlFile("./application-test1.yaml")
	if config.GetValueBoolDefault("base.etcd.enable", false) {
		err := config.GetValueObject("base.etcd", &config.EtcdCfg)
		if err != nil {
			return
		}
	}

	etcdClient, _ := etcd.NewEtcdClient()

	ctx := context.Background()
	etcdClient.Put(ctx, "gobase.k1", "dfsd")
	rsp, _ := etcdClient.Get(ctx, "gobase.k1")
	fmt.Println(rsp)
}
