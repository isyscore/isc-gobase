package test

import (
	"context"
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/extend/etcd"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/time"
	clientv3 "go.etcd.io/etcd/client/v3"
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
	etcdClient.Put(ctx, "test", time.TimeToStringYmdHms(time.Now()))
	rsp, _ := etcdClient.Get(ctx, "test")
	etcdClient.Get(ctx, "test", func(pOp *clientv3.Op) {
		fmt.Println("信息")
		fmt.Println(isc.ToJsonString(&pOp))
	})
	fmt.Println(rsp)
}
